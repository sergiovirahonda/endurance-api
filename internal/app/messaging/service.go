package messaging

import "github.com/labstack/echo/v4"

type MessagingService interface {
	EventsLoop()
	HandleDomainEvent(ctx echo.Context, msg []byte) error
}
