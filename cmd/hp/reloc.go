package main

import (
	"fmt"
	"log/slog"
)

type HunkReloc32 struct {
	Relocs []Reloc32
}

type Reloc32 struct {
	Start      int
	HunkType   uint32
	NumOffsets uint32
	HunkNum    uint32
	Offsets    []uint32
}

func (hr HunkReloc32) Print(l *slog.Logger) {
	start := hr.Relocs[0].Start
	for _, r := range hr.Relocs {
		l.Debug("HunkReloc32", "StartBytes", r.Start, "NumOffsets", r.NumOffsets, "HunkType", r.HunkType, "HunkNum", r.HunkNum, "Offsets", r.Offsets)
	}
	fmt.Printf("┌--------------------------------┐ %0#4x (%4d)\n", start, start)
	fmt.Println("|       Reloc32 = 0x000003EC     |")
	fmt.Println("├--------------------------------┤")
	end := start + 4
	for i, r := range hr.Relocs {
		fmt.Printf("|       Reloc %d                  |\n", i+1)
		fmt.Println("├--------------------------------┤")
		fmt.Printf("|       Hunk Number %d            |\n", r.HunkNum)
		fmt.Println("├--------------------------------┤")
		end += 8
		for _, o := range r.Offsets {
			fmt.Printf("| Offset %4d                    |\n", o)
			end++
		}
		if i < len(hr.Relocs)-1 {
			fmt.Println("├--------------------------------┤")
		}
	}
	fmt.Println("├--------------------------------┤")
	fmt.Println("|               0                |")
	fmt.Println("├--------------------------------┤")
	fmt.Println("|         END = 0x000003F2       |")
	fmt.Printf("└--------------------------------┘ %0#4x(%4d)\n", end+8, end+8)
}

func (hp *HunkParser) readHunkReloc32(i int) (int, error) {
	relocs := make([]Reloc32, 0)
	for {
		start := i

		var numOffsets uint32
		numOffsets, i = hp.readUint32(i)
		if numOffsets == 0 {
			break
		}

		var offsetNum uint32
		offsetNum, i = hp.readUint32(i)

		offsets := make([]uint32, 0, numOffsets)
		for oc := 0; oc < int(numOffsets); oc++ {
			var offset uint32
			offset, i = hp.readUint32(i)
			offsets = append(offsets, offset)
		}

		r := Reloc32{
			Start:      start,
			HunkType:   HUNK_RELOC32,
			NumOffsets: numOffsets,
			HunkNum:    offsetNum,
			Offsets:    offsets,
		}
		relocs = append(relocs, r)
	}

	// check next uint32 is HUNK_END
	end, i := hp.readUint32(i)
	if end != HUNK_END {
		return i, errHunkEndNotSeen
	}

	h := HunkReloc32{Relocs: relocs}

	hp.HunkFile.Hunks = append(hp.HunkFile.Hunks, h)
	h.Print(hp.Logger)

	return i, nil
}
