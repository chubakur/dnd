package async

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chubakur/dnd/llmcore"
	"github.com/chubakur/dnd/transport"
)

type AsyncTask interface {
	Handle(t *transport.Transport) (*llmcore.MessageChain, error)
}

type AsyncTaskCommon struct {
	Type string `json:"type"`
}

func ParseAsync(taskType string, jdata string) (AsyncTask, error) {
	switch taskType {
	case "llmCall":
		var params AsyncTaskChatLlmStruct
		decoder := json.NewDecoder(strings.NewReader(jdata))
		err := decoder.Decode(&params)
		if err != nil {
			return nil, err
		}
		return &params, nil
	}
	return nil, fmt.Errorf("Invalid type: %s", taskType)
}
