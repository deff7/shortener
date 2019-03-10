package main

import (
	"bytes"
)

var base62 = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Generator struct {
	charMap map[byte]uint64
}

func NewGenerator() *Generator {
	g := &Generator{
		charMap: map[byte]uint64{},
	}
	for i, b := range base62 {
		g.charMap[b] = uint64(i)
	}
	return g
}

func (g *Generator) EncodeID(id uint64) []byte {
	var buf bytes.Buffer
	for id >= 0 {
		b := byte(int(id) % len(base62))
		buf.WriteByte(base62[b])
		id /= uint64(len(base62))

		if id == 0 {
			break
		}
	}
	return buf.Bytes()
}

func (g *Generator) DecodeID(raw []byte) uint64 {
	var (
		id    = uint64(0)
		digit = uint64(1)
	)
	for _, b := range raw {
		id += g.charMap[b] * digit
		digit *= uint64(len(base62))
	}
	return id
}
