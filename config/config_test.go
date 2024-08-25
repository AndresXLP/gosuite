package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/andresxlp/gosuite/config"
	"github.com/stretchr/testify/assert"
)

const (
	FILENAME       = ".%s.env"
	PROJECTDIRNAME = "config"
)

type App struct {
	Port            int           `mapstructure:"port"`
	ServiceName     string        `mapstructure:"service_name"`
	IntervalTimeOut time.Duration `mapstructure:"interval_time_out"`
	Required        bool          `mapstructure:"required" validate:"required"`
}

type Configuration struct {
	Host string `mapstructure:"host_dir" validate:"required"`
	App  App    `mapstructure:"app"`
}

type AdditionalConfig struct {
	TestVal int `mapstructure:"test_val"`
}

type MissingMapstructureTag struct {
	Test int
}

type MultipleEnvFiles struct {
	Host     string   `mapstructure:"host_dir" validate:"required"`
	App      App      `mapstructure:"app"`
	DataBase DataBase `mapstructure:"DB" validate:"required"`
}

type DataBase struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type InvalidMoldTags struct {
	Another int `mod:"whatever:9000" mapstructure:"another"`
}

type RequiredTags struct {
	OtherAnother  int `validate:"required" mapstructure:"other_another"`
	OtherAnother2 int `validate:"required" mapstructure:"other_another"`
	OtherAnother3 int `validate:"required" mapstructure:"other_another"`
}

func TestGetConfigFromEnv(t *testing.T) {
	t.Run("should get the config when there are no errors", func(t *testing.T) {
		_ = os.Setenv(kv.GetEnvKey("HOST_DIR"))
		_ = os.Setenv(kv.GetEnvKey("APP_PORT"))
		_ = os.Setenv(kv.GetEnvKey("APP_SERVICE_NAME"))
		_ = os.Setenv(kv.GetEnvKey("APP_INTERVAL_TIME_OUT"))
		_ = os.Setenv(kv.GetEnvKey("APP_REQUIRED"))
		defer func() {
			os.Unsetenv("HOST_DIR")
			os.Unsetenv("APP_PORT")
			os.Unsetenv("APP_SERVICE_NAME")
			os.Unsetenv("APP_INTERVAL_TIME_OUT")
			os.Unsetenv("APP_REQUIRED")
		}()
		cfg := Configuration{}
		if err := config.GetConfigFromEnv(&cfg); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
		}

		assert.Equal(t, kv.GetHostDir(), cfg.Host, fmt.Sprintf(hostDirShouldMsg, kv.GetHostDir(), cfg.Host))
		assert.Equal(t, kv.GetAppPort(), cfg.App.Port, fmt.Sprintf(appPortShouldMsg, kv.GetAppPort(), cfg.App.Port))
		assert.Equal(t, kv.GetAppServiceName(), cfg.App.ServiceName, fmt.Sprintf(appServiceNameShouldMsg, kv.GetAppServiceName(), cfg.App.ServiceName))
		assert.Equal(t, kv.GetAppIntervalTimeOut(), cfg.App.IntervalTimeOut, fmt.Sprintf(appIntervalTimeOutShouldMsg, kv.GetAppIntervalTimeOut(), cfg.App.ServiceName))
		assert.Equal(t, kv.GetAppRequired(), cfg.App.Required, fmt.Sprintf(appRequiredShouldMsg, kv.GetAppRequired(), cfg.App.Required))
	})

	t.Run("should return an error when binding wrong vars", func(t *testing.T) {
		_ = os.Setenv("TEST_VAL", "some")
		defer os.Unsetenv("TEST_VAL")
		cfg := AdditionalConfig{}
		if err := config.GetConfigFromEnv(&cfg); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
		}
	})

	t.Run("should return an error when invalid mold tag", func(t *testing.T) {
		cfg := InvalidMoldTags{}
		if err := config.GetConfigFromEnv(&cfg); err == nil {
			t.Errorf(shouldNotBeError, err)
		}
	})

	t.Run("should return an error when required tag", func(t *testing.T) {
		cfg := RequiredTags{}
		if err := config.GetConfigFromEnv(&cfg); err == nil {
			t.Errorf(shouldNotBeError, err)
		}
	})

	t.Run("should empty config when missing tag: 'mapstructure'", func(t *testing.T) {
		_ = os.Setenv("TEST", "123")
		defer os.Unsetenv("TEST")
		cfg := MissingMapstructureTag{}
		if err := config.GetConfigFromEnv(&cfg); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
			return
		}
		assert.Equal(t, 0, cfg.Test, fmt.Sprintf("TEST should be '0' and is %v", cfg.Test))
	})
}

func TestSetEnvsFromFile(t *testing.T) {
	t.Run("should get the config when there are no errors", func(t *testing.T) {
		cfg := Configuration{}

		if err := config.SetEnvsFromFile(PROJECTDIRNAME, fmt.Sprintf(FILENAME, "testing")); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
			return
		}
		defer func() {
			os.Unsetenv("HOST_DIR")
			os.Unsetenv("APP_PORT")
			os.Unsetenv("APP_SERVICE_NAME")
			os.Unsetenv("APP_INTERVAL_TIME_OUT")
			os.Unsetenv("APP_REQUIRED")
		}()
		if err := config.GetConfigFromEnv(&cfg); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
			return
		}

		assert.Equal(t, kv.GetHostDir(), cfg.Host, fmt.Sprintf(hostDirShouldMsg, kv.GetHostDir(), cfg.Host))
		assert.Equal(t, kv.GetAppPort(), cfg.App.Port, fmt.Sprintf(appPortShouldMsg, kv.GetAppPort(), cfg.App.Port))
		assert.Equal(t, kv.GetAppServiceName(), cfg.App.ServiceName, fmt.Sprintf(appServiceNameShouldMsg, kv.GetAppServiceName(), cfg.App.ServiceName))
		assert.Equal(t, kv.GetAppIntervalTimeOut(), cfg.App.IntervalTimeOut, fmt.Sprintf(appIntervalTimeOutShouldMsg, kv.GetAppIntervalTimeOut(), cfg.App.ServiceName))
		assert.Equal(t, kv.GetAppRequired(), cfg.App.Required, fmt.Sprintf(appRequiredShouldMsg, kv.GetAppRequired(), cfg.App.Required))
	})

	t.Run("should get the config when read multiple envs files", func(t *testing.T) {
		cfg := MultipleEnvFiles{}

		if err := config.SetEnvsFromFile(PROJECTDIRNAME, fmt.Sprintf(FILENAME, "testing"), fmt.Sprintf(FILENAME, "other")); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
			return
		}

		defer func() {
			os.Unsetenv("HOST_DIR")
			os.Unsetenv("APP_PORT")
			os.Unsetenv("APP_SERVICE_NAME")
			os.Unsetenv("APP_INTERVAL_TIME_OUT")
			os.Unsetenv("APP_REQUIRED")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_NAME")
		}()
		if err := config.GetConfigFromEnv(&cfg); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
			return
		}

		assert.Equal(t, kv.GetHostDir(), cfg.Host, fmt.Sprintf(hostDirShouldMsg, kv.GetHostDir(), cfg.Host))
		assert.Equal(t, kv.GetAppPort(), cfg.App.Port, fmt.Sprintf(appPortShouldMsg, kv.GetAppPort(), cfg.App.Port))
		assert.Equal(t, kv.GetAppServiceName(), cfg.App.ServiceName, fmt.Sprintf(appServiceNameShouldMsg, kv.GetAppServiceName(), cfg.App.ServiceName))
		assert.Equal(t, kv.GetAppIntervalTimeOut(), cfg.App.IntervalTimeOut, fmt.Sprintf(appIntervalTimeOutShouldMsg, kv.GetAppIntervalTimeOut(), cfg.App.ServiceName))
		assert.Equal(t, kv.GetAppRequired(), cfg.App.Required, fmt.Sprintf(appRequiredShouldMsg, kv.GetAppRequired(), cfg.App.Required))
		assert.Equal(t, kv.GetDBHost(), cfg.DataBase.Host, fmt.Sprintf(dbHostShouldMsg, kv.GetDBHost(), cfg.DataBase.Host))
		assert.Equal(t, kv.GetDBName(), cfg.DataBase.Name, fmt.Sprintf(dbNameShouldMsg, kv.GetDBName(), cfg.DataBase.Name))
		assert.Equal(t, kv.GetDBPort(), cfg.DataBase.Port, fmt.Sprintf(dbPortShouldMsg, kv.GetDBPort(), cfg.DataBase.Port))
	})

	t.Run("should return an error when missing required variable", func(t *testing.T) {
		cfg := Configuration{}
		if err := config.GetConfigFromEnv(&cfg); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
			return
		}
	})

	t.Run("should return an error opening file", func(t *testing.T) {
		if err := config.SetEnvsFromFile(PROJECTDIRNAME, fmt.Sprintf(FILENAME, "not")); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
			return
		}
	})
}
