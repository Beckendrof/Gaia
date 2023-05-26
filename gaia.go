/*
============================================================================================
This file contains the main function. It initializes the logger and the server and also
handles any port error or server error.
Thank You --- @abhinowP
============================================================================================
*/

package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"sync"

	apostolis "beckendrof/gaia/src/controllers/grpc/apostolis"
	apostolis_pb "beckendrof/gaia/src/services/grpc/apostolis"

	"beckendrof/gaia/src/services/metrics"
	utils "beckendrof/gaia/src/utils"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"beckendrof/gaia/src/services/mjolnir"
)

var (
	wg        sync.WaitGroup
	mainmutex sync.Mutex
	instance  *metrics.Metrics
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	session := mjolnir.GetInstance()
	flag := session.MjolnirInit()
	if flag {
		utils.GaiaLogger = session.Log(mjolnir.Gaia)
		if utils.GaiaLogger != nil {
			utils.GaiaLogger.Info("Gaia Logger Initialized")
		} else {
			log.Fatal("Get Logger Failed...")
		}
	}

	lis, err := net.Listen("tcp", os.Getenv("GAIA_HOST")+":"+os.Getenv("GAIA_PORT")) // server listening at port 8004
	if err != nil {                                                                  // exits program if any error is encountered
		utils.GaiaLogger.Panic("Apostolis Server Boot Failed", err.Error())
	}

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			metrics.CreateInstance("darwinAMD64")
		default:
			utils.GaiaLogger.Panic("Unsupported platform", runtime.GOARCH)
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			switch utils.GetNvidiaModel() {
			case "jetson-xavier":
				metrics.CreateInstance("xavier")
			case "jetson-orin":
				metrics.CreateInstance("orin")
			default:
				utils.GaiaLogger.Panic("Unsupported platform", utils.GetNvidiaModel())
			}
		default:
			utils.GaiaLogger.Panic("Unsupported platform", runtime.GOARCH)
		}
	default:
		utils.GaiaLogger.Panic("Unsupported OS", runtime.GOOS)
	}

	s := grpc.NewServer()                                                 // creates new grpc server
	apostolis_pb.RegisterApostolisServer(s, &apostolis.ApostolisServer{}) // registers the Apostolis service

	utils.GaiaLogger.Info("GRPC server listening at: " + lis.Addr().String())

	mainmutex.Lock()
	wg.Add(1)
	go s.Serve(lis)
	mainmutex.Unlock()
	wg.Wait()
}
