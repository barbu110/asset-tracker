package asset

import "fmt"

const (
	KindUnspecified = iota
	KindItem
	KindContainer
)

type Kind int

func (k Kind) String() string {
	switch k {
	case KindUnspecified:
		return "UNSPECIFIED"
	case KindItem:
		return "ITEM"
	case KindContainer:
		return "CONTAINER"
	}

	panic(fmt.Sprintf("Unknown asset kind: %d", k))
}
