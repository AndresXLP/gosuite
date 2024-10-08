package config

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	DefaultTag = "env"
)

var (
	validate = validator.New()
	conform  = modifiers.New()
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// GetConfigFromEnv reads the fields of the given struct and looks up environment variables
// that match their names, associating the values accordingly. The config argument is a pointer to the struct.
//
// By default, we use the DefaultTag, if you want to override this value
// you can do so using the PersonalTagName parameter.
func GetConfigFromEnv(config interface{}, PersonalTagName ...string) error {
	bindEnvs(config)

	if err := viper.Unmarshal(&config, func(config *mapstructure.DecoderConfig) {
		config.TagName = DefaultTag
		if PersonalTagName != nil {
			config.TagName = PersonalTagName[0]
		}
	}); err != nil {
		return err
	}

	if err := conform.Struct(context.Background(), config); err != nil {
		return err
	}

	if err := validate.Struct(config); err != nil {
		return err
	}

	return nil
}

// SetEnvsFromFile reads a file containing environment variables and adds them to the OS's environment.
//
// - projectDirName parameter specifies the folder name working directory.
//
// - fileNames parameter specifies the names of the files to search.
//
// Use only in development environments.
func SetEnvsFromFile(projectDirName string, fileNames ...string) error {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := string(re.Find([]byte(cwd)))

	filePaths := make([]string, len(fileNames))
	for i, fileName := range fileNames {
		filePaths[i] = filepath.Join(rootPath, fileName)
	}

	if err := godotenv.Load(filePaths...); err != nil {
		return err
	}

	return nil
}

func bindEnvs(config interface{}, parts ...string) {
	vc := reflect.ValueOf(config)

	if vc.Kind() == reflect.Ptr {
		vc = vc.Elem()
	}

	for i := 0; i < vc.NumField(); i++ {
		field := vc.Field(i)
		structField := vc.Type().Field(i)
		value, ok := structField.Tag.Lookup(DefaultTag)

		if !ok {
			continue
		}

		if field.Kind() == reflect.Struct {
			bindEnvs(field.Interface(), append(parts, value)...)
		} else {
			_ = viper.BindEnv(strings.Join(append(parts, value), "."))
		}
	}
}
