package env

import "github.com/caarlos0/env/v6"

type ConfigurationEnvironment struct {
	ApplicationEnvironment
	SystemInformation
	RedisEnvironment
}

type SystemInformation struct {
	SystemName string `env:"SYSTEM_NAME" envDefault:"-"`
	DomainName string `env:"DOMAIN_NAME" envDefault:"localhost:8080"`
}
type ApplicationEnvironment struct {
	AppEnv   string `env:"APP_ENV" envDefault:"development"`
	HttpPort string `env:"PORT" envDefault:"8080"`
}

type RedisEnvironment struct {
	RedisHost     string `env:"REDIS_HOST" envDefault:"0.0.0.0"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

var Conf = ConfigurationEnvironment{}

func LoadEnv() {
	if err := env.Parse(&Conf); err != nil {
		panic(err.Error())
	}
}
