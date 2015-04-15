// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"fmt"
	"strconv"
)

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
// TODO: implement
type Pattern struct {
	Version []byte
	Tempo   float32
	Tracks  []*Track
}

func (p *Pattern) String() string {
	var buf bytes.Buffer
	buf.WriteString("Saved with HW Version: " + string(p.Version[:]))
	buf.WriteByte('\n')
	buf.WriteString("Tempo: " + strconv.FormatFloat(float64(p.Tempo), 'f', -1, 32))
	buf.WriteByte('\n')
	for _, t := range p.Tracks {
		buf.WriteString(t.String())
	}

	return buf.String()
}

type Track struct {
	ID    int32
	Name  []byte
	Steps []bool
}

func (t Track) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("(%d) %s\t", t.ID, t.Name))

	count := 0
	for _, s := range t.Steps {
		if count%4 == 0 {
			buf.WriteByte('|')
		}
		if s {
			buf.WriteByte('x')
		} else {
			buf.WriteByte('-')
		}
		count++
	}

	buf.WriteByte('|')
	buf.WriteByte('\n')

	return buf.String()
}
