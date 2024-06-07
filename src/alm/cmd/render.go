package cmd

import (
	"asset-tracker/src/alm/label"
	"asset-tracker/src/core/asset"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var renderName string
var renderDescription string
var renderWidth int
var renderHeight int
var renderOutputPath string
var renderPng bool

var renderCommand = &cobra.Command{
	Use:   "render",
	Short: "Render the label of a given asset.",
	RunE: func(_ *cobra.Command, _ []string) error {
		outputFile, err := os.Create(renderOutputPath)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer func() {
			if err := outputFile.Close(); err != nil {
				panic(err)
			}
		}()

		a := asset.New(renderName, renderDescription)
		renderParams := label.Params{
			Width:              renderWidth,
			Height:             renderHeight,
			Padding:            10,
			UseWhiteBackground: true,
		}

		render := func() ([]byte, error) {
			if renderPng {
				fmt.Printf("Rendering label as PNG.\n")
				return label.RenderRaster(&a, renderParams)
			} else {
				fmt.Printf("Rendering label as SVG.\n")
				return label.RenderVector(&a, renderParams)
			}
		}
		renderedLabel, err := render()
		if err != nil {
			return fmt.Errorf("failed to render label: %w", err)
		}
		if _, err := outputFile.Write(renderedLabel); err != nil {
			return fmt.Errorf("failed to write rendered label: %w", err)
		}

		return nil
	},
}
