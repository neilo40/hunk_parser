package main

import (
	"fmt"
	"log/slog"
)

type HunkData HunkCode // data looks just like code
func (hc HunkData) Print(l *slog.Logger) {
	l.Debug("HunkData", "StarBytes", hc.Start, "NumWords", hc.NumWords, "NumBytes", hc.NumWords*4, "HunkType", hc.HunkType)
	fmt.Println("┌--------------------------------┐")
	fmt.Println("│         Data = 0x000003EA      │")
	fmt.Println("├--------------------------------┤")
	fmt.Printf("│  Length = %d longwords       │\n", hc.NumWords) // TODO: align border
	fmt.Println("├--------------------------------┤")
	fmt.Println("                ...               ")
	fmt.Println("└--------------------------------┘")
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
