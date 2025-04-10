package cabinets

type storage interface {
	All() []Cabinet

	Save(cabinet Cabinet)

	Remove(cabinet Cabinet)
}
