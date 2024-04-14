package main

import (
	"fmt"
	"log/slog"
)

type HunkHeader struct {
	Start        int
	HunkLen      int
	LibraryNames []string
	TableSize    uint32
	FirstHunk    uint32
	LastHunk     uint32
	HunkSizes    []uint32 // TODO: size is in words or bytes?
}

func (hh HunkHeader) Print(l *slog.Logger) {
	l.Debug("HunkHeader", "StartBytes", hh.Start, "HunkLen", hh.HunkLen, "FirstHunk", hh.FirstHunk,
		"LastHunk", hh.LastHunk, "TableSize", hh.TableSize, "HunkSizes", hh.HunkSizes)
	fmt.Println("┌--------------------------------┐")
	fmt.Println("|       Header = 0x000003F3      |")
	fmt.Println("├--------------------------------┤")
	for _, lib := range hh.LibraryNames {
		fmt.Printf("| TODO: support libs %s |", lib) // TODO: support this
	}
	fmt.Println("|                0               |")
	fmt.Println("├--------------------------------┤")
	fmt.Printf("| Table Size = %d                 |\n", hh.TableSize)
	fmt.Println("├--------------------------------┤")
	fmt.Printf("| First Hunk = %d                 |\n", hh.FirstHunk)
	fmt.Println("├--------------------------------┤")
	fmt.Printf("| Last Hunk = %d                  |\n", hh.LastHunk)
	fmt.Println("├--------------------------------┤")
	for _, hs := range hh.HunkSizes {
		fmt.Printf("| Hunk Size %d                 |\n", hs) // TODO: figure out how to align the right border
	}
	fmt.Println("└--------------------------------┘")

}

func (hp *HunkParser) readMagicCookie() (int, error) {
	mc, i := hp.readUint32(0)
	if mc == HUNK_HEADER {
		hp.HunkFile.MagicCookie = mc
		return i, nil
	} else {
		hp.Logger.Error("Magic Cookie", "status", "fail", "expected", HUNK_HEADER, "got", mc)
		return 0, errBadMagicCookie
	}
}

func (hp *HunkParser) readHeader() (int, error) {
	hunkLen := 0
	hp.HunkFile.Header.Start = 0
	i, err := hp.readMagicCookie()
	if err != nil {
		panic(err)
	}
	hunkLen += 4

	libNames, i, err := hp.readStrings(i)
	if err != nil {
		hp.Logger.Error("Reading library names", "error", err)
		return 0, nil
	}
	hp.HunkFile.Header.LibraryNames = libNames
	hunkLen += 4 // TODO: support library names in determining hunk length

	tableSize, i := hp.readUint32(i)
	hp.HunkFile.Header.TableSize = tableSize
	firstHunk, i := hp.readUint32(i)
	hp.HunkFile.Header.FirstHunk = firstHunk
	lastHunk, i := hp.readUint32(i)
	hp.HunkFile.Header.LastHunk = lastHunk
	hunkLen += 12

	numHunks := lastHunk - firstHunk
	hunkSizes := make([]uint32, 0, numHunks)
	for hn := 0; hn <= int(numHunks); hn++ {
		var hs uint32
		hs, i = hp.readUint32(i)
		// TODO: deal with mem flags
		hunkSizes = append(hunkSizes, hs)
	}
	hp.HunkFile.Header.HunkSizes = hunkSizes
	hp.HunkFile.Hunks = make([]HunkBlock, 0, numHunks)
	hunkLen += 4 * (int(numHunks) + 1)

	hp.HunkFile.Header.HunkLen = hunkLen
	hp.HunkFile.Header.Print(hp.Logger)

	return i, nil
}
