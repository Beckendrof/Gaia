package redis

import (
	"context"
	"encoding/json"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"beckendrof/gaia/src/services/mjolnir/types"
	"beckendrof/gaia/src/services/mjolnir/utils"
)

var (
	workermutex sync.Mutex
	once        sync.Once
	instance    *redis.Client
	ctx         = context.Background()
)

func getRedisClient() *redis.Client {
	_, file, line, _ := runtime.Caller(1)
	once.Do(func() { // <-- atomic, does not allow repeating
		instance = redis.NewClient(&redis.Options{
			Addr:       utils.Host + ":" + utils.Port,
			Password:   "", // no password set
			DB:         0,  // use default DB,
			MaxRetries: 5,
		}) // <-- thread safe
		utils.LogToFile(&types.MjolnirMessage{Level: "info", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Establishing Redis Connection at: " + utils.Host + ":" + utils.Port, Caller: file + ":" + strconv.Itoa(line)})
	})
	return instance
}

func pingRedis() bool {
	_, file, line, _ := runtime.Caller(1)
	_, err := getRedisClient().Ping(ctx).Result()
	if err != nil {
		utils.LogToFile(&types.MjolnirMessage{Level: "panic", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Redis client ping request failed", Caller: file + ":" + strconv.Itoa(line), StackTrace: err.Error()})
		return false
	} else {
		return true
	}
}

func StartSub(wg *sync.WaitGroup, sub string) {
	_, file, line, _ := runtime.Caller(1)
	defer wg.Done()
	if pingRedis() {
		subscriber := instance.Subscribe(ctx, sub)
		for {
			workermutex.Lock()
			msg, recerr := subscriber.ReceiveMessage(ctx)
			if recerr != nil {
				SendToRedisChannel(&types.MjolnirMessage{Level: "error", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Redis Subscribe failed", Caller: file + ":" + strconv.Itoa(line), StackTrace: recerr.Error()})
			} else {
				mMsg, perr := utils.Parser(msg.Payload)
				if perr != nil {
					SendToRedisChannel(&types.MjolnirMessage{Level: "error", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Redis Msg Parse failed", Caller: file + ":" + strconv.Itoa(line), StackTrace: perr.Error()})
				} else {
					SendToRedisChannel(mMsg)
					workermutex.Unlock()
				}
			}
		}
	}
}

func SendToRedisChannel(msg *types.MjolnirMessage) {
	_, file, line, _ := runtime.Caller(1)
	payload, merr := json.Marshal(msg)
	if merr != nil {
		utils.LogToFile(&types.MjolnirMessage{Level: "error", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Msg parsing back to Log message model failed", Caller: file + ":" + strconv.Itoa(line), StackTrace: merr.Error()})
	} else {
		msg.TimeStamp = time.Now().UTC().Format(time.RFC3339)
		if utils.LogToFile(msg) {
			if pingRedis() {
				err := instance.Publish(ctx, "mjolnir", payload).Err()
				if err != nil {
					utils.LogToFile(&types.MjolnirMessage{Level: "panic", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Failed to publish the log to UI", Caller: file + ":" + strconv.Itoa(line), StackTrace: err.Error()})
				}
			}
		} else {
			utils.LogToFile(&types.MjolnirMessage{Level: "warn", Service: "mjolnir", TimeStamp: time.Now().UTC().Format(time.RFC3339), Message: "Service [" + msg.Service + "] / Level [" + msg.Level + "] received doesn't exists", Caller: file + ":" + strconv.Itoa(line)})
		}
	}
}
