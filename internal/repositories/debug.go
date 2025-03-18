package repositories

type DebugRepository interface {
	SizeMB() (float64, error)
}
