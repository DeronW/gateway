package protocol

type CipherKey struct {
	UserKeyIndex int
	IV           string
	IvChr        string
	Iv96str      []byte
	EncryptCtr   uint32
	DecryptCtr   uint32
	UserKey      []byte
}
