package asset

import (
	"fmt"
)

type Asset struct {
	Id          Id
	Name        string
	Description string
	Properties  []CustomProperty
}

type CustomProperty struct {
	Name  string
	Value string
}

func New(name, description string, properties ...CustomProperty) Asset {
	return Asset{
		Id:          RandomId(),
		Name:        name,
		Description: description,
		Properties:  properties,
	}
}

func FromExisting(id, name, description string, properties ...CustomProperty) (a Asset, err error) {
	assetId, e := ParseId(id)
	if e != nil {
		err = fmt.Errorf("invalid ID: %w", e)
		return
	}

	return Asset{
		Id:          assetId,
		Name:        name,
		Description: description,
		Properties:  properties,
	}, nil
}
