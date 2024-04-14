package main

import "encoding/binary"

func (hp *HunkParser) readStrings(i int) ([]string, int, error) {
	s := make([]string, 0)
	// TODO: actually convert each 4 byte sequence to chars
	return s, i + 4, nil
}

func (hp *HunkParser) readUint32(i int) (uint32, int) {
	n := binary.BigEndian.Uint32(hp.FileBytes[i : i+4])
	return n, i + 4
}
