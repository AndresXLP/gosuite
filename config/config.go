package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	validate = validator.New()
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// GetConfigFromEnv reads the fields of the given struct and looks up environment variables
// that match their names, associating the values accordingly. The config argument is a pointer to the struct.
func GetConfigFromEnv(config interface{}) error {
	bindEnvs(config)

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	if err := validate.Struct(config); err != nil {
		var errs []string
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErr {
				errs = append(errs, fmt.Sprintf(fieldError.Error()))
			}
			return fmt.Errorf(strings.Join(errs, ";\n"))
		}

		if invalidValidationErr, ok := err.(*validator.InvalidValidationError); ok {
			return invalidValidationErr
		}

		return err
	}

	return nil
}

// SetEnvsFromFile reads a file that contains environment variables in the format: 'VAR_NAME=var-value'
// and sets these variables in the OS's environment.
//
// The fileName parameter specifies the name of the file to search for in the working directory.
func SetEnvsFromFile(fileNames ...string) error {
	if err := godotenv.Load(fileNames...); err != nil {
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
		value, ok := structField.Tag.Lookup("mapstructure")

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
