package main

import (
	"fmt"
	"log/slog"
)

type HunkCode struct {
	Start    int
	HunkType uint32
	NumWords uint32
	Content  []uint32
}

func (hc HunkCode) Print(l *slog.Logger) {
	l.Debug("HunkCode", "StartBytes", hc.Start, "NumWords", hc.NumWords, "NumBytes", hc.NumWords*4,
		"HunkType", hc.HunkType)
	fmt.Printf("┌--------------------------------┐ %0#4x (%4d)\n", hc.Start, hc.Start)
	fmt.Println("|         Code = 0x000003E9      |")
	fmt.Println("├--------------------------------┤")
	fmt.Printf("|     Length = %4d longwords    |\n", hc.NumWords)
	fmt.Println("├--------------------------------┤")
	fmt.Println("                ...               ")
	fmt.Printf("└--------------------------------┘ %0#4x (%4d)\n", hc.Start+(int(hc.NumWords)*4)+8, hc.Start+(int(hc.NumWords)*4)+8)
	// TODO: write the code section somwehere so we can disassemble it

}

func (hp *HunkParser) readHunkCode(i int) (int, error) {
	start := i - 4 // include the header
	numWords, i := hp.readUint32(i)
	words := make([]uint32, 0, numWords)
	for wc := 0; wc < int(numWords); wc++ {
		var word uint32
		word, i = hp.readUint32(i)
		words = append(words, word)
	}

	h := HunkCode{
		Start:    start,
		HunkType: HUNK_CODE,
		NumWords: numWords,
		Content:  words,
	}

	hp.HunkFile.Hunks = append(hp.HunkFile.Hunks, h)

	h.Print(hp.Logger)
	return i, nil
}
