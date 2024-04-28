package configs

type Config struct {
	AccessTokenSecretKey string         `yaml:"accessTokenSecretKey" json:"accessTokenSecretKey"`
	MongoConfig          MongoConfig    `yaml:"mongo" json:"mongoConfig"`
	RedisConfig          RedisConfig    `yaml:"redis" json:"redisConfig"`
	PostgresConfig       PostgresConfig `yaml:"postgres" json:"postgresConfig"`
}

type MongoConfig struct {
	URI      string ` yaml:"URI" json:"URI"`
	Database string `yaml:"database" json:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"address" json:"addr"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn" json:"dsn"`
}
