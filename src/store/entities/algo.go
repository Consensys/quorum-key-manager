package entities

const (
	Ecdsa string = "ecdsa"
	Eddsa string = "eddsa"

	Bn256     string = "bn256"
	Secp256k1 string = "secp256k1"
)

type Algorithm struct {
	Type          string
	EllipticCurve string
}
