package async

import "github.com/chubakur/dnd/transport"

type AsyncTask interface {
	GetType() string
	Handle(t *transport.Transport) (string, error)
}

type AsyncTaskCommon struct {
	Type string `json:"type"`
}
