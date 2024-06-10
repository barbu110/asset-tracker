package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use:              "alm",
	Short:            "Asset Label Maker",
	TraverseChildren: true,
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	renderCommand.Flags().StringVarP(&renderName, "name", "n", "", "Name of the asset.")
	_ = renderCommand.MarkFlagRequired("name")
	renderCommand.Flags().StringVarP(&renderDescription, "description", "d", "", "One-line description of the asset.")
	_ = renderCommand.MarkFlagRequired("description")

	renderCommand.Flags().IntVar(&renderWidth, "width", 400, "Width of the label, in pixels.")
	renderCommand.Flags().IntVar(&renderHeight, "height", 300, "Height of the label, in pixels.")

	renderCommand.Flags().StringVarP(&renderOutputPath, "output_path", "o", "", "Path to output the rendered label.")
	_ = renderCommand.MarkFlagRequired("output_path")
	_ = renderCommand.MarkFlagFilename("output_path")

	renderCommand.Flags().BoolVar(&renderPng, "png", false, "When passed, label is rendered to PNG.")

	rootCommand.AddCommand(renderCommand)

	saveAssetCommand.Flags().StringVarP(&saveAssetCsvPath, "output_path", "o", "", "Path to a CSV file where assets are to be appended.")
	_ = saveAssetCommand.MarkFlagRequired("output_path")
	_ = saveAssetCommand.MarkFlagFilename("output_path")

	saveAssetCommand.Flags().StringVar(&saveAssetId, "asset_id", "", "ID of the asset.")
	_ = saveAssetCommand.MarkFlagRequired("asset_id")

	saveAssetCommand.Flags().StringVar(&saveAssetName, "asset_name", "", "Name of the asset.")
	_ = saveAssetCommand.MarkFlagRequired("asset_name")

	saveAssetCommand.Flags().StringVar(&saveAssetDescription, "asset_description", "", "Description of the asset.")
	_ = saveAssetCommand.MarkFlagRequired("asset_description")

	rootCommand.AddCommand(saveAssetCommand)
}
