package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// Error when a pattern file's magic number is incorrect
var (
	ErrMagicNumberInvalid = errors.New("pattern: magic number invalid")
)

// StickyErrorReader remembers the first error encountered. Any Read()
// operations following the first error are effectively no-ops.
type StickyErrorReader struct {
	r   io.Reader
	err error
}

// Read implements the io.Reader interface on StickyErrorReader. If an error
// has already occurred it returns that error, otherwise it performs the read.
func (er *StickyErrorReader) Read(p []byte) (n int, err error) {
	if er.err != nil {
		return 0, er.err
	}

	n, err = er.r.Read(p)
	er.err = err

	return
}

// Err returns the first error the StickyErrorReader encountered, if any.
func (er *StickyErrorReader) Err() error {
	return er.err
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := &Pattern{}
	err = NewDecoder(f).Decode(p)

	return p, err
}

// A Decoder reads and decodes Pattern files from an input stream
type Decoder struct {
	r *StickyErrorReader
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: &StickyErrorReader{r: r},
	}
}

// Decode reads a Pattern file from its input and stores it in the value point to by p
func (dec *Decoder) Decode(p *Pattern) error {
	// Read the header information to verify the input is a pattern file and
	// learn its size
	header := struct {
		Magic    [6]byte
		FileSize int64
	}{}
	binary.Read(dec.r, binary.BigEndian, &header)
	if string(header.Magic[:]) != "SPLICE" {
		return ErrMagicNumberInvalid
	}

	// We'll now use the limited reader to ensure we only read the one Pattern file
	r := &io.LimitedReader{
		R: dec.r,
		N: header.FileSize,
	}

	// Version - 32 bytes
	tmpVersion := new([32]byte)
	io.ReadFull(r, tmpVersion[:])
	p.Version = string(bytes.Trim(tmpVersion[:], string(0x00)))

	// Tempo - 4 bytes, float32, BigEndian
	binary.Read(r, binary.LittleEndian, &p.Tempo)

	// Read the tracks associated with this Pattern
	p.Tracks = make([]*Track, 0)
	for r.N > 0 {
		t := &Track{}

		// Read the ID (4 bytes)
		binary.Read(r, binary.LittleEndian, &t.ID)

		// Read the Len of name (1 byte)
		nameLen := make([]byte, 1)
		io.ReadFull(r, nameLen)

		// Read the Name (see prev)
		tmpName := make([]byte, int(nameLen[0]))
		io.ReadFull(r, tmpName)
		t.Name = string(tmpName)

		// Read the steps (16 bytes)
		tmpSteps := make([]byte, 16)
		io.ReadFull(r, tmpSteps)

		// Convert the step into booleans, we don't care about bytes
		t.Steps = make([]bool, 16)
		for i := 0; i < 16; i++ {
			if tmpSteps[i] == 0x01 {
				t.Steps[i] = true
			} else {
				t.Steps[i] = false
			}
		}

		p.Tracks = append(p.Tracks, t)
	}

	return dec.r.Err()
}
