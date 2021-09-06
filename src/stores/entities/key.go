package entities

// Key public part of a key
type Key struct {
	ID          string
	PublicKey   []byte
	Algo        *Algorithm
	Metadata    *Metadata
	Tags        map[string]string
	Annotations *Annotation
}

func (k *Key) IsETHAccount() bool {
	return k.Algo.EllipticCurve == Secp256k1 && k.Algo.Type == Ecdsa
}
