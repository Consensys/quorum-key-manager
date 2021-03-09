package utils

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)


func WaitForServiceLive(ctx context.Context, url, name string, timeout time.Duration) {
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		req, _ := http.NewRequest("GET", url, nil)
		req = req.WithContext(rctx)

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			if resp != nil && resp.StatusCode == 200 {
				log.Infof("service %s is live", name)
				return
			}

			log.WithField("status", resp.StatusCode).Warnf("cannot reach %s service", name)
		}

		if rctx.Err() != nil {
			log.WithError(rctx.Err()).Warnf("cannot reach %s service", name)
			return
		}

		log.Debugf("waiting for 1 s for service %s to start...", name)
		time.Sleep(time.Second)
	}
}
