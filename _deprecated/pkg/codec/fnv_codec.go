package codec

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash"
	"hash/fnv"
)

type FnvCodec[T any] struct {
	hasher hash.Hash64
}

func NewFnvCodec[T any]() Codec[T] {
	return &FnvCodec[T]{
		hasher: fnv.New64a(),
	}
}

func (c *FnvCodec[T]) Hash(data []T) (string, error) {
	c.hasher.Reset()	
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	for _, d := range data {
		err := encoder.Encode(d)
		if err != nil {
			return "", fmt.Errorf("encode: %v", err)
		}
		c.hasher.Write(buf.Bytes())
	}
	return fmt.Sprintf("%x", c.hasher.Sum64()), nil
}
