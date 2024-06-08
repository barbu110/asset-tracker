package asset_manager

import (
	"asset-tracker/src/core/asset"
	"encoding/csv"
	"fmt"
	"io"
)

type CSV struct {
	Assets []asset.Asset
	rw     io.ReadWriter
}

func NewCSV(rw io.ReadWriter) (repo CSV, err error) {
	repo = CSV{rw: rw}
	if e := repo.readAll(rw); e != nil {
		err = e
		return
	}

	return
}

func (repo *CSV) readAll(r io.Reader) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV records: %w", err)
	}

	assets := make([]asset.Asset, len(records))
	for i, record := range records {
		if len(record) < 3 {
			return fmt.Errorf("record at line %v has less than 3 fields", i)
		}

		a, err := asset.FromExisting(record[0], record[1], record[2])
		if err != nil {
			return fmt.Errorf("record at line %v is invalid: %w", i, err)
		}

		assets[i] = a
	}

	repo.Assets = assets
	return nil
}

func (repo *CSV) Add(a asset.Asset) error {
	repo.Assets = append(repo.Assets, a)

	csvWriter := csv.NewWriter(repo.rw)
	if err := csvWriter.Write([]string{asset.EncodeIdToString(a.Id), a.Name, a.Description}); err != nil {
		return err
	}
	csvWriter.Flush()

	return nil
}
