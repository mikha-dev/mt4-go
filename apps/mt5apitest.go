package main

import "C"
import (
	"fmt"
	"mt5manapi"
	"mtconfig"
	"mtenv"
	"mtlog"
	"mtutils/utils"
	"os"

	"github.com/go-redis/redis"
)

type MT5Config struct {
	RedisAddr  string `toml:"redis_addr"`
	ServerAddr string `toml:"server_addr"`
	Login      uint64 `toml:"login"`
	Password   string `toml:"password"`
}

var (
	redisCli         *redis.Client
	mt5factory       = mt5manapi.NewCMTManagerAPIFactory()
	mt5factoryInited = false
)

func initLogger() {
	c := &mtconfig.Common{}

	err := mtconfig.LoadCommonConfig(mtenv.CONFIG_PATH, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log, err := mtlog.NewLogger(c.LogPath, c.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mtlog.SetDefault(log)

	mtlog.Info("root: \"%s\"", mtenv.DIR)
	mtlog.Info("config path: \"%s\"", mtenv.CONFIG_PATH)
	mtlog.Info("log path: \"%s\" with level \"%s\"", c.LogPath, c.LogLevel)

	fmt.Println("log path: ", c.LogPath)
}

func CreateIManagerApi() (mt5manapi.IMTManagerAPI, error) {

	if !mt5factoryInited {

		err := mt5factory.Initialize()

		if err != 0 {
			return nil, fmt.Errorf("Failed to load MT5APIManager dll. err: ", err)
		}

		var version uint
		errno := mt5factory.Version(&version)
		if errno != 0 {
			return nil, fmt.Errorf("Failed to load version.")
		}

		if int(version) < mt5manapi.MTManagerAPIVersion {
			return nil, fmt.Errorf("Wrong Manager API version %u, version %u required\n", version, mt5manapi.MTManagerAPIVersion)
		}

		mt5factoryInited = true
	}

	mtapi := mt5factory.CreateManager(uint(mt5manapi.MTManagerAPIVersion))

	if mtapi == nil {
		return nil, fmt.Errorf("Failed to Create manager")
	}

	return mtapi, nil
}

func main() {

	initLogger()

	config := &MT5Config{}
	err := mtconfig.LoadConfig(mtenv.CONFIG_PATH, "mt5config", config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	man, err := CreateIManagerApi()
	if err != nil {
		panic(err)
	}

	errno := man.Connect(config.ServerAddr, config.Login, config.Password, "", 0, 10)
	if errno != 0 {
		panic(errno)
	}

	/*var server_ref string
	if man.NetworkServer(server_ref) == 0 && man.NetworkRescan(0, 10000) == 0 {

		man.Disconnect()

		errno := man.Connect(server_ref, config.ServerAddr, config.Password, "", 0xffffffff, 10)
		if errno != 0 {
			panic(errno)
		}
	} */

	var ticks_manager mt5manapi.IMTTickSink
	man.TickSubscribe(ticks_manager)
	man.SelectedAddAll()

	fmt.Println("redis connecting to ", config.RedisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//checkError(err)
	fmt.Println("Redis PING - ", pong)

	fmt.Println("working...")
	utils.HandleSignal()
}
