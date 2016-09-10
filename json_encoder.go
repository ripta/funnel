package funnel

import (
	"encoding/json"
	"time"
)

func jsonLogTransformer(src []byte, message []byte) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"@timestamp": time.Now().UnixNano(),
		"@source":    string(src),
		"message":    string(message),
	})
}
