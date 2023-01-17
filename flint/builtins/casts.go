package builtins

func castToStringSlice(v []any) ([]string, bool) {
	strs := make([]string, len(v))
	for i, vv := range v {
		str, ok := vv.(string)
		if !ok {
			return nil, false
		}

		strs[i] = str
	}

	return strs, true
}

func castToStringMap(v map[string]any) (map[string]string, bool) {
	strMap := make(map[string]string, len(v))
	for k, vv := range v {
		str, ok := vv.(string)
		if !ok {
			return nil, false
		}

		strMap[k] = str
	}

	return strMap, true
}

func castToMapSlice(v []any) ([]map[string]string, bool) {
	maps := make([]map[string]string, len(v))
	for i, vv := range v {
		m, ok := vv.(map[string]any)
		if !ok {
			return nil, false
		}

		strMap, ok := castToStringMap(m)
		if !ok {
			return nil, false
		}

		maps[i] = strMap
	}

	return maps, true
}
