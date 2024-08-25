package deep_folder

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/andresxlp/gosuite/config"
)

const (
	fileName       = ".%s.env"
	projectDirName = "config"
)

type App struct {
	Port            int           `env:"port"`
	ServiceName     string        `env:"service_name"`
	IntervalTimeOut time.Duration `env:"interval_time_out"`
	Required        bool          `env:"required" validate:"required"`
}

type Configuration struct {
	Host string `env:"host_dir" validate:"required"`
	App  App    `env:"app"`
}

func TestSetEnvsFromFile(t *testing.T) {
	t.Run("should get the config when there are no errors in deepFolder", func(t *testing.T) {
		cfg := Configuration{}

		if err := config.SetEnvsFromFile(projectDirName, fmt.Sprintf(fileName, "testing")); err != nil {
			t.Errorf(fmt.Sprintf("It shouldn't be a error and it is %v", err))
			return
		}
		defer os.Clearenv()
		if err := config.GetConfigFromEnv(&cfg); err != nil {
			t.Errorf(fmt.Sprintf("It shouldn't be a error and it is %v", err))
			return
		}
	})
}
