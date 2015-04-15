package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var MAGIC_NUM = []byte{0x53, 0x50, 0x4c, 0x49, 0x43, 0x45}

var ErrMagicNumberInvalid = errors.New("magic number invalid")

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

type Decoder struct {
	r   io.Reader
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:   r,
		err: nil,
	}
}

func (dec *Decoder) Decode(p *Pattern) error {
	if dec.err != nil {
		return dec.err
	}

	header := struct {
		Magic    [6]byte
		FileSize int64
	}{}
	err := binary.Read(dec.r, binary.BigEndian, &header)
	if string(header.Magic[:]) != "SPLICE" {
		return ErrMagicNumberInvalid
	}

	dec.r = &io.LimitedReader{
		R: dec.r,
		N: header.FileSize,
	}

	// Version - 32 bytes
	tmpVersion := new([32]byte)
	_, err = io.ReadFull(dec.r, tmpVersion[:])
	if err != nil {
		return err
	}
	p.Version = string(bytes.Trim(tmpVersion[:], string(0x00)))

	// Tempo - 4 bytes, float32, BigEndian
	err = binary.Read(dec.r, binary.LittleEndian, &p.Tempo)
	if err != nil {
		return err
	}

	// FOR EACH TRACK
	p.Tracks = make([]*Track, 0)
	for {
		t := &Track{}
		// Read the ID (4 bytes)
		err = binary.Read(dec.r, binary.LittleEndian, &t.ID)
		if err != nil {
			break
		}

		// Read the Len of name (1 byte)
		nameLen := make([]byte, 1)
		_, err = io.ReadFull(dec.r, nameLen)
		if err != nil {
			break
		}

		// Read the Name (see prev)
		tmpName := make([]byte, int(nameLen[0]))
		_, err = io.ReadFull(dec.r, tmpName)
		if err != nil {
			break
		}
		t.Name = string(tmpName)

		// Read the steps (16 bytes)
		tmpSteps := make([]byte, 16)
		_, err = io.ReadFull(dec.r, tmpSteps)
		if err != nil {
			break
		}

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

	if err == io.EOF {
		return nil
	} else {
		dec.err = err
		return err
	}
}
