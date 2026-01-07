package structs

type Parser[T any] interface {
	Parse(jsonData []byte) []T
}
