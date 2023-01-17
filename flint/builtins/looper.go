package builtins

import "github.com/hay-kot/flint/pkgs/frontmatter"

// looper defines a set of handler functions for data types that abstracts
// away the type switch/cast logic. It is used to iterate over a list of
// fields and call the appropriate handler function for each field.
//
// The following types are supported:
type looper struct {
	string    func(field string, data string, idx int)
	stringMap func(field string, data map[string]string, idx int)
}

func (l *looper) Do(fields []string, fm *frontmatter.FrontMatter) {
	for _, field := range fields {
		v, ok := fm.Get(field)
		if !ok {
			continue
		}

		switch v := v.(type) {
		// Supports the following types
		//
		// String
		// ---
		// key: value
		case string:
			if l.string == nil {
				continue
			}

			l.string(field, v, -1)

		// Supports the following list-like types
		//
		// String Slice
		// ---
		// key:
		//   - value1
		//   - value2
		//
		// Map Slice
		// ---
		// key:
		//   - key1: value1
		//     key2: value2
		//   - key1: value1
		//     key2: value2
		case []any:
			if l.string != nil {
				strs, ok := castToStringSlice(v)
				if !ok {
					continue
				}

				for i, str := range strs {
					l.string(field, str, i)
				}
			}

			if l.stringMap != nil {
				maps, ok := castToMapSlice(v)
				if !ok {
					continue
				}

				for i, m := range maps {
					l.stringMap(field, m, i)
				}
			}

		// Supports the following map-like types
		//
		// String Map
		// ---
		// key:
		//   key1: value1
		//   key2: value2
		case map[string]any:
			if l.stringMap == nil {
				continue
			}

			m, ok := castToStringMap(v)
			if !ok {
				continue
			}

			l.stringMap(field, m, -1)
		}

	}
}
