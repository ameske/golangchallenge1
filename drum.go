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
type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []*Track
}

// String returns a string representation of a Pattern
func (p *Pattern) String() string {
	var buf bytes.Buffer

	buf.WriteString("Saved with HW Version: ")
	buf.WriteString(p.Version)
	buf.WriteByte('\n')
	buf.WriteString("Tempo: ")
	buf.WriteString(strconv.FormatFloat(float64(p.Tempo), 'f', -1, 32))
	buf.WriteByte('\n')

	for _, t := range p.Tracks {
		buf.WriteString(t.String())
	}

	return buf.String()
}

// Track is the high level representation of one musical
// instrument contained in a Pattern
type Track struct {
	ID    int32
	Name  string
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

	buf.WriteString("|\n")

	return buf.String()
}
