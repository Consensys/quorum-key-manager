package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadJSON(resp *http.Response, msg interface{}) error {
	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(msg)
	default:
		return fmt.Errorf("invalid response Content-Type: %v", contentType)
	}
}
