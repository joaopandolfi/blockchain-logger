package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	c "github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
)

// Configs Struct
type Config struct {
	File            map[string]string
	AESKey          string
	JWTSecret       string
	DefaultPassword string
	BcryptCost      int
	SystemID        string
	Propertyes      c.Configurations
	BlockChain      blockChain
	PostgreSQL      string
	Server          server `json:"server"`
	SnakeByDefault  bool
}

type server struct {
	Port         string
	Host         string
	TimeoutWrite time.Duration
	TimeoutRead  time.Duration
	Debug        bool
	Security     security
}

type security struct {
	TLSCert    string
	TLSKey     string
	Opsec      secure.Options
	BcryptCost int //10,11,12,13,14
	JWTSecret  string
	AESKey     string
}

type blockChain struct {
	PrivKey    string
	PubKey     string
	Passphrase string
}

// Config global
var cfg *Config

// Get Config
func Get() Config {
	if cfg == nil {
		panic(fmt.Errorf("config not loaded"))
	}

	return *cfg
}

func (c *Config) getEnvOrFile(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return cfg.File[key]
}

// Load config
func Load(args []string) {
	godotenv.Load(".env")
	cfile := "./config.json"
	cfg = &Config{
		File: c.LoadJsonFile(cfile),
	}
	c.LoadConfig(c.LoadFromFile(cfile))
	cfg.Propertyes = c.Configuration
	cfg.SystemID = cfg.getEnvOrFile("SYSTEM_ID")

	cfg.AESKey = cfg.getEnvOrFile("AES_KEY")
	cfg.Server.Security.AESKey = cfg.AESKey
	cfg.JWTSecret = cfg.getEnvOrFile("JWT_SECRET")
	cfg.SnakeByDefault, _ = strconv.ParseBool(cfg.getEnvOrFile("SNAKE_DEFAULT"))

	cfg.PostgreSQL =
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.getEnvOrFile("POSTGRESQL_HOST"), cfg.getEnvOrFile("POSTGRESQL_PORT"), cfg.getEnvOrFile("POSTGRESQL_USER"),
			cfg.getEnvOrFile("POSTGRESQL_PASSWORD"), cfg.getEnvOrFile("POSTGRESQL_DB"))

	cfg.BcryptCost, _ = strconv.Atoi(cfg.getEnvOrFile("BCRYPT_COST"))
	cfg.DefaultPassword = cfg.getEnvOrFile("DEFAULT_PASSWORD")

	// slack
	c.Configuration.SlackChannel = cfg.getEnvOrFile("SLACK_CHANNEL")
	c.Configuration.SlackToken = cfg.getEnvOrFile("SLACK_TOKEN")
	c.Configuration.SlackWebHook = []string{cfg.getEnvOrFile("SLACK_WEBHOOK")}
	c.Configuration.Security.JWTSecret = cfg.JWTSecret

	cfg.BlockChain.PubKey = cfg.getEnvOrFile("BLOCKCHAIN_PUB_KEY")
	cfg.BlockChain.PrivKey = cfg.getEnvOrFile("BLOCKCHAIN_PRIV_KEY")
	cfg.BlockChain.Passphrase = cfg.getEnvOrFile("BLOCKCHAIN_PRIV_KEY_PASS")

	// Load And Inject Jaeger Envs
	os.Setenv("JAEGER_SERVICE_NAME", fmt.Sprintf("%s%s", cfg.SystemID, cfg.getEnvOrFile("JAEGER_ENVIRONMENT")))
	os.Setenv("JAEGER_AGENT_HOST", cfg.getEnvOrFile("JAEGER_AGENT_HOST"))
	os.Setenv("JAEGER_SAMPLER_TYPE", cfg.getEnvOrFile("JAEGER_SAMPLER_TYPE"))
	os.Setenv("JAEGER_SAMPLER_PARAM", cfg.getEnvOrFile("JAEGER_SAMPLER_PARAM"))
	os.Setenv("JAEGER_REPORTER_LOG_SPANS", cfg.getEnvOrFile("JAEGER_REPORTER_LOG_SPANS"))
}

func Inject(c *Config) {
	cfg = c
}
