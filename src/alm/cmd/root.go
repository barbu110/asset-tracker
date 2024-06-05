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
	renderCommand.Flags().IntVar(&renderWidth, "width", 400, "Width of the label, in pixels.")
	renderCommand.Flags().IntVar(&renderHeight, "height", 300, "Height of the label, in pixels.")

	rootCommand.AddCommand(renderCommand)
}
