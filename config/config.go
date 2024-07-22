package config

import "github.com/spf13/viper"

type Config struct {
	PG     PostgresConfig
	Server Server
	Redis  RedisConfig
}

type PostgresConfig struct {
	DSN    string
	Schema string
}

type Server struct {
	ServerAddr string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func ParseEnv() *Config {
	conf := &Config{}
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// postgreSQL config
	conf.PG = PostgresConfig{
		DSN:    viper.GetString("PG_DSN"),
		Schema: viper.GetString("PG_SCHEMA"),
	}

	// fasthttpServer config
	conf.Server = Server{
		//Debug:      viper.GetBool("DEBUG"),
		ServerAddr: viper.GetString("SERVER_ADDR"),
	}

	// redis config
	conf.Redis = RedisConfig{
		Address:  viper.GetString("REDIS_ADDR"),
		Password: viper.GetString("REDIS_PSWD"), // "" - no password
		DB:       viper.GetInt("REDIS_DB"),      // 0 - default DB
	}

	return conf
}
