package configs

type Config struct {
	MongoConfig    MongoConfig        `yaml:"mongo" json:"mongoConfig"`
	RedisConfig    RedisConfig        `yaml:"redis" json:"redisConfig"`
	PostgresConfig PostgresConfig     `yaml:"postgres" json:"postgresConfig"`
	Keys           Keys               `yaml:"keys" json:"keys"`
	Address        Address            `yaml:"address" json:"address"`
	Inbound        map[string]Address `yaml:"inbound" json:"inbound"`
	Outbound       map[string]Address `yaml:"outbound" json:"outbound"`
}

type Address struct {
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
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

type Keys struct {
	ApiGateway           string `yaml:"apiGateway" json:"apiGateway"`
	AccessTokenSecretKey string `yaml:"accessTokenSecretKey" json:"accessTokenSecretKey"`
}
