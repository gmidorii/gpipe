package gpipe

import (
	"bufio"
	"fmt"
	"io"
)

// Parser is parse func interface
type Parser interface {
	//Parse(param []interface{}) (interface{}, error)
	Parse(in interface{}) (interface{}, error)
}

// Input is createing input reader
// return input chank to response channel
func Input(r io.Reader, size int) chan []byte {
	in := make(chan []byte, size)
	go func(r io.Reader, in chan []byte) {
		defer close(in)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			in <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			fmt.Errorf("error: %v", err)
		}
	}(r, in)
	return in
}

// Output is output to io.Writer
func Output(w io.Writer, in chan []byte, conv func(b []byte) interface{}, parsers []Parser) error {
	for {
		select {
		case b, ok := <-in:
			if !ok {
				break
			}
			out := conv(b)
			for _, p := range parsers {
				output, err := p.Parse(out)
				if err != nil {
					return err
				}
				out = output
			}
			w.Write(out.([]byte))
		}
	}
	return nil
}
