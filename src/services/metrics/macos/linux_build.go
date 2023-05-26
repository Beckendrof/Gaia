//go:build linux
// +build linux

package macos

import (
	apostolis_pb "beckendrof/gaia/src/services/grpc/apostolis"
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

func (m *AMD64Controller) GetMemoryStats() *apostolis_pb.MemoryReply {
	return &m.MemoryReply
}

func (m *AMD64Controller) GetDiskStats() *apostolis_pb.DiskReply {
	return &m.DiskReply
}

func (m *AMD64Controller) GetNetworkStats() *apostolis_pb.NetReply {
	return &m.NetReply
}

func (m *AMD64Controller) GetLoadAVGStats() *apostolis_pb.LoadReply {
	return &m.LoadReply
}

func (m *AMD64Controller) GetCPUStats() *apostolis_pb.CPUReply {
	return &m.CPUReply
}

func (m *AMD64Controller) GetGPUStats() *apostolis_pb.GPUReply {
	return &m.GPUReply
}
