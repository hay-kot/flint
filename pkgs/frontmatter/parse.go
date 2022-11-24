package frontmatter

import (
	"regexp"
)

var (
	YAMLSeparator = []byte("---")
	NewLine       = []byte("\n")

	reKeyMatch = regexp.MustCompile(`(?m)^\s*\w+:`)
)

// type WriteOpts func(*writeOpts)

// type writeOpts struct {
// 	order []string
// }

// func WithOrder(order ...string) WriteOpts {
// 	return func(opts *writeOpts) {
// 		opts.order = order
// 	}
// }

// func Write(w io.Writer, m map[string]interface{}, optfuncs ...WriteOpts) (n int, err error) {
// 	opts := &writeOpts{}
// 	for _, optfunc := range optfuncs {
// 		optfunc(opts)
// 	}

// 	bits, err := yaml.Marshal(m)
// 	if err != nil {
// 		return n, err
// 	}

// 	if len(opts.order) > 0 {
// 		out := make([]byte, 0, len(bits))
// 		out = append(out, YAMLSeparator...)
// 		out = append(out, NewLine...)

// 		lines := bytes.Split(bits, NewLine)

// 		for _, key := range opts.order {
// 		inner:
// 			for i, line := range lines {
// 				if !bytes.HasPrefix(line, []byte(key+":")) {
// 					continue
// 				}

// 				out = append(out, line...)
// 				out = append(out, NewLine...)

// 				for _, l := range lines[i+1:] {
// 					if reKeyMatch.Match(l) { // break on next key
// 						break inner
// 					}

// 					if l == nil || bytes.Equal(l, NewLine) {
// 						continue
// 					}

// 					out = append(out, append(l, NewLine...)...)
// 				}
// 			}
// 		}

// 		out = append(out, YAMLSeparator...)
// 		out = append(out)

// 		return w.Write(out)
// 	}

// 	return w.Write(bits)
// }
