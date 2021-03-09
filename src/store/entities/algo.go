package entities

// Algo
type Algo struct {
	// Type of key (e.g. ecdsa, rsa)
	Type string

	// EllipticCurve (e.g secp256k1)
	EllipticCurve string

	// Size of the key
	Size int
}

