package app


func getMapString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getMapInt(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return defaultValue
}

func getMapBool(m map[string]interface{}, key string, defaultVal bool) bool {
    if val, ok := m[key]; ok {
        if b, ok := val.(bool); ok {
            return b
        }
    }
    return defaultVal
}