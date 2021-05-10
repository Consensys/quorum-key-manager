package response

// Proxy creates a preparer for proxying request
func Proxy(cfg *ProxyConfig) Modifier {
	var modifiers = []Modifier{
		BackendServer(),
		Headers(cfg.Headers),
		GZIP(),
	}

	return CombineModifier(modifiers...)
}
