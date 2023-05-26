package mjolnir

import (
	"runtime"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	"golang.org/x/exp/slices"

	"beckendrof/gaia/src/services/mjolnir/redis"
	"beckendrof/gaia/src/services/mjolnir/types"
	"beckendrof/gaia/src/services/mjolnir/utils"
	"beckendrof/gaia/src/services/mjolnir/zap"
)

type MjolnirSession struct{}

type MjolnirLogger struct {
	logData types.MjolnirMessage
}

type Services int

const (
	Gaia Services = iota
	Bifrost
	Asgard
	Genesis
	Malenia
	Animus
	Mjolnir
	Strato
)

var (
	lock           = &sync.Mutex{}
	wg             sync.WaitGroup
	mjolnirSession *MjolnirSession
)

func GetInstance() *MjolnirSession {
	if mjolnirSession == nil {
		lock.Lock()
		defer lock.Unlock()
		if mjolnirSession == nil {
			mjolnirSession = &MjolnirSession{}
		}
	}
	return mjolnirSession
}

func (m *MjolnirSession) MjolnirInit() bool {
	if configInit(utils.GetEnvKey("REDIS_HOST"), utils.GetEnvKey("REDIS_PORT")) {
		if initializeLogger() {
			return true
		}
	}
	return false
}

func (m *MjolnirSession) Log(service Services) *MjolnirLogger {
	servName := ""
	switch service {
	case 0:
		servName = "gaia"
	case 1:
		servName = "bifrost"
	case 2:
		servName = "asgard"
	case 3:
		servName = "genesis"
	case 4:
		servName = "malenia"
	case 5:
		servName = "animus"
	case 6:
		servName = "mjolnir"
	case 7:
		servName = "strato"
	}

	if slices.Contains(utils.GetServices(), servName) {
		mLog := &MjolnirLogger{}
		mLog.logData.Service = servName
		return mLog
	} else {
		return nil
	}
}

func (ml *MjolnirLogger) Debug(msg string) {
	_, file, line, _ := runtime.Caller(1)
	ml.logData.Level = "debug"
	ml.logData.Message = msg
	ml.logData.Caller = file + ":" + strconv.Itoa(line)
	redis.SendToRedisChannel(&ml.logData)
}

func (ml *MjolnirLogger) Info(msg string) {
	_, file, line, _ := runtime.Caller(1)
	ml.logData.Level = "info"
	ml.logData.Message = msg
	ml.logData.Caller = file + ":" + strconv.Itoa(line)
	redis.SendToRedisChannel(&ml.logData)
}

func (ml *MjolnirLogger) Warn(msg string) {
	_, file, line, _ := runtime.Caller(1)
	ml.logData.Level = "warn"
	ml.logData.Message = msg
	ml.logData.Caller = file + ":" + strconv.Itoa(line)
	redis.SendToRedisChannel(&ml.logData)
}

func (ml *MjolnirLogger) Error(msg string, trace string) {
	_, file, line, _ := runtime.Caller(1)
	ml.logData.Level = "error"
	ml.logData.Message = msg
	ml.logData.Caller = file + ":" + strconv.Itoa(line)
	ml.logData.StackTrace = trace
	redis.SendToRedisChannel(&ml.logData)
}

func (ml *MjolnirLogger) Fatal(msg string, trace string) {
	_, file, line, _ := runtime.Caller(1)
	ml.logData.Level = "fatal"
	ml.logData.Message = msg
	ml.logData.Caller = file + ":" + strconv.Itoa(line)
	ml.logData.StackTrace = trace
	redis.SendToRedisChannel(&ml.logData)
}

func (ml *MjolnirLogger) Panic(msg string, trace string) {
	_, file, line, _ := runtime.Caller(1)
	ml.logData.Level = "panic"
	ml.logData.Message = msg
	ml.logData.Caller = file + ":" + strconv.Itoa(line)
	ml.logData.StackTrace = trace
	redis.SendToRedisChannel(&ml.logData)
}

func configInit(host string, port string) bool {
	viper.AddConfigPath("src/services/mjolnir")
	viper.SetConfigName("mjolnir")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		return false
	}
	var config types.MjolnirConfig
	perr := viper.Unmarshal(&config)
	if perr != nil {
		return false
	}

	utils.Host = host
	utils.Port = port
	utils.Config = config
	return true
}

func initializeLogger() bool {
	if zap.InitializeLogger() {
		wg.Add(1)
		go redis.StartSub(&wg, utils.Config.RedisSubChannel)
		return true
	} else {
		return false
	}
}
