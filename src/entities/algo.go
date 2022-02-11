package entities

type Curve string
type KeyType string

const (
	Ecdsa KeyType = "ecdsa"
	Eddsa KeyType = "eddsa"

	Babyjubjub Curve = "babyjubjub"
	Secp256k1  Curve = "secp256k1"
	X25519     Curve = "x25519"
)

type Algorithm struct {
	Type          KeyType
	EllipticCurve Curve
}
