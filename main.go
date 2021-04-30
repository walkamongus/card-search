package main

import (
	"errors"
	goflag "flag"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/walkamongus/card-search/internal/hsapi"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// globalLog sets up global log settings for the ap
var globalLog = logf.Log.WithName("app")

// client is a pointer to a global API client available throughout the app
var client *hsapi.HSClient

// router represents the global Gin routing engine
var router *gin.Engine

func main() {
	// Inject zap logging flags
	opts := zap.Options{}
	opts.BindFlags(goflag.CommandLine)

	// Set logging backend to zap
	logger := zap.New(zap.UseFlagOptions(&opts))
	logf.SetLogger(logger)

	// Read in optional configuration file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			globalLog.Error(err, "Error parsing config file")
		}
	}

	// Enable configuration ENV variables
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// Set up supported flags
	flag.String("client-id", "", "API client ID")
	viper.BindPFlag("client-id", flag.Lookup("client-id"))
	flag.String("client-secret", "", "API client secret")
	viper.BindPFlag("client-secret", flag.Lookup("client-secret"))

	flag.Parse()

	if viper.GetString("client-id") == "" {
		globalLog.Error(errors.New("Required API client-id missing"), "Missing required config")
		os.Exit(2)
	}
	if viper.GetString("client-secret") == "" {
		globalLog.Error(errors.New("Required API client-secret missing"), "Missing required config")
		os.Exit(2)
	}

	// Set up Gin routes, templates, and initialize
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")
	initializeRoutes()
	router.Run()
}
