package beinagrind

import (
	"bytes"
	"encoding/gob"
)

type Hnit struct {
	PosX int
	PosY int
}

type Leikur struct {
	Hnit
}

func (leikur Leikur) ToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	err := gob.
		NewEncoder(&buf).
		Encode(leikur)

	return buf.Bytes(),  err
}

func (leikur Leikur) FromBytes(leikjaByte []byte) (Leikur, error) {
	return leikur, gob.
		NewDecoder(bytes.NewReader(leikjaByte)).
		Decode(&leikur)
}