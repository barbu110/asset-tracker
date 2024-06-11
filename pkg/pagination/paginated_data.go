package pagination

type PaginatedData[T any] struct {
	Items        []T
	NextToken    string
	HasNextToken bool
}

func NewEmpty[T any]() PaginatedData[T] {
	return PaginatedData[T]{}
}
