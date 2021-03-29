package request

// Proxy creates a preparer for proxying request
func Proxy(cfg *ProxyConfig) Preparer {
	cfg.SetDefault()

	var preparers = []Preparer{
		Headers(cfg.Headers),
		ExtractURI(true),
		HTTPProtocol(1, 1),
		UserAgent(""),
	}

	if cfg.PassHostHeader != nil && !*cfg.PassHostHeader {
		preparers = append(preparers, Host(nil))
	}

	preparers = append(
		preparers,
		BasicAuth(cfg.BasicAuth),
		Body(),
	)

	return CombinePreparer(preparers...)
}
