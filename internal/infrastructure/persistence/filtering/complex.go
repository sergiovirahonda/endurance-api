package filtering

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/domain/entities"
	"github.com/sergiovirahonda/endurance-api/internal/lib"
	"gorm.io/gorm"
)

// Structs

type ComplexFilters struct {
	context        echo.Context
	filters        map[string]interface{}
	orderBy        string
	orderDirection string
	pagination     ComplexFiltersPagination
}

type ComplexFiltersPagination struct {
	Page       int
	PageSize   int
	TotalItems int
}

func NewComplexFilter(
	ctx echo.Context,
	filters map[string]interface{},
	orderBy string,
	orderDirection string,
	page int,
	limit int,
) ComplexFilters {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if orderDirection == "" {
		orderDirection = "desc"
	}
	if page < 0 {
		page = 0
	}
	if limit <= 0 {
		limit = 100
	}
	cf := ComplexFilters{
		context:        ctx,
		filters:        filters,
		orderBy:        orderBy,
		orderDirection: orderDirection,
		pagination: ComplexFiltersPagination{
			Page:       page,
			PageSize:   limit,
			TotalItems: 0,
		},
	}
	cf.Validate()
	return cf
}

func (cf *ComplexFilters) Validate() error {
	if cf.orderBy == "" {
		return errors.New("order_by is required")
	}
	if cf.orderDirection == "" {
		return errors.New("order_direction is required")
	}
	if !lib.SliceContains([]string{"asc", "desc"}, cf.orderDirection) {
		return errors.New("order_direction must be either asc or desc")
	}
	return nil
}

func (cf *ComplexFilters) GetFilters() map[string]interface{} {
	return cf.filters
}

func (cf *ComplexFilters) OverrideFilters(filters map[string]interface{}) {
	cf.filters = filters
}

func (cf *ComplexFilters) AddFilter(key string, value interface{}) {
	cf.filters[key] = value
}

func (cf *ComplexFilters) AddFilters(filters map[string]interface{}) {
	for key, value := range filters {
		cf.filters[key] = value
	}
}

func (cf *ComplexFilters) SetOrderBy(orderBy string) {
	cf.orderBy = orderBy
	delete(cf.filters, "order_by")
}

func (cf *ComplexFilters) SetOrderDirection(orderDirection string) {
	cf.orderDirection = orderDirection
	delete(cf.filters, "order_direction")
}

func (cf *ComplexFilters) SetOffset(offset int) {
	cf.pagination.Page = offset
	delete(cf.filters, "offset")
}

func (cf *ComplexFilters) SetLimit(limit int) {
	cf.pagination.PageSize = limit
	delete(cf.filters, "limit")
}

func (cf *ComplexFilters) GetPagination() ComplexFiltersPagination {
	return cf.pagination
}

func (cf *ComplexFilters) GetOrdering() string {
	return fmt.Sprintf("%s %s", cf.orderBy, strings.ToUpper(cf.orderDirection))
}

func (cf *ComplexFilters) Paginate() map[string]interface{} {
	page, err := strconv.Atoi(cf.context.QueryParam("offset"))
	if err != nil {
		page = 0
	}
	if page < 0 {
		page = 0
	}
	pageSize, err := strconv.Atoi(cf.context.QueryParam("limit"))
	if err != nil {
		pageSize = 100
	} else {
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
	}
	orderBy := cf.context.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "created_at"
	}
	orderDirection := cf.context.QueryParam("order_direction")
	if orderDirection == "" {
		orderDirection = "desc"
	}
	if !lib.SliceContains([]string{"asc", "desc"}, orderDirection) {
		orderDirection = "desc"
	}
	cf.SetOrderBy(orderBy)
	cf.SetOrderDirection(orderDirection)
	cf.SetOffset(page)
	cf.SetLimit(pageSize)
	return cf.filters
}

func (cf *ComplexFilters) NarrowUserFilters(key string) error {
	if cf.context.Get("user") == nil {
		return errors.New("user not in context")
	}
	user := cf.context.Get("user").(*entities.User)
	keys := make([]string, 0, len(cf.filters))
	for key := range cf.filters {
		keys = append(keys, key)
	}
	switch key {
	case "user_id":
		if !lib.SliceContains(keys, "user_id") {
			cf.filters["user_id"] = user.ID
		} else {
			if user.ID != cf.filters["user_id"].(uuid.UUID) {
				return errors.New("forbidden")
			}
		}
	case "id":
		if !lib.SliceContains(keys, "id") {
			cf.filters["id"] = user.ID
		} else {
			if user.ID != cf.filters["id"].(uuid.UUID) {
				return errors.New("forbidden")
			}
		}
	default:
		return errors.New("unsupported")
	}
	return nil
}

func (cf *ComplexFilters) BindFilters() map[string]interface{} {
	logger := config.GetLoggerFromContext(cf.context)
	rawQuery := cf.context.Request().URL.RawQuery
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		logger.Errorf("Error parsing query: %v", err)
		return cf.filters
	}
	var unparsedFilters = make(map[string]interface{})
	for key, value := range values {
		val := strings.TrimSpace(value[0])
		val = strings.ReplaceAll(val, "*", "%")
		unparsedFilters[key] = val
	}
	var filters = make(map[string]interface{})
	// Parse boolean values
	for key, value := range unparsedFilters {
		if lib.SliceContains([]string{"true", "false"}, value.(string)) {
			boolValue, err := strconv.ParseBool(value.(string))
			if err != nil {
				logger.Errorf("Error parsing boolean: %v", err)
				continue
			}
			filters[key] = boolValue
		} else {
			filters[key] = value
		}
	}
	cf.filters = filters
	cf.Paginate()
	return cf.filters
}

func (cf *ComplexFilters) QueryFromFilter(query *gorm.DB) *gorm.DB {
	logger := config.GetLogger()
	filters := cf.GetFilters()
	for key, value := range filters {
		logger.Infof("Key: %s, Value: %v", key, value)
		if strings.HasSuffix(key, "__gt") {
			key = strings.ReplaceAll(key, "__gt", "")
			query = query.Where(key+" > ?", value)
		} else if strings.HasSuffix(key, "__gte") {
			key = strings.ReplaceAll(key, "__gte", "")
			query = query.Where(key+" >= ?", value)
		} else if strings.HasSuffix(key, "__lt") {
			key = strings.ReplaceAll(key, "__lt", "")
			query = query.Where(key+" < ?", value)
		} else if strings.HasSuffix(key, "__lte") {
			key = strings.ReplaceAll(key, "__lte", "")
			query = query.Where(key+" <= ?", value)
		} else if strings.HasSuffix(key, "__like") {
			key = strings.ReplaceAll(key, "__like", "")
			query = query.Where(
				fmt.Sprintf("unaccent(%s)", key)+" ILIKE unaccent(?)",
				value,
			)
		} else if strings.HasSuffix(key, "__in") {
			// If type of value is a string, split it by comma
			switch v := value.(type) {
			default:
				logger.Infof("Type: %T", v)
			case string:
				value = strings.Split(value.(string), ",")
			case []string:
				value = value.([]string)
			case []uuid.UUID:
				value = value.([]uuid.UUID)
			case []bool:
				value = value.([]bool)
			case []int:
				value = value.([]int)
			}
			key = strings.ReplaceAll(key, "__in", "")
			query = query.Where(key+" IN (?)", value)
		} else {
			if key != "offset" && key != "limit" && key != "order_by" && key != "order_direction" {
				query = query.Where(key+" = ?", value)
			}
		}
	}
	return query
}

// Raw filtering functions

func (cf *ComplexFilters) SetMetaParameters() map[string]interface{} {
	if cf.pagination.Page == 0 {
		cf.SetOffset(0)
	}
	if cf.pagination.PageSize == 0 {
		cf.SetLimit(100)
	}
	if cf.orderDirection == "" {
		cf.SetOrderDirection("desc")
	}
	if cf.orderBy == "" {
		cf.SetOrderBy("created_at")
	}
	return cf.filters
}
