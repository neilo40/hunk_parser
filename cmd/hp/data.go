package main

import (
	"fmt"
	"log/slog"
)

type HunkData HunkCode // data looks just like code
func (hc HunkData) Print(l *slog.Logger) {
	l.Debug("HunkData", "StarBytes", hc.Start, "NumWords", hc.NumWords, "NumBytes", hc.NumWords*4, "HunkType", hc.HunkType)
	fmt.Printf("┌--------------------------------┐ %0#4x (%4d)\n", hc.Start, hc.Start)
	fmt.Println("│         Data = 0x000003EA      │")
	fmt.Println("├--------------------------------┤")
	fmt.Printf("│     Length = %4d longwords    │\n", hc.NumWords)
	fmt.Println("├--------------------------------┤")
	fmt.Println("                ...               ")
	fmt.Printf("└--------------------------------┘%0#4x (%4d)\n", hc.Start+(int(hc.NumWords)*4)+8, hc.Start+(int(hc.NumWords)*4)+8)
}

// TODO - avoid this duplication.  code and data look the same so make a generic function?
func (hp *HunkParser) readHunkData(i int) (int, error) {
	start := i - 4
	numWords, i := hp.readUint32(i)
	words := make([]uint32, 0, numWords)
	for wc := 0; wc < int(numWords); wc++ {
		var word uint32
		word, i = hp.readUint32(i)
		words = append(words, word)
	}

	h := HunkData{
		Start:    start,
		HunkType: HUNK_DATA,
		NumWords: numWords,
		Content:  words,
	}

	hp.HunkFile.Hunks = append(hp.HunkFile.Hunks, h)

	h.Print(hp.Logger)
	return i, nil
}
