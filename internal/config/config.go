/*
Инициализация конфигурации
*/
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	/*
		теги:
			yaml:"env" - сопоставляет названию в файле config.yml
			env-default:"local"` - устанавливает значение по умолчанию в случве отсутствия
			env-required:"true" - приложение не запустится если парамт отсутствует в config.yml
	*/
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	AliasLength int    `yaml:"alias_length" env-default:"6"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

/*
Загружает конфигурацию в cleanenv
Приставка Must означает что фукция будет ПАНИКОВАТЬ, а не возвращать ошибку
*/
func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//	check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", err)

	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
