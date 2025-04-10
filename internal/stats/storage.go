package stats

type storage interface {
	All() []Stat

	Save(stat *Stat)

	Remove(stat *Stat)
}
