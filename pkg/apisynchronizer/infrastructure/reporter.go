package infrastructure

type Reporter interface {
	Info(...interface{})
}
