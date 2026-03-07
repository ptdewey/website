package generator

import (
	"encoding/json"
	"os"
)

func WriteJSONFile(data any, outputPath string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, jsonData, 0644)
}
