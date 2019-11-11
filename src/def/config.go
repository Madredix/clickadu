package def

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type (
	Config struct {
		Environment string `mapstructure:"environment"`

		Log struct {
			Level  string `mapstructure:"level"`
			Format string `mapstructure:"format"`
			Output string `mapstructure:"output"`
		} `mapstructure:"log"`

		HTTP struct {
			Port int `mapstructure:"port"`
		} `mapstructure:"http"`

		baseDir string
	}
)

const (
	CfgDef    = "config"
	CfgType   = "json"
	importOpt = "imports"
)

func init() {
	Register(func(builder *Builder, params map[string]interface{}) error {

		// change current directory for tests
		dir, err := os.Getwd()
		if err != nil {
			return errors.New("can't get current dir")
		}
		if strings.Count(dir, `src`) == 2 {
			dir = dir[:strings.LastIndex(dir, `/src`)]
			if err = os.Chdir(dir); err != nil {
				return errors.New("can't change dir for test")
			}
		}

		var ok bool
		if _, ok = params["configFile"]; !ok {
			return errors.New("can't get required parameter config path")
		}

		var path string
		if path, ok = params["configFile"].(string); !ok {
			return errors.New(`parameter "configFile" should be string`)
		}

		return builder.Add(Definition{
			Name: CfgDef,
			Build: func(ctx Context) (interface{}, error) {
				var config Config
				var err error

				config.baseDir, err = filepath.Abs(filepath.Dir(path))
				if err != nil {
					panic(err.Error())
				}

				err = config.importFile(path, true)
				if err != nil {
					panic(err.Error())
				}

				return config, nil
			},
		})
	})
}

func (c *Config) importFile(path string, isMain bool) error {
	viperInst := viper.New()
	viperInst.SetConfigType(CfgType)

	replacer := strings.NewReplacer(".", "_")
	viperInst.SetEnvKeyReplacer(replacer)
	viperInst.SetEnvPrefix(`core`)
	viperInst.AutomaticEnv()

	var cfgPath = path
	var err error

	if !filepath.IsAbs(path) {
		if isMain {
			cfgPath, err = filepath.Abs(path)
			if err != nil {
				return err
			}
		} else {
			cfgPath, err = filepath.Abs(c.baseDir + string(filepath.Separator) + filepath.Clean(path))
			if err != nil {
				return err
			}
		}
	}
	viperInst.SetConfigFile(cfgPath)

	if err := viperInst.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: '%s'", err)
	}

	imports := viperInst.GetStringSlice(importOpt)
	for _, filePath := range imports {
		if err := c.importFile(filePath, false); err != nil {
			return err
		}
	}

	if err := viperInst.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
