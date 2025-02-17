package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	AccessKey string
	SecretKey string
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
		c := Config{}
		return json.NewEncoder(f).Encode(c)
	}

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

//----------------

func login(ctx context.Context, uri string, username, password string) (string, error) {
	var body bytes.Buffer
	var token string
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, &body)
	if err != nil {
		return token, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	}
	switch res.StatusCode {
	case 200, 202:
		

	}
	return "", nil
}
