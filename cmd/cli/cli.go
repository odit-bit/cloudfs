package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/odit-bit/cloudfs/rpc"
	"google.golang.org/grpc"
)

const (
	_default_endpoint = "localhost:6969"
)

var defaultConfig = Config{
	Endpoint: _default_endpoint,
}

type Config struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Token     string
}

var (
	_dir  = ".cloudfs"
	_file = "config.json"
	// defaultConfigPath = ".cloudfs/config.json"
)

// return path to default config file path
func getDefaultConfigPath() (string, error) {
	ucd, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ucd, _dir, _file), nil
}

func createConfigFile(toDir string) error {

	if toDir == "" {
		s, err := getDefaultConfigPath()
		if err != nil {
			return err
		}
		f, err := os.Create(s)
		if err != nil {
			return err
		}
		return f.Close()
	}

	//check path
	confDir := filepath.Join(toDir, _dir)
	_, err := os.Stat(confDir)
	if err != nil {
		if os.IsNotExist(err) {
			//create dir
			if err := os.MkdirAll(toDir, 0o777); err != nil {
				return fmt.Errorf("failed create config dir: %v", err)
			}
		}
	}

	//create File
	if f, err := os.Create(filepath.Join(confDir, _file)); err != nil {
		return fmt.Errorf("failed create config file: %v", err)
	} else {
		c := Config{
			Endpoint: _default_endpoint,
		}
		return json.NewEncoder(f).Encode(c)
	}

}

func saveConfig(toDir string, conf *Config) error {
	if toDir == "" {
		dir, err := getDefaultConfigPath()
		if err != nil {
			return err
		}
		toDir = dir
	}

	f, err := os.Create(toDir)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(conf); err != nil {
		return err
	}

	return nil
}

func loadConfig(toDir string) (*Config, error) {
	if toDir == "" {
		s, err := getDefaultConfigPath()
		if err != nil {
			return nil, err
		}
		toDir = s
	} else {
		//check given path is exist
		info, err := os.Stat(toDir)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("path is directory not file")
		}
	}

	f, err := os.Open(filepath.Join(toDir, _dir, _file))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf Config
	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil

}

type cmd struct {
	cli *rpc.CloudfsClient
}

//----------------

func New(endpoint string) (*cmd, error) {

	conn, err := grpc.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	cli := rpc.NewCloudfsClient(conn)

	return &cmd{
		cli: cli,
	}, err
}

func (c *cmd) Signup(ctx context.Context, username, password string) error {
	_, err := c.cli.Register(ctx, rpc.RegisterParam{Username: username, Password: password})
	if err != nil {
		return err
	}
	return nil
}

func (c *cmd) Login(ctx context.Context, uri string, username, password string) (string, error) {
	var token string

	res, err := c.cli.BasicAuth(ctx, username, password)
	if err != nil {
		return token, err
	}

	token = res.Token
	return token, nil
}
