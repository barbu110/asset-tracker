package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var renderWidth int
var renderHeight int

var renderCommand = &cobra.Command{
	Use:   "render",
	Short: "Render the label of a given asset.",
	Run: func(cmd *cobra.Command, _ []string) {
		fmt.Printf("%v %v\n", renderWidth, renderHeight)
	},
}
