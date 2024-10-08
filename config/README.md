# 🔐 Config Environments

This package provides convenient methods to manage environment configurations.

It builds upon the functionality of [Viper-go](https://github.com/spf13/viper), enhancing the handling of environment
variables in your application.

---

## Installation

```bash
  go get github.com/andresxlp/gosuite@latest
```

---

## Usage

### 🔐💻

#### Use it when your variables are set in your environment

```go
package main

import (
	"fmt"
	"os"

	"github.com/andresxlp/gosuite/config"
)

type Config struct {
	ServerName string `env:"server_name"`
	App        App    `env:"app"`
}

type App struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}

func main() {
	_ = os.Setenv("SERVER_NAME", "server-test")
	_ = os.Setenv("APP_HOST", "192.168.25.255")
	_ = os.Setenv("APP_PORT", "3001")

	cfg := Config{}

	if err := config.GetConfigFromEnv(&cfg); err != nil {
		fmt.Printf("error setting environment variables in struct \n Error: %v", err)
		return
	}

	fmt.Printf("The config with envs is %+v", cfg)
	// output: The config with envs is {ServerName:server-test App:{Host:192.168.25.255 Port:3001}}

}

```

---

### 📄➡️💻

#### Use it when your variables are in an .env file

```dotenv
SERVER_NAME=server-env
APP_HOST=127.0.0.1
APP_PORT=9000
```

---

```go
package main

import (
	"fmt"
	"os"

	"github.com/andresxlp/gosuite/config"
)

type Config struct {
	ServerName string `env:"server_name"`
	App        App    `env:"app"`
}

type App struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}

func main() {
	if err := config.SetEnvsFromFile("server",".develop.env"); err != nil {
		fmt.Printf("error setting variables in environment \n Error: %v", err)
		return
	}

	cfg := Config{}

	if err := config.GetConfigFromEnv(&cfg); err != nil {
		fmt.Printf("error setting environment variables in struct \n Error: %v", err)
		return
	}

	fmt.Printf("The config with envs is %+v", cfg)
	// output: The config with envs is {ServerName:server-env App:{Host:127.0.0.1 Port:9000}}
}
```

---

## Additional Info

### Nested Struct

For this type of `struct`, for example:
```go
type Config struct {
	ServerName string `env:"server_name"`
	App        App    `env:"app"`
}

type App struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}
```
The keys in the .env file must be configured as the nested structure as follows:
```dotenv
SERVER_NAME=sever-name
APP_HOST=0.0.0.0
APP_PORT=8080
```


### Tags
- `env`: `env:"server_name"` tag that allows us to bind our environment variable to the fields of the structure.
  - If you want to use a custom tag, you can use the `PersonalTagName` parameter in the `GetConfigFromEnv` function to change the default tag.
- `validate`: `validate:"required"` tag that allow use the [validator](https://github.com/go-playground/validator).
- `mod`: `mod:"default=true"` tag that allow use the [mold](https://github.com/go-playground/mold).

---

## Package Dependency

- 🐍 [Viper-go](https://github.com/spf13/viper)
- 🔐 [godotenv](https://github.com/joho/godotenv)
- 🔍 [validator-go](https://github.com/go-playground/validator)
- ✍️ [mold-go](https://github.com/go-playground/mold)

---

## Authors

- [@andresxlp](https://www.github.com/andresxlp)

---

## License

The project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/)

