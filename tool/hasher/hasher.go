package hasher

type Hasher interface {
	Hash(str string) (string, error)
	Compare(str string, hash string) error
}
