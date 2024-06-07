package asset

import "encoding"

type Asset struct {
	Id          encoding.BinaryMarshaler
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
