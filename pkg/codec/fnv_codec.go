package codec

import (
	"fmt"
	"hash"
	"hash/fnv"
)

type FnvCodec struct {
	hasher hash.Hash64
}

func NewFnvCodec() *FnvCodec {
	return &FnvCodec{
		hasher: fnv.New64a(),
	}
}

func (c *FnvCodec) Hash(data [][]string) string {
	c.hasher.Reset()
	for _, row := range data {
		for _, str := range row {
			_, _ = c.hasher.Write([]byte(str))
		}
	}
	return fmt.Sprintf("%x", c.hasher.Sum64())
}
