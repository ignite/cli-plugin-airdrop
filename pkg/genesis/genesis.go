package genesis

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

type GenState map[string]json.RawMessage

// GetGenStateFromPath returns a JSON genState message from inputted path.
func GetGenStateFromPath(genesisFilePath string) (genState GenState, err error) {
	genesisFile, err := os.Open(filepath.Clean(genesisFilePath))
	if err != nil {
		return genState, err
	}
	defer genesisFile.Close()

	byteValue, err := io.ReadAll(genesisFile)
	if err != nil {
		return genState, err
	}

	var doc tmtypes.GenesisDoc
	err = tmjson.Unmarshal(byteValue, &doc)
	if err != nil {
		return genState, err
	}

	err = json.Unmarshal(doc.AppState, &genState)
	if err != nil {
		return genState, err
	}
	return genState, nil
}
