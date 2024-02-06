package main

func safeConvert(data interface{}, key string) string {
	if data == nil {
		return ""
	}
	if value, ok := data.(map[string]interface{})[key]; ok {
		return value.(string)
	}
	return ""
}
