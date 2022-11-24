package builtins

import "github.com/hay-kot/flint/pkgs/set"

func extractKeys(mp map[string]any) *set.Set[string] {
	keys := set.New[string]()
	for k := range mp {
		v := mp[k]

		switch v := v.(type) {
		case map[string]any:
			for _, key := range extractKeys(v).Slice() {
				keys.Insert(k + "." + key)
			}
		default:
			keys.Insert(k)
		}
	}

	return keys
}

func extractValue(mp map[string]any, parts []string) (any, bool) {
	if len(parts) == 1 {
		v, ok := mp[parts[0]]
		if !ok {
			return "", false
		}

		return v, true
	}

	v, ok := mp[parts[0]]
	if !ok {
		return "", false
	}

	switch v := v.(type) {
	case map[string]any:
		return extractValue(v, parts[1:])
	default:
		return "", false
	}
}
