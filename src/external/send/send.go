package send

import (
	"encoding/json"
	"fmt"
	"time"

	"Robot-Project/external/router"
	"Robot-Project/internal/model"
)

// SendMessage builds a command message and posts it to the given bot ID.
func Message(host, id string, cmd model.Command, arg string, from string) ([]byte, int, error) {
	msg := model.Message{
		From:    from,
		Payload: cmd,
		Arg:     arg,
		SentAt:  time.Now(),
	}

	pack, err := json.Marshal(msg)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal message: %w", err)
	}

	return router.Post(host, "/id/"+id+"/message", pack)
}
