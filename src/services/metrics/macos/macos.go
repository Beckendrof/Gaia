//go:build darwin
// +build darwin

package macos

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/dkorunic/iSMC/smc"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/maxbeatty/golang-book/chapter11/math"
	CPU "github.com/shirou/gopsutil/cpu"
	Disk "github.com/shirou/gopsutil/v3/disk"

	apostolis_pb "beckendrof/gaia/src/services/grpc/apostolis"
	"beckendrof/gaia/src/utils"
)

type MacOS struct {
	AMD64 AMD64Controller
	// ARM   ARMController
}

type AMD64Controller struct {
	apostolis_pb.MemoryReply
	apostolis_pb.CPUReply
	apostolis_pb.GPUReply
	apostolis_pb.DiskReply
	apostolis_pb.NetReply
	apostolis_pb.LoadReply
}

var (
	rxBytesIdx, txBytesIdx, Pagesize                      int
	Speculative, Wired, FileBacked, Compressed, Purgeable float64
)

func (m *AMD64Controller) GetCPUStats() *apostolis_pb.CPUReply {
	cpu_stats, err := cpu.Get()
	m.CPUReply.User = utils.ToFixed(float64(cpu_stats.User), 2)
	m.CPUReply.System = utils.ToFixed(float64(cpu_stats.System), 2)
	m.CPUReply.Idle = utils.ToFixed(float64(cpu_stats.Idle), 2)
	m.CPUReply.Total = utils.ToFixed(float64(cpu_stats.Total), 2)

	// Temperature
	temp_stats := smc.GetTemperature()
	var temps []float64
	for k := range temp_stats {
		if strings.Contains(k, "CPU Core") {
			value, err := strconv.ParseFloat(temp_stats[k].(map[string]interface{})["value"].(string)[0:4], 64)
			if err != nil {
				utils.GaiaLogger.Error("Error parsing cpu temp: ", err.Error())
			}
			temps = append(temps, utils.ToFixed(value, 2))
		}
	}
	m.CPUReply.Temp = utils.ToFixed(math.Average(temps), 2)
	// m.CPUReply.CoreTemps = temps

	// Power
	power := smc.GetPower()
	for k := range power {
		if strings.Contains(k, "CPU Package Total") {
			cpu_power, err := strconv.ParseFloat(power[k].(map[string]interface{})["value"].(string)[0:3], 64)
			if err != nil {
				utils.GaiaLogger.Error("Error parsing cpu power: ", err.Error())
			}
			m.CPUReply.Power = utils.ToFixed(cpu_power, 2)
		}
	}

	// Utilization per core
	cores, err := CPU.Percent(0, true)
	if err != nil {
		utils.GaiaLogger.Error("Error getting cpu percent: ", err.Error())
	}

	m.CPUReply.Cpus = nil
	for _, core := range cores {
		m.CPUReply.Cpus = append(m.CPUReply.Cpus, utils.ToFixed(core, 2))
	}

	m.CPUReply.User = utils.ToFixed(m.CPUReply.User/m.CPUReply.Total*100, 2)     // User CPU usage
	m.CPUReply.System = utils.ToFixed(m.CPUReply.System/m.CPUReply.Total*100, 2) // System CPU usage
	m.CPUReply.Idle = utils.ToFixed(m.CPUReply.Idle/m.CPUReply.Total*100, 2)     // Idle CPU usage
	return &m.CPUReply
}

func (m *AMD64Controller) GetGPUStats() *apostolis_pb.GPUReply {
	cmd := exec.Command("powermetrics", "--samplers", "smc,gpu_power", "-n", "1")
	out, _, _ := utils.ReturnExitCode(cmd)
	gpu_stats := strings.Split(out, "**** GPU usage ****\n\n")[1]

	tmp := strings.FieldsFunc(strings.Split(strings.Split(gpu_stats, "average active frequency as fraction of nominal")[1], "\n")[0], utils.Split)
	gpu_tf, err := strconv.ParseFloat(tmp[1][:5], 64)
	if err != nil {
		utils.GaiaLogger.Error("Error parsing gpu total frequency: ", err.Error())
	}
	m.GPUReply.Total = utils.ToFixed(gpu_tf, 2)

	m.GPUReply.Percent, err = strconv.ParseFloat(tmp[2][1:len(tmp[2])-2], 64)
	if err != nil {
		utils.GaiaLogger.Error("Error parsing gpu percentage used: ", err.Error())
	}
	m.GPUReply.Percent = utils.ToFixed(m.GPUReply.Percent, 2)

	m.GPUReply.Used, err = strconv.ParseFloat(tmp[3][:len(tmp[3])-3], 64)
	if err != nil {
		utils.GaiaLogger.Error("Error parsing gpu average active frequency: ", err.Error())
	}

	// Temperature
	temp_stats := smc.GetTemperature()
	m.GPUReply.Temp = nil
	m.GPUReply.DeviceName = nil
	for k := range temp_stats {
		if strings.Contains(k, "GPU") {
			if !strings.Contains(k, "GPU Proximity") {
				value, err := strconv.ParseFloat(temp_stats[k].(map[string]interface{})["value"].(string)[0:4], 64)
				if err != nil {
					utils.GaiaLogger.Error("Error parsing gpu temperature: ", err.Error())
				}
				m.GPUReply.Temp = append(m.GPUReply.Temp, value)
				m.GPUReply.DeviceName = append(m.GPUReply.DeviceName, strings.Split(k, "GPU ")[1])
			}
		}
	}

	// Power
	power_stats := smc.GetPower()
	m.GPUReply.Power = nil
	for k := range power_stats {
		if strings.Contains(k, "GPU") {
			value, err := strconv.ParseFloat(power_stats[k].(map[string]interface{})["value"].(string)[0:len(power_stats[k].(map[string]interface{})["value"].(string))-2], 64)
			if err != nil {
				utils.GaiaLogger.Error("Error parsing gpu power: ", err.Error())
			}
			m.GPUReply.Power = append(m.GPUReply.Power, value)
		}
	}

	return &m.GPUReply
}

func (m *AMD64Controller) GetMemoryStats() *apostolis_pb.MemoryReply {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "vm_stat")
	out, err := cmd.StdoutPipe()
	if err != nil {
		utils.GaiaLogger.Panic("Error getting memory stats: ", "")
	}
	if err := cmd.Start(); err != nil {
		utils.GaiaLogger.Panic("Error getting memory stats: ", "")
	}
	scanner := bufio.NewScanner(out)
	if !scanner.Scan() {
		utils.GaiaLogger.Panic("failed to scan output of vm_stat", "")
	}
	line := scanner.Text()
	if _, err := fmt.Sscanf(line, "Mach Virtual Memory Statistics: (page size of %d bytes)", &Pagesize); err != nil {
		utils.GaiaLogger.Panic("unexpected output of vm_stat: "+line, err.Error())
	}

	memStats := map[string]*float64{
		"Pages free":                   &m.MemoryReply.Free,
		"Pages active":                 &m.MemoryReply.Active,
		"Pages inactive":               &m.MemoryReply.Inactive,
		"Pages speculative":            &Speculative,
		"Pages wired down":             &Wired,
		"Pages purgeable":              &Purgeable,
		"File-backed pages":            &FileBacked,
		"Pages occupied by compressor": &Compressed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		if ptr := memStats[line[:i]]; ptr != nil {
			val := strings.TrimRight(strings.TrimSpace(line[i+1:]), ".")
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				*ptr = v * float64(Pagesize)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("scan error for vm_stat: ", "")
	}

	m.MemoryReply.Cached = Purgeable + FileBacked
	m.MemoryReply.Used = Wired + Compressed + m.MemoryReply.Active + m.MemoryReply.Inactive + Speculative - m.MemoryReply.Cached
	m.MemoryReply.Total = m.MemoryReply.Used + m.MemoryReply.Cached + m.MemoryReply.Free

	m.MemoryReply.Total = utils.ToFixed(m.MemoryReply.Total/1000000, 2)
	m.MemoryReply.Used = utils.ToFixed(m.MemoryReply.Used/1000000, 2)
	m.MemoryReply.Free = utils.ToFixed(m.MemoryReply.Free/1000000, 2)
	m.MemoryReply.Cached = utils.ToFixed(m.MemoryReply.Cached/1000000, 2)
	m.MemoryReply.Active = utils.ToFixed(m.MemoryReply.Active/1000000, 2)
	m.MemoryReply.Inactive = utils.ToFixed(m.MemoryReply.Inactive/1000000, 2)

	return &m.MemoryReply
}

func (m *AMD64Controller) GetDiskStats() *apostolis_pb.DiskReply {
	disks, err := Disk.IOCounters()
	if err != nil {
		utils.GaiaLogger.Error("Error getting disk stats: ", "")
	}

	m.DiskReply.Name = nil
	m.DiskReply.ReadsCompleted = nil
	m.DiskReply.WritesCompleted = nil

	for _, d := range disks {
		m.DiskReply.Name = append(m.DiskReply.Name, d.Name)
		m.DiskReply.ReadsCompleted = append(m.DiskReply.ReadsCompleted, d.ReadCount)
		m.DiskReply.WritesCompleted = append(m.DiskReply.WritesCompleted, d.WriteCount)
	}

	return &m.DiskReply
}

func (m *AMD64Controller) GetNetworkStats() *apostolis_pb.NetReply {
	time.Sleep(1 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Reference: man 1 netstat
	cmd := exec.CommandContext(ctx, "netstat", "-bni")
	out, err := cmd.StdoutPipe()
	if err != nil {
		utils.GaiaLogger.Panic("Error getting network stats: ", err.Error())
	}
	if err := cmd.Start(); err != nil {
		utils.GaiaLogger.Panic("Error getting network stats: ", err.Error())
	}
	scanner := bufio.NewScanner(out)

	if !scanner.Scan() {
		utils.GaiaLogger.Panic("failed to scan output of netstat", "")
	}
	line := scanner.Text()
	if !strings.HasPrefix(line, "Name") {
		utils.GaiaLogger.Panic("unexpected output of netstat -bni: ", line)
	}
	fields := strings.Fields(line)
	fieldsCount := len(fields)

	for i, field := range fields {
		switch field {
		case "Ibytes":
			rxBytesIdx = i
		case "Obytes":
			txBytesIdx = i
		}
	}

	if rxBytesIdx == 0 || txBytesIdx == 0 {
		utils.GaiaLogger.Panic("unexpected output of netstat -bni: ", line)
	}

	m.NetReply.Name = nil
	m.NetReply.RxBytes = nil
	m.NetReply.TxBytes = nil
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := strings.TrimSuffix(fields[0], "*")
		m.NetReply.Name = append(m.NetReply.Name, name)
		if strings.HasPrefix(name, "lo") || !strings.HasPrefix(fields[2], "<Link#") {
			continue
		}
		rxBytesIdx, txBytesIdx := rxBytesIdx, txBytesIdx
		if len(fields) < fieldsCount { // Address can be empty
			rxBytesIdx, txBytesIdx = rxBytesIdx-1, txBytesIdx-1
		}
		rxBytes, err := strconv.ParseUint(fields[rxBytesIdx], 10, 64)
		m.NetReply.RxBytes = append(m.NetReply.RxBytes, rxBytes)
		if err != nil {
			utils.GaiaLogger.Panic("failed to parse Ibytes of %s", err.Error())
		}
		txBytes, err := strconv.ParseUint(fields[txBytesIdx], 10, 64)
		m.NetReply.TxBytes = append(m.NetReply.TxBytes, txBytes)
		if err != nil {
			utils.GaiaLogger.Panic("failed to parse Obytes of %s", err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("scan error for netstat: %s", err.Error())
	}

	return &m.NetReply
}

func (m *AMD64Controller) GetLoadAVGStats() *apostolis_pb.LoadReply {
	load_stats, err := loadavg.Get()
	if err != nil {
		utils.GaiaLogger.Error("Error getting loadavg: ", err.Error())
	}
	m.LoadReply = apostolis_pb.LoadReply{
		Loadavg1:  utils.ToFixed(load_stats.Loadavg1, 2),
		Loadavg5:  utils.ToFixed(load_stats.Loadavg5, 2),
		Loadavg15: utils.ToFixed(load_stats.Loadavg15, 2),
	}

	return &m.LoadReply
}
