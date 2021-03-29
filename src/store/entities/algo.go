package entities

type SigningAlgorithm string
type EllipticCurve string

const (
	Ecdsa SigningAlgorithm = "ecdsa"
	Eddsa SigningAlgorithm = "eddsa"
)

const (
	Bn256     EllipticCurve = "bn256"
	Secp256k1 EllipticCurve = "secp256k1"
)

type Algo struct {
	Type          SigningAlgorithm
	EllipticCurve string
	Size          int
}
