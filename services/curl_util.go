package services

import (
	"bytes"
	"encoding/json"

	"github.com/geekymedic/neon-cli/types/xast"
)

func convertCurl(astTree *xast.TopNode) (string, error) {
	var data bytes.Buffer
	if err := json.NewEncoder(&data).Encode(InjectAstTree(astTree)); err != nil {
		return "", err
	}
	return data.String(), nil
}
