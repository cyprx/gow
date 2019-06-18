package pool

type Work interface {
	Execute() interface{}
}
