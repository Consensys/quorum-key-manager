package entities

const (
	Ecdsa string = "ecdsa"
	Eddsa string = "eddsa"

	Bn254     string = "bn254"
	Secp256k1 string = "secp256k1"
)

type Algorithm struct {
	Type          string
	EllipticCurve string
}
