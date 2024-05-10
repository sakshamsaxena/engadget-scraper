package config

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

var configMap map[string]any

func init() {
	mode := os.Getenv("APP_ENV")
	configFile, openErr := os.OpenFile("config/"+mode+".json", os.O_RDONLY, 0444)
	if openErr != nil {
		panic(openErr)
	}
	configBytes, _ := io.ReadAll(configFile)
	unmarshalErr := json.Unmarshal(configBytes, &configMap)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
}

func Get(key string) any {
	return get(key, configMap)
}

func get(key string, m map[string]any) any {
	allKeys := strings.Split(key, ".")
	if len(allKeys) == 1 {
		return m[key]
	}
	return get(strings.Join(allKeys[1:], "."), m[allKeys[0]].(map[string]any))
}
