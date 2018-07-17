package gpipe

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

type Parse func(b []byte) (fmt.Stringer, error)

// Converter is parse func interface
type Converter interface {
	Convert(out fmt.Stringer, params []interface{}) (fmt.Stringer, error)
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
func Output(ctx context.Context, w io.Writer, in chan []byte, parse Parse, converters []Converter) error {
	for {
		select {
		case b, ok := <-in:
			if !ok {
				return nil
			}
			out, err := parse(b)
			if err != nil {
				break
			}
			for _, p := range converters {
				output, err := p.Convert(out, nil)
				if err != nil {
					return err
				}
				out = output
			}
			fmt.Fprintln(w, out.String())
		case <-ctx.Done():
			return nil
		}
	}
}
