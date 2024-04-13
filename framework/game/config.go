package game

import (
	"common/logs"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
)

var Conf *Config

const (
	gameConfig = "gameConfig.json"
	servers    = "servers.json"
)

type Config struct {
	GameConfig  map[string]GameConfigValue `json:"gameConfig"`
	ServersConf ServersConf                `json:"serversConf"`
}
type ServersConf struct {
	Nats       NatsConfig         `json:"nats" `
	Connector  []*ConnectorConfig `json:"connector" `
	Servers    []*ServersConfig   `json:"servers" `
	TypeServer map[string][]*ServersConfig
}

type ServersConfig struct {
	ID               string `json:"id" `
	ServerType       string `json:"serverType" `
	HandleTimeOut    int    `json:"handleTimeOut" `
	RPCTimeOut       int    `json:"rpcTimeOut" `
	MaxRunRoutineNum int    `json:"maxRunRoutineNum" `
}

type ConnectorConfig struct {
	ID         string `json:"id" `
	Host       string `json:"host" `
	ClientPort int    `json:"clientPort" `
	Frontend   bool   `json:"frontend" `
	ServerType string `json:"serverType" `
}
type NatsConfig struct {
	Url string `json:"url" mapstructure:"db"`
}

type GameConfigValue map[string]any

func InitConfig(configDir string) {
	Conf = &Config{}
	dir, err := os.ReadDir(configDir)
	if err != nil {
		logs.Fatal("read config dir fail err:%v", err)
	}
	for _, file := range dir {
		if !file.IsDir() {
			configFile := path.Join(configDir, file.Name())
			if file.Name() == gameConfig {
				readGameConfig(configFile)
			}

			if file.Name() == servers {
				readServersConfig(configFile)
			}
		}
	}
}

func readServersConfig(file string) {
	conf := ServersConf{}

	v := viper.New()
	v.SetConfigFile(file)

	// 坚挺配置文件修改
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		// 解析
		err := v.Unmarshal(&conf)
		if err != nil {
			panic(fmt.Errorf("Fatal error unmarshal config: %s \n", err))
		}
		Conf.ServersConf = conf
	})
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 解析
	err = v.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Fatal error unmarshal config: %s \n", err))
	}

	Conf.ServersConf = conf
	typeServerConfig()
}

func typeServerConfig() {
	if len(Conf.ServersConf.Servers) > 0 {
		if Conf.ServersConf.TypeServer == nil {
			Conf.ServersConf.TypeServer = make(map[string][]*ServersConfig)
		}
		for _, v := range Conf.ServersConf.Servers {
			if Conf.ServersConf.TypeServer[v.ServerType] == nil {
				Conf.ServersConf.TypeServer[v.ServerType] = make([]*ServersConfig, 0)
			}
			Conf.ServersConf.TypeServer[v.ServerType] = append(Conf.ServersConf.TypeServer[v.ServerType], v)
		}
	}
}

func readGameConfig(configFile string) {
	//conf := make(map[string]GameConfigValue)
	//
	//v := viper.New()
	//v.SetConfigFile(file)
	//
	//// 坚挺配置文件修改
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//	// 解析
	//	err := v.Unmarshal(&conf)
	//	if err != nil {
	//		panic(fmt.Errorf("Fatal error unmarshal config: %s \n", err))
	//	}
	//	Conf.GameConfig = conf
	//})
	//err := v.ReadInConfig()
	//if err != nil {
	//	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	//}
	//
	//// 解析
	//err = v.Unmarshal(&conf)
	//if err != nil {
	//	panic(fmt.Errorf("Fatal error unmarshal config: %s \n", err))
	//}
	//
	//Conf.GameConfig = conf
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var gc map[string]GameConfigValue
	err = json.Unmarshal(data, &gc)
	if err != nil {
		panic(err)
	}
	Conf.GameConfig = gc
}

func (c *Config) GetConnector(serverId string) *ConnectorConfig {
	for _, v := range c.ServersConf.Connector {
		if v.ID == serverId {
			return v
		}
	}
	return nil
}

func (c *Config) GetConnectorByServerType(serverType string) *ConnectorConfig {
	for _, v := range c.ServersConf.Connector {
		if v.ServerType == serverType {
			return v
		}
	}
	return nil
}

func (c *Config) GetFrontGameConfig() map[string]any {
	result := make(map[string]any)
	for k, v := range c.GameConfig {
		value, ok := v["value"]
		backend := false
		_, exist := v["backend"]
		if exist {
			backend = v["backend"].(bool)
		}
		if ok && !backend {
			result[k] = value
		}
	}
	return result
}
