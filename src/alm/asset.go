package main

import (
	"bytes"
	svg "github.com/ajstarks/svgo"
	"github.com/google/uuid"
)

type Asset struct {
	Id          uuid.UUID
	Name        string
	Description string
	Properties  []AssetProperty
}

type AssetProperty struct {
	Name  string
	Value string
}

type RenderLabelParams struct {
	width  int
	height int
}

func (asset *Asset) RenderLabel(params RenderLabelParams) []byte {
	buffer := bytes.Buffer{}
	canvas := svg.New(&buffer)

	canvas.Start(params.width, params.height)
	canvas.End()

	return buffer.Bytes()
}
