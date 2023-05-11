/*******************************************************************************
 * YAML config loader and CLI argument parser
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2023-05-11
 ******************************************************************************/

package main

import (
	"flag"
	"github.com/spf13/viper"
	"time"
)

var (
	createSnapshots *bool
)

func setDefaultConfig() {
	viper.SetDefault("loglevel", 1)
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 80)
	viper.SetDefault("server.base_path", "/")
	viper.SetDefault("storage.data", "/tmp/maxima-data")
	viper.SetDefault("storage.workspace", "/tmp")
	viper.SetDefault("job.command", "maxima")
	viper.SetDefault("job.timeout", 30*time.Second)
}

func loadConfig() error {
	setDefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	configPath := flag.String("config", "", "Path to config.yaml")
	createSnapshots = flag.Bool("create-snapshots", false, "Create snapshots; must run before normal application mode")
	flag.Parse()
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
	} else {
		viper.AddConfigPath("/etc/maxima-pool/")
		viper.AddConfigPath(".")
	}

	return viper.ReadInConfig()
}
