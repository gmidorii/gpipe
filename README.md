## gpipe

### Overview
pipeline implementation for go.

### Usage
```go
func main() {
	in := gpipe.Input(os.Stdin, 10)
	converters := []gpipe.Converter{
		// Converter type
  }
  parseFunc := func(b []byte) (fmt.Stringer, error) {
    // binary to input type
  }
  gpipe.Output(os.Stdout, in, parseFunc, converters)
}
```

Passed type from converter to convert is `fmt.Stringer`.  