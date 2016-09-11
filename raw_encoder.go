package funnel

import (
	"fmt"
	"time"
)

func rawLogTransformer(src []byte, message []byte) ([]byte, error) {
	payload := fmt.Sprintf("%d %s %s", time.Now().UnixNano(), string(src), string(message))
	return []byte(payload), nil
}
