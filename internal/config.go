/*
Copyright 2022 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
	"errors"
	"fmt"
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

type configData struct {
	Mail     *string `yaml:"mail"`
	Password *string `yaml:"password"`
	Device   *string `yaml:"device"`
}

func getClient(debug bool) (jdownloader.JdClient, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	var logger *zap.Logger
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	return jdownloader.NewClient(*cfg.Mail, *cfg.Password, logger.Sugar(),
		jdownloader.ClientOptionTimeout(30*time.Second),
		jdownloader.ClientOptionAppKey("jdcli")), nil
}

func loadConfig() (*configData, error) {
	var cfg configData
	cfgDir, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(cfgDir)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if len(*cfg.Mail) == 0 || len(*cfg.Password) == 0 {
		return nil, errors.New("credentials are not specified. Use 'jdcli login' to populate them")
	}
	return &cfg, nil
}

func saveConfig(cfg *configData) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, 0o644)
}

func getConfigPath() (string, error) {
	cfgPath, ok := os.LookupEnv("JD_CONFIG")
	if !ok {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/jdconfig.yaml", cfgDir), nil
	}
	return cfgPath, nil
}
