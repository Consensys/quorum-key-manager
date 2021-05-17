package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func WriteJSON(req *http.Request, msg interface{}) error {
	// Prepare header
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set("Content-Type", "application/json")

	// Write message into buffer
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(b)

	// Set body
	req.Body = ioutil.NopCloser(body)
	req.ContentLength = int64(body.Len())
	buf := body.Bytes()
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(buf)
		return ioutil.NopCloser(r), nil
	}

	return nil
}
