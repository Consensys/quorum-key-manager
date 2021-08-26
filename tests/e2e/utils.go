package e2e

import (
	"time"

	"github.com/consensys/quorum-key-manager/pkg/client"
)

const MAX_RETRIES = 5

type callFunc func() error
type logFunc func(format string, args ...interface{})

func retryOn(call callFunc, logger logFunc, errMsg string, httpStatusCode, retries int) error {
	for {
		err := call()
		if httpError, ok := err.(*client.ResponseError); retries <= 0 || !ok || httpError.StatusCode != httpStatusCode {
			if err != nil {
				return err
			}
			break
		}

		logger("%s (retrying in 1 second...)", errMsg)
		time.Sleep(time.Second)
		retries--
	}
	
	return nil
}
