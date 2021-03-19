package entities

type SigningAlgorithm string
type EllipticCurve string

const (
	Ecdsa SigningAlgorithm = "ecdsa"
	Eddsa SigningAlgorithm = "eddsa"
)

const (
	Bn256 EllipticCurve = "bn256"
	Secp256k1 EllipticCurve = "secp256k1"
)

// Algo
type Algo struct {
	// Type of key (e.g. ecdsa, rsa)
	Type SigningAlgorithm

	// EllipticCurve (e.g secp256k1)
	EllipticCurve string

	// Size of the key
	Size int
}
