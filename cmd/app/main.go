package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/odit-bit/cloudfs/handler/app"
	"github.com/odit-bit/cloudfs/internal/blob/store/minioblob"
	"github.com/odit-bit/cloudfs/internal/blob/store/tokenpg"
	"github.com/odit-bit/cloudfs/internal/user/pguser"
	"github.com/odit-bit/cloudfs/service"
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
	viper.BindEnv("storage.user.uri", USER_DB_URI)
	viper.BindEnv("storage.token.uri", TOKEN_DB_URI)
	viper.BindEnv("http.host", HOST)
	viper.BindEnv("http.port", PORT)

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

	root.PersistentFlags().StringVar(&_host, "host", "", "host name")
	viper.BindPFlag("http.host", root.Flags().Lookup("host"))

	var conf config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	if err := conf.validate(); err != nil {
		log.Fatal(err)
	}

	appCMD := setupWebAppCmd(&conf)
	root.AddCommand(appCMD)
	if err := root.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}

func setupWebAppCmd(conf *config) *cobra.Command {
	cmd := cobra.Command{
		Use:  "run",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if port, err := cmd.Flags().GetString("port"); err != nil {
				log.Println(err)
				return
			} else if port != "" {
				conf.HTTP.Port = port
			}

			if host, err := cmd.Flags().GetString("host"); err != nil {
				log.Println(err)
				return
			} else if host != "" {
				conf.HTTP.Host = host
			}

			cmdCtx, cancel := context.WithTimeout(cmd.Context(), 3*time.Second)
			defer cancel()

			//setup user database
			userDB, err := pguser.NewDB(
				cmdCtx,
				conf.Storage.User.URI,
			)
			if err != nil {
				log.Println(err)
				return
			}
			defer userDB.Close()

			//setup blob storage
			blobStore, err := minioblob.New(
				conf.Storage.Blob.Endpoint,
				conf.Storage.Blob.AccessKey,
				conf.Storage.Blob.SecretKey,
			)
			if err != nil {
				log.Println(err)
				return
			}
			//setup blob token storage
			tokenStore, err := tokenpg.NewDB(cmdCtx, conf.Storage.Token.URI)
			if err != nil {
				log.Println(err)
				return
			}
			defer tokenStore.Close()

			//setup logger
			logger := slog.New(slog.NewTextHandler(cmd.OutOrStdout(), &slog.HandlerOptions{
				// AddSource: true,
				Level: slog.LevelDebug,
			}))

			// setup app
			api, err := service.NewCloudfs(tokenStore, blobStore, userDB)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			//setup session
			sess := scs.New()
			sess.Lifetime = conf.Session.Expired

			//setup application
			h := app.New(api, sess, logger)
			if err := h.Run(
				cmdCtx,
				conf.Address(),
				RateLimit(conf.HTTP.Limit, conf.HTTP.Burst),
				LogRequest(logger),
			); err != nil {
				panic(err)
			}

		},
	}

	return &cmd
}
