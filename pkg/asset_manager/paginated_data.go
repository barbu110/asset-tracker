package asset_manager

type PaginatedData[T any] struct {
	Items        []T
	NextToken    string
	HasNextToken bool
}
