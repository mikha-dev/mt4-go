package main

import "C"
import (
	"encoding/json"
	"fmt"
	"mtconfig"
	"mtenv"
	"mtlog"
	"mtmanapi"
	"mtutils/utils"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type TickInfo struct {
	Time   int
	Symbol string
	Bid    float64
	Ask    float64
}

type Config struct {
	RedisAddr  string `toml:"redis_addr"`
	ServerAddr string `toml:"server_addr"`
	Login      int    `toml:"login"`
	Password   string `toml:"password"`
}

var (
	redisClient   *redis.Client
	factory       = mtmanapi.NewCManagerFactory()
	factoryInited = false
	apiVer        = makelong(
		mtmanapi.ManAPIProgramBuild,
		mtmanapi.ManAPIProgramVersion,
	)
)

func makelong(a, b int) int {
	return int(uint32(a) | uint32(b)<<16)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLog() {
	c := &mtconfig.Common{}

	err := mtconfig.LoadCommonConfig(mtenv.CONFIG_PATH, c)
	checkError(err)

	log, err := mtlog.NewLogger(c.LogPath, c.LogLevel)
	checkError(err)

	mtlog.SetDefault(log)

	mtlog.Info("root: \"%s\"", mtenv.DIR)
	mtlog.Info("config path: \"%s\"", mtenv.CONFIG_PATH)
	mtlog.Info("log path: \"%s\" with level \"%s\"", c.LogPath, c.LogLevel)

	fmt.Println("log path: ", c.LogPath)
}

func CreateManagerApi() (mtmanapi.CManagerInterface, error) {
	if !factoryInited {
		factory.Init()

		if factory.IsValid() == 0 {
			return nil, fmt.Errorf("Failed to load mtmanapi dll.")
		}

		if factory.WinsockStartup() != mtmanapi.RET_OK {
			return nil, fmt.Errorf("WinsockStartup failed")
		}
		factoryInited = true
	}

	mtapi := factory.Create(apiVer)
	if mtapi == nil {
		return nil, fmt.Errorf("Failed to create manager interface.")
	}

	return mtapi, nil
}

type PumpReceiver struct {
	mtapi mtmanapi.CManagerInterface
	redis *redis.Client
	//apiHelper mtmanapi.ManApiHelper
}

func (r *PumpReceiver) OnPump(code int, typ int, data, param uintptr) {
	// fmt.Println(time.Now(), "Pumping", code, typ)
	switch {
	case code == mtmanapi.PUMP_START_PUMPING:
		fmt.Println(time.Now(), "START PUMPING")
		// symbols := r.apiHelper.SymbolsGetAll()
		// for i := 0; i < int(symbols.Size()); i++ {
		// 	s := symbols.Get(i)
		// 	fmt.Println("Symbol: ", s.GetSymbol())
		// }
		// users := r.apiHelper.UsersGet()
		// for i := 0; i < int(users.Size()); i++ {
		// 	user := users.Get(i)
		// 	fmt.Println("User:", user.GetLogin(), user.GetName())
		// }
		// groups := r.apiHelper.GroupsGet()
		// for i := 0; i < int(groups.Size()); i++ {
		// 	group := groups.Get(i)
		// 	fmt.Println("Group:", group.GetGroup())
		// }
		// trades := r.apiHelper.TradesGet()
		// for i := 0; i < int(trades.Size()); i++ {
		// 	trade := trades.Get(i)
		// 	fmt.Println("Trade:", trade.GetOrder())
		// }
	case code == mtmanapi.PUMP_UPDATE_SYMBOLS:
		//fmt.Println(time.Now(), "PUMP_UPDATE_SYMBOLS")
	case code == mtmanapi.PUMP_UPDATE_BIDASK:
		//fmt.Println(time.Now(), "PUMP_UPDATE_BIDASK")

		nTotalTicksGot := 0
		pTicksInfo := r.mtapi.TickInfoLast("", &nTotalTicksGot)
		//fmt.Println("nTotalTicksGot: ", nTotalTicksGot)

		//tis := make([]*TickInfo, 0, nTotalTicksGot)

		for i := 0; i < nTotalTicksGot; i++ {
			ti := mtmanapi.TickInfoArray_getitem(pTicksInfo, i)
			if len(ti.GetSymbol()) > 1 {
				ts := &TickInfo{
					Symbol: ti.GetSymbol(),
					Bid:    ti.GetBid(),
					Ask:    ti.GetAsk(),
					Time:   ti.GetCtm(),
				}

				//		tis = append(tis, ts)
				jti, _ := json.Marshal(ts)
				fmt.Println(string(jti))
				err := r.redis.Publish("bid_ask:"+ts.Symbol, string(jti)).Err()
				if err != nil {
					fmt.Errorf("Failed to publish. ", err)
				}

			}
			//break
		}

		r.mtapi.MemFree(pTicksInfo.Swigcptr())

	case code == mtmanapi.PUMP_UPDATE_TRADES:
		//fmt.Println(time.Now(), "PUMP_UPDATE_TRADES")
		//trades := r.apiHelper.TradesGet()
		//for i := 0; i < int(trades.Size()); i++ {
		////	trade := trades.Get(i)
		//	fmt.Println("Update trade:", trade.GetOrder())
		//}
	default:
	}
}

func main() {

	initLog()

	config := &Config{}
	err := mtconfig.LoadConfig(mtenv.CONFIG_PATH, "config", config)
	checkError(err)

	receiver := &PumpReceiver{}
	mtmanapi.SetGlobalPumper(
		mtmanapi.NewDirectorPumpReceiver(receiver),
	)

	mtapi, err := CreateManagerApi()
	if err != nil {
		panic(err)
	}

	receiver.mtapi = mtapi
	//receiver.apiHelper = mtmanapi.NewManApiHelper(mtapi)
	fmt.Println("Connecting: ", config.ServerAddr)
	errno := mtapi.Connect(config.ServerAddr)
	if errno != mtmanapi.RET_OK {
		panic(fmt.Sprintf("%d %s", errno, mtapi.ErrorDescription(errno)))
	}

	fmt.Println("Connected")

	fmt.Println("Logining in login/pwd: ", config.Login, "/", config.Password)
	errno = mtapi.Login(config.Login, config.Password)
	if errno != mtmanapi.RET_OK {
		panic(fmt.Sprintf("%d %s", errno, mtapi.ErrorDescription(errno)))
	}

	fmt.Println("Loged in")

	flags := mtmanapi.CLIENT_FLAGS_HIDENEWS |
		mtmanapi.CLIENT_FLAGS_HIDEMAIL |
		mtmanapi.CLIENT_FLAGS_HIDEONLINE

	fmt.Println("redis connecting to ", config.RedisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	pong, err := redisClient.Ping().Result()
	checkError(err)
	fmt.Println("Redis PING - ", pong)

	receiver.redis = redisClient

	mtmanapi.PumpingSwitchEx(mtapi, flags, 0)

	//	amqpClient = NewAmqpClient(config.AmqpUrl, config.AmqpPrefix)
	//time.Sleep(time.Hour)
	fmt.Println("working...")
	utils.HandleSignal()

	//	fmt.Println("stopping...")
	//	mtapi.Disconnect()
	//	mtapi.Release()

	//	fmt.Println("stopped")
}
