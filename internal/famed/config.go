package famed

type Config struct {
	Label    string
	Currency string
	Rewards  map[IssueSeverity]float64
}
