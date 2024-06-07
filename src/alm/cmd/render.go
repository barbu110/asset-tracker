package cmd

import (
	"asset-tracker/src/alm/asset"
	"asset-tracker/src/alm/label"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var renderName string
var renderDescription string
var renderWidth int
var renderHeight int
var renderOutputPath string

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
		renderedLabel, err := label.Render(&a, label.Params{
			Width:              renderWidth,
			Height:             renderHeight,
			Padding:            10,
			UseWhiteBackground: true,
		})
		if err != nil {
			return fmt.Errorf("failed to render label: %w", err)
		}
		if _, err := outputFile.Write(renderedLabel); err != nil {
			return fmt.Errorf("failed to write rendered label: %w", err)
		}

		return nil
	},
}
