package asset

import (
	"fmt"
	"time"
)

type Asset struct {
	Id Id
	// ID of another asset (e.g. a box) in which this one is placed.
	ContainerId           Id
	Kind                  Kind
	Name                  string
	Description           string
	Properties            []CustomProperty
	CreatedAt, ModifiedAt time.Time
}

type CustomProperty struct {
	Name  string
	Value string
}

func NewItem(name, description string, containerID *Id, properties ...CustomProperty) Asset {
	actualContainerId := RootContainerId()
	if containerID != nil {
		actualContainerId = *containerID
	}

	return Asset{
		Id:          RandomId(),
		ContainerId: actualContainerId,
		Kind:        KindItem,
		Name:        name,
		Description: description,
		Properties:  properties,
		CreatedAt:   time.Now().UTC(),
		ModifiedAt:  time.Now().UTC(),
	}
}

func NewContainer(name, description string, containerID *Id, properties ...CustomProperty) Asset {
	actualContainerId := RootContainerId()
	if containerID != nil {
		actualContainerId = *containerID
	}

	return Asset{
		Id:          RandomId(),
		ContainerId: actualContainerId,
		Kind:        KindContainer,
		Name:        name,
		Description: description,
		Properties:  properties,
		CreatedAt:   time.Now().UTC(),
		ModifiedAt:  time.Now().UTC(),
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
