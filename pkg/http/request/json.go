package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func WriteJSON(req *http.Request, msg interface{}) error {
	// Prepare header
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set("Content/Type", "application/json")

	// Write message into buffer
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(msg)
	if err != nil {
		return err
	}

	// Set body
	req.Body = ioutil.NopCloser(buf)
	req.ContentLength = int64(buf.Len())

	return nil
}
