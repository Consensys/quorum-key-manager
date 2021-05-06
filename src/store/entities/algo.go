package entities

type Curve string
type KeyType string

const (
	Ecdsa KeyType = "ecdsa"
	Eddsa KeyType = "eddsa"

	Bn254     Curve = "bn254"
	Secp256k1 Curve = "secp256k1"
)

type Algorithm struct {
	Type          KeyType
	EllipticCurve Curve
}
