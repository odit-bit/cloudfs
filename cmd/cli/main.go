package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	c := &cobra.Command{}
	c.AddCommand(&loginCMD, &registerCMD)
	if err := c.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

var (
	_default_dir = ""
)

func init() {
	var err error
	dir, err := getDefaultConfigPath()
	if err != nil {
		panic(err)
	}
	_default_dir = dir
}

var registerCMD = cobra.Command{
	Use:  "signup username password",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := loadConfig(_default_dir)
		if err != nil {
			return err
		}

		cli, err := New(conf.Endpoint)
		if err != nil {
			return err
		}
		if err := cli.Signup(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		log.Println("success")
		return nil
	},
}

var loginCMD = cobra.Command{
	Use: "login",
	RunE: func(cmd *cobra.Command, args []string) error {
		var username, password, endpoint string
		//bind flag to username, password
		username, _ = cmd.Flags().GetString("username")
		password, _ = cmd.Flags().GetString("password")
		endpoint, _ = cmd.Flags().GetString("address")

		conf, err := loadConfig(_default_dir)
		if err != nil {
			return err
		}

		if username == "" && password == "" {
			username = conf.AccessKey
			password = conf.SecretKey
		}

		if endpoint == "" {
			endpoint = conf.Endpoint
		}

		cli, err := New(conf.Endpoint)
		if err != nil {
			return err
		}
		token, err := cli.Login(cmd.Context(), endpoint, username, password)
		if err != nil {
			return err
		} else {
			conf.Token = token
		}

		conf.AccessKey = username
		conf.SecretKey = password
		conf.Endpoint = endpoint
		conf.Token = token

		return saveConfig(_default_dir, conf)

	},
}
