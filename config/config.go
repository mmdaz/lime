package config

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// Register sqlite3 driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Provider is a ...
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var (
	defaultConfig *viper.Viper

	// DB is a ...
	DB *gorm.DB
)

// Config is a ...
func Config() Provider {
	return defaultConfig
}

// LoadConfigProvider is a ...
func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	var err error
	defaultConfig = readViperConfig("LIME")

	DB, err = gorm.Open("sqlite3", defaultConfig.GetString("db_path"))
	if err != nil {
		fmt.Printf("Cannot connect to sqlite3 database")
		log.Fatal("This is the error:", err)
	}
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	v.SetDefault("mode", "debug") // release
	v.SetDefault("port", ":8080")

	v.SetDefault("admin_username", "admin")
	v.SetDefault("admin_password", "admin")

	v.SetDefault("cookie_secret", "TGq7dTjt@G.vkuDYwQfdf7uZvmwr@MzV.r2r6NGtPF")
	v.SetDefault("cookie_name", "console")

	v.SetDefault("db_path", "./gorm.db")

	return v
}
