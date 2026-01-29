package backup

type Strategy string

const (
	StrategyFull         Strategy = "full"
	StrategyIncremental  Strategy = "incremental"
	StrategyDifferential Strategy = "differential"
)