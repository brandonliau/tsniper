package codec

type Codec interface {
	Hash(data [][]string) string
}
