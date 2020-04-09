package main

import "C"
import (
	"fmt"
	"mtconfig"
	"mtenv"
	"mtlog"
	"mtmanapi"
	"mtutils/utils"
	"os"
	"reportdb"
	"time"
)

type Config struct {
	MysqlUri   string `toml:"mysql_uri"`
	ServerAddr string `toml:"server_addr"`
	Login      int    `toml:"login"`
	Password   string `toml:"password"`
}

var (
	db            *reportdb.Db
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
	mtapi    mtmanapi.CManagerInterface
	reportDb *reportdb.Db
	//apiHelper mtmanapi.ManApiHelper
}

func (r *PumpReceiver) OnPump(code int, typ int, data, param uintptr) {
	switch {
	case code == mtmanapi.PUMP_START_PUMPING:
		fmt.Println(time.Now(), "START PUMPING")
	case code == mtmanapi.PUMP_UPDATE_TRADES:
		fmt.Println(time.Now(), "PUMP_UPDATE_TRADES")
		//r.db.AddTrade(&Trade{Ticket})

		total := 0
		trades := r.mtapi.TradesGet(&total)
		//fmt.Println("total: ", total)

		if total < 1 {
			return
		}

		t := mtmanapi.TradeRecordArray_getitem(trades, 0)

		//fmt.Println("ticket: ", t.GetOrder())
		//fmt.Println("pl: ", t.GetProfit())
		//fmt.Println("symbol: ", t.GetSymbol())
		if err := r.reportDb.AddTrade(reportdb.Trade{
			Ticket:     t.GetOrder(),
			Symbol:     t.GetSymbol(),
			Login:      t.GetLogin(),
			Cmd:        t.GetCmd(),
			Volume:     t.GetVolume(),
			OpenTime:   time.Unix(int64(t.GetOpen_time()), 0),
			OpenPrice:  t.GetOpen_price(),
			Sl:         t.GetSl(),
			Tp:         t.GetTp(),
			CloseTime:  time.Unix(int64(t.GetClose_time()), 0),
			ClosePrice: t.GetClose_price(),
			Profit:     t.GetProfit(),
			Magic:      t.GetMagic(),
			Comment:    t.GetComment(),
		}); err != nil {
			fmt.Println("Error while saving: ", err)
		}
		r.mtapi.MemFree(trades.Swigcptr())
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

	fmt.Println("Logining in with login/pwd: ", config.Login, "/", config.Password)
	errno = mtapi.Login(config.Login, config.Password)
	if errno != mtmanapi.RET_OK {
		panic(fmt.Sprintf("%d %s", errno, mtapi.ErrorDescription(errno)))
	}

	fmt.Println("Logged in")

	flags := mtmanapi.CLIENT_FLAGS_HIDENEWS |
		mtmanapi.CLIENT_FLAGS_HIDEMAIL |
		mtmanapi.CLIENT_FLAGS_HIDEONLINE

	fmt.Println("mysql connecting to ", config.MysqlUri)
	db, err = reportdb.NewDb(config.MysqlUri)
	checkError(err)
	defer db.Close()

	err = db.Ping()
	checkError(err)

	receiver.reportDb = db

	mtmanapi.PumpingSwitchEx(mtapi, flags, 0)

	fmt.Println("working...")
	utils.HandleSignal()
}
