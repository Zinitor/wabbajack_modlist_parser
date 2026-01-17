package parser

type Parser[T any] interface {
	Parse() []T
}

func GetBase[T any](p Parser[T]) <-chan []T {
	ch := make(chan []T, 1)
	go func() {
		defer close(ch)
		ch <- p.Parse()
	}()
	return ch
}
