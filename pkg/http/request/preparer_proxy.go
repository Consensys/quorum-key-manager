package request

import "net/url"

// Proxy creates a preparer for proxying request
func Proxy(cfg *ProxyConfig) (Preparer, error) {
	cfg = cfg.Copy().SetDefault()

	var preparers []Preparer
	if cfg.Addr != "" {
		u, err := url.Parse(cfg.Addr)
		if err != nil {
			return nil, err
		}
		preparers = append(preparers, URL(u))
	}

	preparers = append(
		preparers,
		Headers(cfg.Headers),
		ExtractURI(true),
		HTTPProtocol(1, 1),
		UserAgent(""),
	)

	if cfg.PassHostHeader != nil && !*cfg.PassHostHeader {
		preparers = append(preparers, Host(nil))
	}

	preparers = append(
		preparers,
		BasicAuth(cfg.BasicAuth),
		Body(),
		WebSocketHeaders(),
	)

	return CombinePreparer(preparers...), nil
}
