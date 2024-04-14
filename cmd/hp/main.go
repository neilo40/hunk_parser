package main

// references: http://amiga-dev.wikidot.com/file-format:hunk

import (
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
	//Hunks       []Hunk // TODO: nest hunks, i.e. group relocs, debugs etc with the associated code/data hunk
	Hunks []HunkBlock
}

type Hunk struct {
	Blocks []HunkBlock
}

type HunkBlock interface {
	Print(l *slog.Logger)
}

const (
	HUNK_HEADER  = 0x000003F3
	HUNK_UNIT    = 0x000003E7 // TODO: support this
	HUNK_CODE    = 0x000003E9 // 1001
	HUNK_DATA    = 0x000003EA // 1002
	HUNK_BSS     = 0x000003EB // 1003
	HUNK_RELOC32 = 0x000003EC // 1004
	HUNK_END     = 0x000003F2 // 1010
)

var (
	errBadMagicCookie = errors.New("bad magic cookie")
	errHunkEndNotSeen = errors.New("hunk_end was not seen at the end of the current hunk")
)

func main() {

	hp := setup()

	i, err := hp.readHeader()
	if err != nil {
		panic(err)
	}

	//for range hp.HunkFile.Header.HunkSizes{
	for {
		if i >= len(hp.FileBytes) {
			break
		}
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
		case HUNK_DATA:
			i, err = hp.readHunkData(i)
			if err != nil {
				panic(err)
			}
		case HUNK_BSS: // no reloc blocks follow this
			i, err = hp.readHunkBss(i)
			if err != nil {
				panic(err)
			}
		default:
			slog.Error("Unknown hunk type", "type", hunkType, "byte", i)
			panic(true)
		}
	}

	hp.Logger.Info("Parsing Complete", "bytes_parsed", i, "total_bytes", len(hp.FileBytes))
}

func setup() *HunkParser {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, logOpts)
	logger := slog.New(handler)
	hp := &HunkParser{
		Logger:   logger,
		HunkFile: &HunkFile{},
	}

	err := hp.loadHunkFile("../../test/hello")
	if err != nil {
		panic(err)
	}

	return hp
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
