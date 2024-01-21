package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/integralist/go-findroot/find"
	"github.com/lockp111/go-easyzap"
	"github.com/spf13/viper"
)

const (
	ENV_PROFILE_LOCAL = "local"
)

var (
	runOnce sync.Once
	config  Config
)

type Config struct {
	AppName      string
	ServerPort   string
	HealthPort   string
	TimeLocation string
	DatabaseConfig
	AwsConfig
}

type DatabaseConfig struct {
	PostgresDriver string
	User           string
	Host           string
	Port           int
	Password       string
	DbName         string
	Conn
	DatabaseConnStr string
}

type Conn struct {
	Max int
}

type AwsConfig struct {
	Url           string
	Region        string
	QueueCampaing string
	Credentials   aws.AnonymousCredentials
}

func initConfig() {
	envProfile := os.Getenv("ENV_PROFILE")
	if envProfile == ENV_PROFILE_LOCAL {
		setEnvsByFile()
	}
	viper.AutomaticEnv()
}

func setEnvsByFile() {
	root, _ := find.Repo()

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(root.Path + "/build/package/env/local")

	if err := viper.ReadInConfig(); err != nil {
		easyzap.Panic(context.Background(), err, "failed reading config file")
	}
}

func GetConfig() Config {
	runOnce.Do(func() {
		initConfig()
		config = Config{
			AppName:      viper.GetString("APPLICATION_NAME"),
			ServerPort:   viper.GetString("SERVER_PORT"),
			HealthPort:   viper.GetString("HEALTH_PORT"),
			TimeLocation: viper.GetString("TIME_LOCATION"),
			DatabaseConfig: DatabaseConfig{
				PostgresDriver: viper.GetString("DATABASE_POSTGRESDRIVER"),
				User:           viper.GetString("DATABASE_USER"),
				Host:           viper.GetString("DATABASE_HOST"),
				Port:           viper.GetInt("DATABASE_PORT"),
				DbName:         viper.GetString("DATABASE_NAME"),
				Conn: Conn{
					Max: viper.GetInt("DATABASE_CON_MAX"),
				},
			},
			AwsConfig: AwsConfig{
				Url:           viper.GetString("AWS_URL"),
				Region:        viper.GetString("AWS_REGION"),
				QueueCampaing: viper.GetString("AWS_QUEUE_CAMPAING"),
			},
		}
		setEnvValues()
		config.DatabaseConfig.DatabaseConnStr = buildDatabaseConnString(config.DatabaseConfig)
	})
	return config
}

func setEnvValues() {
	if len(os.Getenv("DB_PASS_CAMPAING_CONSUMER_API")) > 0 {
		config.DatabaseConfig.Password = os.Getenv("DB_PASS_CAMPAING_CONSUMER_API")
	}
}

func buildDatabaseConnString(dbCfg DatabaseConfig) string {
	connectionDSN := fmt.Sprintf("user=%s host=%s port=%v  "+
		"password=%s dbname=%s connect_timeout=%v sslmode=disable",
		dbCfg.User,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Password,
		dbCfg.DbName,
		dbCfg.Conn.Max,
	)

	return connectionDSN
}
