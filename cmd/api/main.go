package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	var (
		_port        string
		_host        string
		_config_path string
	)

	viper.BindEnv("storage.blob.endpoint", BLOB_STORAGE_ENDPOINT)
	viper.BindEnv("storage.blob.accessKey", BLOB_STORAGE_ACCESS_KEY)
	viper.BindEnv("storage.blob.secretKey", BLOB_STORAGE_SECRET_KEY)
	viper.BindEnv("storage.user.uri", USER_DATABASE_URI)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigFile(wd + "/cloudfs.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	root := cobra.Command{}
	root.PersistentFlags().StringVarP(&_config_path, "config", "c", "", "path to config file")
	viper.BindPFlag("configPath", root.Flags().Lookup("config"))

	root.PersistentFlags().StringVarP(&_port, "port", "p", "8181", "port to listen-to")
	viper.BindPFlag("http.port", root.Flags().Lookup("port"))

	root.PersistentFlags().StringVar(&_host, "host", "localhost", "host name")
	viper.BindPFlag("http.host", root.Flags().Lookup("host"))

	var conf config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	if err := conf.validate(); err != nil {
		log.Fatal(err)
	}

	runCmd := setupRunCmd(&conf)
	root.AddCommand(runCmd)
	root.Execute()
}

func setupRunCmd(conf *config) *cobra.Command {
	cmd := cobra.Command{
		Use:  "run",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if port, err := cmd.Flags().GetString("port"); err != nil {
				log.Println(err)
			} else {
				conf.HTTP.PORT = port
			}

			if host, err := cmd.Flags().GetString("host"); err != nil {
				log.Println(err)
			} else {
				conf.HTTP.HOST = host

			}

			app := NewApiHandler(conf)
			httpServer(conf.Address(), app)

		},
	}

	return &cmd
}
