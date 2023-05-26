package metrics

import (
	"sync"

	apostolis_pb "beckendrof/gaia/src/services/grpc/apostolis"
	"beckendrof/gaia/src/services/metrics/macos"
	"beckendrof/gaia/src/services/metrics/nvidia"
	"beckendrof/gaia/src/utils"
)

var lock = &sync.Mutex{}

type Metrics struct {
	macos.MacOS
	nvidia.Nvidia

	Platform string

	MemoryStats *apostolis_pb.MemoryReply
	CpuStats    *apostolis_pb.CPUReply
	GpuStats    *apostolis_pb.GPUReply
	DiskStats   *apostolis_pb.DiskReply
	NetStats    *apostolis_pb.NetReply
	LoadStats   *apostolis_pb.LoadReply
}

var Instance *Metrics

func CreateInstance(platform string) {
	if Instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if Instance == nil {
			Instance = &Metrics{}
			Instance.Platform = platform
		} else {
			utils.GaiaLogger.Info("Single instance already created.")
		}
	}
}

func (m *Metrics) GetAllStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.MemoryStats = m.AMD64.GetMemoryStats()
		m.CpuStats = m.AMD64.GetCPUStats()
		m.GpuStats = m.AMD64.GetGPUStats()
		m.DiskStats = m.AMD64.GetDiskStats()
		m.NetStats = m.AMD64.GetNetworkStats()
		m.LoadStats = m.AMD64.GetLoadAVGStats()

	case "xavier":
		m.MemoryStats = m.Xavier.GetMemoryStats()
		m.CpuStats = m.Xavier.GetCPUStats()
		m.GpuStats = m.Xavier.GetGPUStats()
		m.DiskStats = m.Xavier.GetDiskStats()
		m.NetStats = m.Xavier.GetNetworkStats()
		m.LoadStats = m.Xavier.GetLoadAVGStats()
	case "orin":
		m.MemoryStats = m.Orin.GetMemoryStats()
		m.CpuStats = m.Orin.GetCPUStats()
		m.GpuStats = m.Orin.GetGPUStats()
		m.DiskStats = m.Orin.GetDiskStats()
		m.NetStats = m.Orin.GetNetworkStats()
		m.LoadStats = m.Orin.GetLoadAVGStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}

func (m *Metrics) GetMemoryStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.MemoryStats = m.AMD64.GetMemoryStats()
	case "xavier":
		m.MemoryStats = m.Xavier.GetMemoryStats()
	case "orin":
		m.MemoryStats = m.Orin.GetMemoryStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}

func (m *Metrics) GetCPUStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.CpuStats = m.AMD64.GetCPUStats()
	case "xavier":
		m.CpuStats = m.Xavier.GetCPUStats()
	case "orin":
		m.CpuStats = m.Orin.GetCPUStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}

func (m *Metrics) GetLoadAVGStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.LoadStats = m.AMD64.GetLoadAVGStats()
	case "xavier":
		m.LoadStats = m.Xavier.GetLoadAVGStats()
	case "orin":
		m.LoadStats = m.Orin.GetLoadAVGStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}

func (m *Metrics) GetNetworkStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.NetStats = m.AMD64.GetNetworkStats()
	case "xavier":
		m.NetStats = m.Xavier.GetNetworkStats()
	case "orin":
		m.NetStats = m.Orin.GetNetworkStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}

func (m *Metrics) GetGPUStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.GpuStats = m.AMD64.GetGPUStats()
	case "xavier":
		m.GpuStats = m.Xavier.GetGPUStats()
	case "orin":
		m.GpuStats = m.Orin.GetGPUStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}

func (m *Metrics) GetDiskStats() {
	switch Instance.Platform {
	case "darwinAMD64":
		m.DiskStats = m.AMD64.GetDiskStats()
	case "xavier":
		m.DiskStats = m.Xavier.GetDiskStats()
	case "orin":
		m.DiskStats = m.Orin.GetDiskStats()
	default:
		utils.GaiaLogger.Panic("Unknown platform: ", Instance.Platform)
	}
}
