package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"config"
	"github.com/stretchr/testify/assert"
)

const fileName = ".%s.env"

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

func TestGetConfigFromEnv(t *testing.T) {
	t.Run("should get the config when there are no errors", func(t *testing.T) {
		_ = os.Setenv(kv.GetEnvKey("HOST_DIR"))
		_ = os.Setenv(kv.GetEnvKey("APP_PORT"))
		_ = os.Setenv(kv.GetEnvKey("APP_SERVICE_NAME"))
		_ = os.Setenv(kv.GetEnvKey("APP_INTERVAL_TIME_OUT"))
		_ = os.Setenv(kv.GetEnvKey("APP_REQUIRED"))
		defer os.Clearenv()
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
		defer os.Clearenv()
		cfg := AdditionalConfig{}
		if err := config.GetConfigFromEnv(&cfg); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
		}
	})

	t.Run("should empty config when missing tag: 'mapstructure'", func(t *testing.T) {
		_ = os.Setenv("TEST", "123")
		defer os.Clearenv()
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

		if err := config.SetEnvsFromFile(fmt.Sprintf(fileName, "testing")); err != nil {
			t.Errorf(fmt.Sprintf(shouldNotBeError, err))
			return
		}
		defer os.Clearenv()
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

	t.Run("should return an error when missing required variable", func(t *testing.T) {
		cfg := Configuration{}
		if err := config.GetConfigFromEnv(&cfg); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
			return
		}
	})

	t.Run("should return an error opening file", func(t *testing.T) {
		defer os.Clearenv()
		if err := config.SetEnvsFromFile(fmt.Sprintf(fileName, "not")); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
			return
		}
	})

	t.Run("should return an error wrong file content", func(t *testing.T) {
		defer os.Clearenv()
		if err := config.SetEnvsFromFile(fmt.Sprintf(fileName, "wrong")); err == nil {
			t.Errorf(fmt.Sprintf(shouldBeError, err))
			return
		}
	})
}
