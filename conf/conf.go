package conf

import (
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"gopkg.in/ini.v1"
	"strings"
)

var Config = &conf{}

type conf struct {
	Path     string
	Data     map[string]interface{}
	IniObj   *ini.File
	ConfType ConfigType
}

type ConfigType string

var (
	Yaml ConfigType = "yaml"
	Ini  ConfigType = "ini"
)

// NewConf 默认yaml, 如果文件是ini 则读取ini文件
func NewConf(appConfigPath string) error {
	Config = &conf{
		Path: appConfigPath,
		Data: make(map[string]interface{}),
	}
	var err error
	if strings.Index(appConfigPath, ".yaml") != -1 {
		err = Config.InitYaml()
		Config.ConfType = Yaml
	}
	if strings.Index(appConfigPath, ".ini") != -1 {
		err = Config.InitIni()
		Config.ConfType = Ini
	}
	return err
}

func (c *conf) GetInt(key string) int {
	switch c.ConfType {
	case Yaml:
		if c.Data == nil {
			_ = c.InitYaml()
		}
		return utils.AnyToInt(c.Data[key])
	case Ini:
		value, err := c.IniObj.Section("").Key(key).Int()
		if err != nil {
			log.Error(err)
		}
		return value
	}
	return 0
}

func (c *conf) Get(key string) interface{} {
	switch c.ConfType {
	case Yaml:
		if c.Data == nil {
			_ = c.InitYaml()
		}
		return c.Data[key]
	case Ini:
		value := c.IniObj.Section("").Key(key).Value()
		return value
	}
	return nil
}

func (c *conf) GetStr(key string) string {
	switch c.ConfType {
	case Yaml:
		if c.Data == nil {
			_ = c.InitYaml()
		}
		return utils.AnyToString(c.Data[key])
	case Ini:
		value := c.IniObj.Section("").Key(key).Value()
		return value
	}
	return ""
}
