package main

// references: http://amiga-dev.wikidot.com/file-format:hunk

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"os"
)

type HunkParser struct {
	Logger    *slog.Logger
	HunkFile  *HunkFile
	FileBytes []byte
}

type HunkFile struct {
	MagicCookie uint32
	Header      HunkHeader
	Hunks       []Hunk
}

type HunkHeader struct {
	Start        int
	LibraryNames []string
	TableSize    uint32
	FirstHunk    uint32
	LastHunk     uint32
	HunkSizes    []uint32
}

type Hunk interface{}
type HunkCode struct {
	Start    int
	HunkType uint32
	NumWords uint32
	Content  []uint32
}
type HunkReloc32 struct {
	Start      int
	HunkType   uint32
	NumOffsets uint32
	HunkNum    uint32
	Offsets    []uint32
}
type HunkEnd struct {
	Start    int
	HunkType uint32
}

const (
	HUNK_HEADER  = 0x000003F3
	HUNK_CODE    = 0x000003E9
	HUNK_DATA    = 0x000003EA
	HUNK_BSS     = 0x000003EB
	HUNK_RELOC32 = 0x000003EC
	HUNK_END     = 0x000003F2
)

var (
	errBadMagicCookie = errors.New("bad magic cookie")
)

func main() {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, logOpts)
	logger := slog.New(handler)
	hp := HunkParser{
		Logger:   logger,
		HunkFile: &HunkFile{},
	}

	err := hp.loadHunkFile("../../test/Vim")
	if err != nil {
		panic(err)
	}

	i, err := hp.readMagicCookie()
	if err != nil {
		panic(err)
	}

	i, err = hp.readHeader(i)
	if err != nil {
		panic(err)
	}

	for range hp.HunkFile.Header.HunkSizes {
		var hunkType uint32
		hunkType, i = hp.readUint32(i)
		switch hunkType {
		// TODO: Implement all the hunk types
		case HUNK_CODE:
			i, err = hp.readHunkCode(i)
			if err != nil {
				panic(err)
			}
		case HUNK_RELOC32:
			i, err = hp.readHunkReloc32(i)
			if err != nil {
				panic(err)
			}
		default:
			slog.Error("Unknown hunk type", "type", hunkType, "byte", i)
			panic(true)
		}
	}

	hp.Logger.Info("Parsing Complete", "bytes_parsed", i, "total_bytes", len(hp.FileBytes))
	hp.Logger.Debug("hunkfile", "content", *hp.HunkFile)
}

func (hp *HunkParser) readStrings(i int) ([]string, int, error) {
	s := make([]string, 0)
	// TODO: actually convert each 4 byte sequence to chars
	return s, i + 4, nil
}

func (hp *HunkParser) readUint32(i int) (uint32, int) {
	n := binary.BigEndian.Uint32(hp.FileBytes[i : i+4])
	return n, i + 4
}

func (hp *HunkParser) loadHunkFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		hp.Logger.Error("Could not open file", "filename", filename, "error", err)
		return err
	}

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		hp.Logger.Error("Could not read file", "filename", filename, "error", err)
		return err
	}

	hp.FileBytes = fileBytes
	return nil
}

func (hp *HunkParser) readMagicCookie() (int, error) {
	mc, i := hp.readUint32(0)
	if mc == HUNK_HEADER {
		hp.Logger.Info("Magic Cookie", "status", "ok")
		hp.HunkFile.MagicCookie = mc
		return i, nil
	} else {
		hp.Logger.Error("Magic Cookie", "status", "fail", "expected", HUNK_HEADER, "got", mc)
		return 0, errBadMagicCookie
	}
}

func (hp *HunkParser) readHeader(i int) (int, error) {
	hp.HunkFile.Header.Start = i
	libNames, i, err := hp.readStrings(i)
	if err != nil {
		hp.Logger.Error("Reading library names", "error", err)
		return 0, nil
	}
	hp.HunkFile.Header.LibraryNames = libNames

	tableSize, i := hp.readUint32(i)
	hp.HunkFile.Header.TableSize = tableSize
	firstHunk, i := hp.readUint32(i)
	hp.HunkFile.Header.FirstHunk = firstHunk
	lastHunk, i := hp.readUint32(i)
	hp.HunkFile.Header.LastHunk = lastHunk

	numHunks := lastHunk - firstHunk
	hunkSizes := make([]uint32, 0, numHunks)
	for hn := 0; hn <= int(numHunks); hn++ {
		var hs uint32
		hs, i = hp.readUint32(i)
		// TODO: deal with mem flags
		hunkSizes = append(hunkSizes, hs)
	}
	hp.HunkFile.Header.HunkSizes = hunkSizes
	hp.HunkFile.Hunks = make([]Hunk, 0, numHunks)

	return i, nil
}

func (hp *HunkParser) readHunkCode(i int) (int, error) {
	start := i
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

	return i, nil
}

func (hp *HunkParser) readHunkReloc32(i int) (int, error) {
	start := i
	numOffsets, i := hp.readUint32(i)
	offsetNum, i := hp.readUint32(i)
	offsets := make([]uint32, 0, numOffsets)
	for oc := 0; oc < int(numOffsets); oc++ {
		var offset uint32
		offset, i = hp.readUint32(i)
		offsets = append(offsets, offset)
	}

	h := HunkReloc32{
		Start:      start,
		HunkType:   HUNK_RELOC32,
		NumOffsets: numOffsets,
		HunkNum:    offsetNum,
		Offsets:    offsets,
	}

	hp.HunkFile.Hunks = append(hp.HunkFile.Hunks, h)

	// TODO: repeat until HUNK_END is seen.
	return i, nil
}
