package cmd

import (
	"asset-tracker/src/asset_manager"
	"asset-tracker/src/core/asset"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var saveAssetCsvPath string
var saveAssetId string
var saveAssetName string
var saveAssetDescription string

var saveAssetCommand = &cobra.Command{
	Use:   "save_asset",
	Short: "Save an asset into the database.",
	RunE: func(_ *cobra.Command, _ []string) error {
		csvFile, err := os.OpenFile(saveAssetCsvPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("could not open CSV file: %w", err)
		}

		defer func() {
			if err := csvFile.Close(); err != nil {
				return
			}
		}()

		repo, err := asset_manager.NewCSV(csvFile)
		if err != nil {
			return fmt.Errorf("could not read assets: %w", err)
		}

		a, err := asset.FromExisting(saveAssetId, saveAssetName, saveAssetDescription)
		if err != nil {
			return fmt.Errorf("could not construct asset: %w", err)
		}

		if err := repo.Add(a); err != nil {
			return err
		}

		return nil
	},
}
