package config_test

import (
	"fmt"
	"time"
)

var kv = envKeyValue{value: map[string]interface{}{
	"HOST_DIR":              "0.0.0.0",
	"APP_PORT":              8080,
	"APP_SERVICE_NAME":      "service-name",
	"APP_INTERVAL_TIME_OUT": "30s",
	"APP_REQUIRED":          true,
	"DB_HOST":               "192.168.1.1",
	"DB_PORT":               5432,
	"DB_NAME":               "configs",
},
}

type envKeyValue struct {
	value map[string]interface{}
}

func (v *envKeyValue) GetEnvKey(key string) (string, string) {
	return key, fmt.Sprintf("%v", v.value[key])
}

func (v *envKeyValue) GetHostDir() string {
	return v.value["HOST_DIR"].(string)
}

func (v *envKeyValue) GetAppPort() int {
	return v.value["APP_PORT"].(int)
}

func (v *envKeyValue) GetAppServiceName() string {
	return v.value["APP_SERVICE_NAME"].(string)
}

func (v *envKeyValue) GetAppIntervalTimeOut() time.Duration {
	d, _ := time.ParseDuration(v.value["APP_INTERVAL_TIME_OUT"].(string))
	return d
}

func (v *envKeyValue) GetAppRequired() bool {
	return v.value["APP_REQUIRED"].(bool)
}

func (v *envKeyValue) GetDBHost() string {
	return v.value["DB_HOST"].(string)
}

func (v *envKeyValue) GetDBPort() int {
	return v.value["DB_PORT"].(int)
}

func (v *envKeyValue) GetDBName() string {
	return v.value["DB_NAME"].(string)
}
