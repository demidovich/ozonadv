package ozon

type Debug interface {
	Println(m ...any)
	Printf(format string, m ...any)
}
