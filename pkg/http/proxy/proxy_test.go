package proxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oxtoacart/bpool"
)

type staticTransport struct {
	res *http.Response
}

func (t *staticTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.res.Request = r
	return t.res, nil
}

func BenchmarkProxy(b *testing.B) {
	res := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://foo.bar/", nil)

	proxy, err := New(
		&Config{PassHostHeader: Bool(false)},
		&staticTransport{res},
		bpool.NewBytePool(32, 1024),
	)

	if err != nil {
		b.Errorf("Could not build proxy: %v", err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		proxy.ServeHTTP(w, req)
	}
}
