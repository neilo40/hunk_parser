package main

import (
	"fmt"
	"log/slog"
)

type HunkBSS struct {
	Start     int
	BlockSize uint32
}

func (hb HunkBSS) Print(l *slog.Logger) {
	l.Debug("HunkBSS", "StartBytes", hb.Start, "BlockSize", hb.BlockSize)
	fmt.Println("┌--------------------------------┐")
	fmt.Println("|         BSS = 0x000003EB       |")
	fmt.Println("├--------------------------------┤")
	fmt.Printf("|  Length = %d longwords       |\n", hb.BlockSize) // TODO: align border
	fmt.Println("├--------------------------------┤")
	fmt.Println("                ...               ")
	fmt.Println("└--------------------------------┘")
}

func (hp *HunkParser) readHunkBss(i int) (int, error) {
	start := i - 4
	blockSize, i := hp.readUint32(i)

	h := HunkBSS{
		Start:     start,
		BlockSize: blockSize,
	}
	hp.HunkFile.Hunks = append(hp.HunkFile.Hunks, h)

	h.Print(hp.Logger)

	// check next uint32 is HUNK_END
	end, i := hp.readUint32(i)
	if end != HUNK_END {
		return i, errHunkEndNotSeen
	}

	return i, nil
}
