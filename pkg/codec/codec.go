package codec

type Codec[T any] interface {
	Hash(data []T) (string, error)
}
