package valueobjects

type SymbolScore struct {
	Symbol  string
	Score   float64
	Ranking int
}

type SymbolScores []SymbolScore
