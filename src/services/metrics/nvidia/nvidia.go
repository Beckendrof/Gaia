package nvidia

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	CPU "github.com/shirou/gopsutil/cpu"

	apostolis_pb "beckendrof/gaia/src/services/grpc/apostolis"

	"beckendrof/gaia/src/utils"
)

// #include <stdlib.h>
import "C"

var (
	cmd                 *exec.Cmd
	out                 string
	loadavgs            [3]C.double
	StatCount, CPUCount int
	Nice, Iowait, Irq, Softirq, Steal, Guest, GuestNice, PageTables,
	Committed, VmallocUsed, Buffers, SwapCached, SwapFree, SwapTotal,
	Shmem, Slab, Mapped float64
	MemAvailableEnabled bool
)

type Nvidia struct {
	Xavier XavierController
	Orin   OrinController
}

type XavierController struct {
	apostolis_pb.MemoryReply
	apostolis_pb.CPUReply
	apostolis_pb.GPUReply
	apostolis_pb.DiskReply
	apostolis_pb.NetReply
	apostolis_pb.LoadReply
}

type OrinController struct {
	apostolis_pb.MemoryReply
	apostolis_pb.CPUReply
	apostolis_pb.GPUReply
	apostolis_pb.DiskReply
	apostolis_pb.NetReply
	apostolis_pb.LoadReply
}

type cpuStat struct {
	name string
	ptr  *float64
}

func (m *OrinController) GetCPUStats() *apostolis_pb.CPUReply {
	file, err := os.Open("/proc/stat")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/stat", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	cpuStats := []cpuStat{
		{"user", &m.CPUReply.User},
		{"nice", &Nice},
		{"system", &m.CPUReply.System},
		{"idle", &m.CPUReply.Idle},
		{"iowait", &Iowait},
		{"irq", &Irq},
		{"softirq", &Softirq},
		{"steal", &Steal},
		{"guest", &Guest},
		{"guest_nice", &GuestNice},
	}

	if !scanner.Scan() {
		utils.GaiaLogger.Panic("Failed to scan /proc/stat", err.Error())
	}

	valStrs := strings.Fields(scanner.Text())[1:]
	StatCount = len(valStrs)
	for i, valStr := range valStrs {
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse /proc/stat", err.Error())
		}
		*cpuStats[i].ptr = val
		m.CPUReply.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	m.CPUReply.Total -= Guest

	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	m.CPUReply.Total -= GuestNice

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			CPUCount++
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/stat", err.Error())
	}

	m.CPUReply.User = utils.ToFixed(m.CPUReply.User/m.CPUReply.Total*100, 2)
	m.CPUReply.System = utils.ToFixed(m.CPUReply.System/m.CPUReply.Total*100, 2)
	m.CPUReply.Idle = utils.ToFixed(m.CPUReply.Idle/m.CPUReply.Total*100, 2)
	m.CPUReply.Total = utils.ToFixed(m.CPUReply.Total, 2)

	// CORE Usage
	cores, _ := CPU.Percent(0, true)

	m.CPUReply.Cpus = nil
	for _, core := range cores {
		m.CPUReply.Cpus = append(m.CPUReply.Cpus, utils.ToFixed(core, 2))
	}

	// Temperature
	cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	out, _, _ = utils.ReturnExitCode(cmd)
	temp, err := strconv.ParseFloat(out[:3], 64)
	if err != nil {
		utils.GaiaLogger.Error("Failed to parse /sys/class/thermal/thermal_zone0/temp", err.Error())
	}
	m.CPUReply.Temp = utils.ToFixed(temp/10, 2)

	// Power
	cmd = exec.Command("cat", "/sys/class/hwmon/hwmon3/in2_input")
	out, _, _ = utils.ReturnExitCode(cmd)
	cpu_volt, _ := strconv.ParseFloat(out[:len(out)-1], 64)

	cmd = exec.Command("cat", "/sys/class/hwmon/hwmon3/curr2_input")
	out, _, _ = utils.ReturnExitCode(cmd)
	cpu_amps, _ := strconv.ParseFloat(out[:len(out)-1], 64)

	m.CPUReply.Power = utils.ToFixed(cpu_volt*cpu_amps/1000000, 2)
	return &m.CPUReply
}

func (m *XavierController) GetCPUStats() *apostolis_pb.CPUReply {
	file, err := os.Open("/proc/stat")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/stat", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	cpuStats := []cpuStat{
		{"user", &m.CPUReply.User},
		{"nice", &Nice},
		{"system", &m.CPUReply.System},
		{"idle", &m.CPUReply.Idle},
		{"iowait", &Iowait},
		{"irq", &Irq},
		{"softirq", &Softirq},
		{"steal", &Steal},
		{"guest", &Guest},
		{"guest_nice", &GuestNice},
	}

	if !scanner.Scan() {
		utils.GaiaLogger.Panic("Failed to scan /proc/stat", err.Error())
	}

	valStrs := strings.Fields(scanner.Text())[1:]
	StatCount = len(valStrs)
	for i, valStr := range valStrs {
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse /proc/stat", err.Error())
		}
		*cpuStats[i].ptr = val
		m.CPUReply.Total += val
	}

	// Since cpustat[CPUTIME_USER] includes cpustat[CPUTIME_GUEST], subtract the duplicated values from total.
	m.CPUReply.Total -= Guest

	// cpustat[CPUTIME_NICE] includes cpustat[CPUTIME_GUEST_NICE]
	m.CPUReply.Total -= GuestNice

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			CPUCount++
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/stat", err.Error())
	}

	m.CPUReply.User = utils.ToFixed(m.CPUReply.User/m.CPUReply.Total*100, 2)
	m.CPUReply.System = utils.ToFixed(m.CPUReply.System/m.CPUReply.Total*100, 2)
	m.CPUReply.Idle = utils.ToFixed(m.CPUReply.Idle/m.CPUReply.Total*100, 2)
	m.CPUReply.Total = utils.ToFixed(m.CPUReply.Total, 2)

	// CORE Usage
	cores, _ := CPU.Percent(0, true)

	m.CPUReply.Cpus = nil
	for _, core := range cores {
		m.CPUReply.Cpus = append(m.CPUReply.Cpus, utils.ToFixed(core, 2))
	}

	// Temperature
	cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	out, _, _ = utils.ReturnExitCode(cmd)
	temp, err := strconv.ParseFloat(out[:3], 64)
	if err != nil {
		utils.GaiaLogger.Error("Failed to parse temperature", err.Error())
	}
	m.CPUReply.Temp = utils.ToFixed(temp/10, 2)

	// Power
	cmd = exec.Command("cat", "/sys/bus/i2c/drivers/ina3221x/1-0040/iio:device0/in_power1_input")
	out, _, _ = utils.ReturnExitCode(cmd)
	m.CPUReply.Power, err = strconv.ParseFloat(out[:len(out)-1], 64)
	if err != nil {
		utils.GaiaLogger.Error("Failed to parse power", err.Error())
	}

	return &m.CPUReply
}

func (m *XavierController) GetMemoryStats() *apostolis_pb.MemoryReply {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/meminfo", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	memStats := map[string]*float64{
		"MemTotal":     &m.MemoryReply.Total,
		"MemFree":      &m.MemoryReply.Free,
		"MemAvailable": &m.MemoryReply.Available,
		"Buffers":      &Buffers,
		"Cached":       &m.MemoryReply.Cached,
		"Active":       &m.MemoryReply.Active,
		"Inactive":     &m.MemoryReply.Inactive,
		"SwapCached":   &SwapCached,
		"SwapTotal":    &SwapTotal,
		"SwapFree":     &SwapFree,
		"Mapped":       &Mapped,
		"Shmem":        &Shmem,
		"Slab":         &Slab,
		"PageTables":   &PageTables,
		"Committed_AS": &Committed,
		"VmallocUsed":  &VmallocUsed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		fld := line[:i]
		if ptr := memStats[fld]; ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				*ptr = v * 1024
			}
			if fld == "MemAvailable" {
				MemAvailableEnabled = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/meminfo", err.Error())
	}

	if MemAvailableEnabled {
		m.MemoryReply.Used = m.MemoryReply.Total - m.MemoryReply.Available
	} else {
		m.MemoryReply.Used = m.MemoryReply.Total - m.MemoryReply.Free - Buffers - m.MemoryReply.Cached
	}

	m.MemoryReply.Total = utils.ToFixed(m.MemoryReply.Total/1000000, 2)
	m.MemoryReply.Used = utils.ToFixed(m.MemoryReply.Used/1000000, 2)
	m.MemoryReply.Free = utils.ToFixed(m.MemoryReply.Free/1000000, 2)
	m.MemoryReply.Cached = utils.ToFixed(m.MemoryReply.Cached/1000000, 2)
	m.MemoryReply.Active = utils.ToFixed(m.MemoryReply.Active/1000000, 2)
	m.MemoryReply.Inactive = utils.ToFixed(m.MemoryReply.Inactive/1000000, 2)

	return &m.MemoryReply
}

func (m *OrinController) GetMemoryStats() *apostolis_pb.MemoryReply {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/meminfo", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	memStats := map[string]*float64{
		"MemTotal":     &m.MemoryReply.Total,
		"MemFree":      &m.MemoryReply.Free,
		"MemAvailable": &m.MemoryReply.Available,
		"Buffers":      &Buffers,
		"Cached":       &m.MemoryReply.Cached,
		"Active":       &m.MemoryReply.Active,
		"Inactive":     &m.MemoryReply.Inactive,
		"SwapCached":   &SwapCached,
		"SwapTotal":    &SwapTotal,
		"SwapFree":     &SwapFree,
		"Mapped":       &Mapped,
		"Shmem":        &Shmem,
		"Slab":         &Slab,
		"PageTables":   &PageTables,
		"Committed_AS": &Committed,
		"VmallocUsed":  &VmallocUsed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		fld := line[:i]
		if ptr := memStats[fld]; ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				*ptr = v * 1024
			}
			if fld == "MemAvailable" {
				MemAvailableEnabled = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/meminfo", err.Error())
	}

	if MemAvailableEnabled {
		m.MemoryReply.Used = m.MemoryReply.Total - m.MemoryReply.Available
	} else {
		m.MemoryReply.Used = m.MemoryReply.Total - m.MemoryReply.Free - Buffers - m.MemoryReply.Cached
	}

	m.MemoryReply.Total = utils.ToFixed(m.MemoryReply.Total/1000000, 2)
	m.MemoryReply.Used = utils.ToFixed(m.MemoryReply.Used/1000000, 2)
	m.MemoryReply.Free = utils.ToFixed(m.MemoryReply.Free/1000000, 2)
	m.MemoryReply.Cached = utils.ToFixed(m.MemoryReply.Cached/1000000, 2)
	m.MemoryReply.Active = utils.ToFixed(m.MemoryReply.Active/1000000, 2)
	m.MemoryReply.Inactive = utils.ToFixed(m.MemoryReply.Inactive/1000000, 2)

	return &m.MemoryReply
}

func (m *OrinController) GetDiskStats() *apostolis_pb.DiskReply {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/diskstats", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	m.DiskReply.Name = nil
	m.DiskReply.ReadsCompleted = nil
	m.DiskReply.WritesCompleted = nil
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		m.DiskReply.Name = append(m.DiskReply.Name, fields[2])

		readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse reads completed of "+fields[2], err.Error())
		}
		m.DiskReply.ReadsCompleted = append(m.DiskReply.ReadsCompleted, readsCompleted)

		writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse writes completed of "+fields[2], err.Error())
		}
		m.DiskReply.WritesCompleted = append(m.DiskReply.WritesCompleted, writesCompleted)
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/diskstats", err.Error())
	}
	return &m.DiskReply
}

func (m *XavierController) GetDiskStats() *apostolis_pb.DiskReply {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/diskstats", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	m.DiskReply.Name = nil
	m.DiskReply.ReadsCompleted = nil
	m.DiskReply.WritesCompleted = nil
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		m.DiskReply.Name = append(m.DiskReply.Name, fields[2])

		readsCompleted, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse reads completed of "+fields[2], err.Error())
		}
		m.DiskReply.ReadsCompleted = append(m.DiskReply.ReadsCompleted, readsCompleted)

		writesCompleted, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse writes completed of "+fields[2], err.Error())
		}
		m.DiskReply.WritesCompleted = append(m.DiskReply.WritesCompleted, writesCompleted)
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/diskstats", err.Error())
	}
	return &m.DiskReply
}

func (m *OrinController) GetNetworkStats() *apostolis_pb.NetReply {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/net/dev", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	m.NetReply.Name = nil
	m.NetReply.RxBytes = nil
	m.NetReply.TxBytes = nil
	for scanner.Scan() {
		// Reference: dev_seq_printf_stats in Linux source code
		kv := strings.SplitN(scanner.Text(), ":", 2)
		if len(kv) != 2 {
			continue
		}
		fields := strings.Fields(kv[1])
		if len(fields) < 16 {
			continue
		}
		name := strings.TrimSpace(kv[0])
		m.NetReply.Name = append(m.NetReply.Name, name)
		if name == "lo" {
			continue
		}
		rxBytes, err := strconv.ParseUint(fields[0], 10, 64)
		m.NetReply.RxBytes = append(m.NetReply.RxBytes, rxBytes)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse rxBytes of "+name, err.Error())
		}
		txBytes, err := strconv.ParseUint(fields[8], 10, 64)
		m.NetReply.TxBytes = append(m.NetReply.TxBytes, txBytes)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse txBytes of "+name, err.Error())
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/net/dev", err.Error())
	}
	return &m.NetReply
}

func (m *XavierController) GetNetworkStats() *apostolis_pb.NetReply {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		utils.GaiaLogger.Panic("Failed to open /proc/net/dev", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	m.NetReply.Name = nil
	m.NetReply.RxBytes = nil
	m.NetReply.TxBytes = nil
	for scanner.Scan() {
		// Reference: dev_seq_printf_stats in Linux source code
		kv := strings.SplitN(scanner.Text(), ":", 2)
		if len(kv) != 2 {
			continue
		}
		fields := strings.Fields(kv[1])
		if len(fields) < 16 {
			continue
		}
		name := strings.TrimSpace(kv[0])
		m.NetReply.Name = append(m.NetReply.Name, name)
		if name == "lo" {
			continue
		}
		rxBytes, err := strconv.ParseUint(fields[0], 10, 64)
		m.NetReply.RxBytes = append(m.NetReply.RxBytes, rxBytes)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse rxBytes of "+name, err.Error())
		}
		txBytes, err := strconv.ParseUint(fields[8], 10, 64)
		m.NetReply.TxBytes = append(m.NetReply.TxBytes, txBytes)
		if err != nil {
			utils.GaiaLogger.Error("Failed to parse txBytes of "+name, err.Error())
		}
	}
	if err := scanner.Err(); err != nil {
		utils.GaiaLogger.Panic("Failed to scan /proc/net/dev", err.Error())
	}
	return &m.NetReply
}

func (m *OrinController) GetGPUStats() *apostolis_pb.GPUReply {
	time.Sleep(1 * time.Second)

	// Load
	cmd = exec.Command("sudo", "-S", "cat", "/sys/devices/gpu.0/load")
	out, _, _ = utils.ReturnExitCode(cmd)

	m.GPUReply.Percent, _ = strconv.ParseFloat(out[:len(out)-1], 64)

	// Temperature
	m.GPUReply.Temp = nil
	cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone1/temp")
	out, _, _ = utils.ReturnExitCode(cmd)
	temp, err := strconv.ParseFloat(out[:3], 64)
	if err != nil {
		utils.GaiaLogger.Error("Failed to parse Temp", err.Error())
	}
	m.GPUReply.Temp = append(m.GPUReply.Temp, utils.ToFixed(temp/10, 2))

	// Name
	m.GPUReply.DeviceName = nil
	m.GPUReply.DeviceName = append(m.GPUReply.DeviceName, "NVIDIA Tegra Orin")

	// Power
	m.GPUReply.Power = nil
	cmd = exec.Command("cat", "/sys/class/hwmon/hwmon3/in1_input")
	out, _, _ = utils.ReturnExitCode(cmd)
	gpu_volt, _ := strconv.ParseFloat(out[:len(out)-1], 64)

	cmd = exec.Command("cat", "/sys/class/hwmon/hwmon3/curr1_input")
	out, _, _ = utils.ReturnExitCode(cmd)
	gpu_amps, err := strconv.ParseFloat(out[:len(out)-1], 64)
	if err != nil {
		utils.GaiaLogger.Error("Failed to parse current", err.Error())
	}

	m.GPUReply.Power = append(m.GPUReply.Power, utils.ToFixed(gpu_volt*gpu_amps/1000000, 2))

	return &m.GPUReply
}

func (m *XavierController) GetGPUStats() *apostolis_pb.GPUReply {
	time.Sleep(1 * time.Second)
	cmd = exec.Command("cat", "/sys/devices/gpu.0/load")
	out, _, _ = utils.ReturnExitCode(cmd)
	m.GPUReply.Percent, _ = strconv.ParseFloat(out[:len(out)-1], 64)

	// Temp
	m.GPUReply.Temp = nil
	cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone1/temp")
	out, _, _ = utils.ReturnExitCode(cmd)
	temp, _ := strconv.ParseFloat(out, 64)
	m.GPUReply.Temp = append(m.GPUReply.Temp, utils.ToFixed(temp/10, 2))

	// Name
	m.GPUReply.DeviceName = nil
	m.GPUReply.DeviceName = append(m.GPUReply.DeviceName, "NVIDIA Tegra Xavier")

	// Power
	m.GPUReply.Power = nil
	cmd = exec.Command("cat", "/sys/bus/i2c/drivers/ina3221x/1-0040/iio:device0/in_power0_input")
	out, _, _ = utils.ReturnExitCode(cmd)

	gpu_power, _ := strconv.ParseFloat(out[:len(out)-1], 64)
	m.GPUReply.Power = append(m.GPUReply.Power, utils.ToFixed(gpu_power/1000, 4))

	return &m.GPUReply
}

func (m *OrinController) GetLoadAVGStats() *apostolis_pb.LoadReply {
	ret := C.getloadavg(&loadavgs[0], 3)
	if ret != 3 {
		utils.GaiaLogger.Error("Error getting loadavg: ", "")
	}

	m.LoadReply.Loadavg1 = utils.ToFixed(float64(loadavgs[0]), 2)
	m.LoadReply.Loadavg5 = utils.ToFixed(float64(loadavgs[1]), 2)
	m.LoadReply.Loadavg15 = utils.ToFixed(float64(loadavgs[2]), 2)
	return &m.LoadReply
}

func (m *XavierController) GetLoadAVGStats() *apostolis_pb.LoadReply {
	ret := C.getloadavg(&loadavgs[0], 3)
	if ret != 3 {
		utils.GaiaLogger.Error("Error getting loadavg: ", "")
	}

	m.LoadReply.Loadavg1 = utils.ToFixed(float64(loadavgs[0]), 2)
	m.LoadReply.Loadavg5 = utils.ToFixed(float64(loadavgs[1]), 2)
	m.LoadReply.Loadavg15 = utils.ToFixed(float64(loadavgs[2]), 2)
	return &m.LoadReply
}
