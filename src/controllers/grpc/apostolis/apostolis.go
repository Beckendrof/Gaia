package apostolis_grpc_controller

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apostolis_pb "beckendrof/gaia/src/services/grpc/apostolis"
	"beckendrof/gaia/src/services/metrics"
)

// server is used to implement apostolis_pb.AtratoServer.
type ApostolisServer struct {
	apostolis_pb.UnimplementedApostolisServer // Embed the unimplemented server
}

func (s *ApostolisServer) System(in *apostolis_pb.ApostolisRequest, stream apostolis_pb.Apostolis_SystemServer) error {
	if in.GetMetric() == apostolis_pb.Stats_Memory {
		for {
			time.Sleep(1 * time.Second)
			metrics.Instance.GetMemoryStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_MemoryUsage{
						MemoryUsage: &apostolis_pb.MemoryReply{
							Total:    metrics.Instance.MemoryStats.Total,
							Used:     metrics.Instance.MemoryStats.Used,
							Cached:   metrics.Instance.MemoryStats.Cached,
							Free:     metrics.Instance.MemoryStats.Free,
							Active:   metrics.Instance.MemoryStats.Active,
							Inactive: metrics.Instance.MemoryStats.Inactive,
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	} else if in.GetMetric() == apostolis_pb.Stats_Disk {
		for {
			time.Sleep(1 * time.Second)
			metrics.Instance.GetDiskStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_DiskUsage{
						DiskUsage: &apostolis_pb.DiskReply{
							Name:            metrics.Instance.DiskStats.Name,
							ReadsCompleted:  metrics.Instance.DiskStats.ReadsCompleted,
							WritesCompleted: metrics.Instance.DiskStats.WritesCompleted,
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	} else if in.GetMetric() == apostolis_pb.Stats_LoadAvg {
		for {
			time.Sleep(1 * time.Second)
			metrics.Instance.GetLoadAVGStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_LoadUsage{
						LoadUsage: &apostolis_pb.LoadReply{
							Loadavg1:  metrics.Instance.LoadStats.Loadavg1,
							Loadavg5:  metrics.Instance.LoadStats.Loadavg5,
							Loadavg15: metrics.Instance.LoadStats.Loadavg15,
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	} else if in.GetMetric() == apostolis_pb.Stats_CPU {
		for {
			time.Sleep(1 * time.Second)
			metrics.Instance.GetCPUStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_CpuUsage{
						CpuUsage: &apostolis_pb.CPUReply{
							Total:  metrics.Instance.CpuStats.Total,
							User:   metrics.Instance.CpuStats.User,
							System: metrics.Instance.CpuStats.System,
							Idle:   metrics.Instance.CpuStats.Idle,
							Cpus:   metrics.Instance.CpuStats.Cpus,
							Temp:   metrics.Instance.CpuStats.Temp,
							Power:  metrics.Instance.CpuStats.Power,
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	} else if in.GetMetric() == apostolis_pb.Stats_Net {
		for {
			time.Sleep(1 * time.Second)
			metrics.Instance.GetNetworkStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_NetUsage{
						NetUsage: &apostolis_pb.NetReply{
							Name:    metrics.Instance.NetStats.Name,
							RxBytes: metrics.Instance.NetStats.RxBytes,
							TxBytes: metrics.Instance.NetStats.TxBytes,
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	} else if in.GetMetric() == apostolis_pb.Stats_GPU {
		for {
			metrics.Instance.GetGPUStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_GpuUsage{
						GpuUsage: &apostolis_pb.GPUReply{
							DeviceName: metrics.Instance.GpuStats.DeviceName,
							Temp:       metrics.Instance.GpuStats.Temp,
							Percent:    metrics.Instance.GpuStats.Percent,
							Used:       metrics.Instance.GpuStats.Used,
							Total:      metrics.Instance.GpuStats.Total,
							Power:      metrics.Instance.GpuStats.Power,
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	} else {
		for {
			metrics.Instance.GetAllStats()
			err := stream.Send(
				&apostolis_pb.ApostolisReply{
					Status:  0,
					Message: "Successful",
					Data: &apostolis_pb.ApostolisReply_All{
						All: &apostolis_pb.AllUsageStats{
							MemoryUsage: &apostolis_pb.MemoryReply{
								Total:    metrics.Instance.MemoryStats.Total,
								Used:     metrics.Instance.MemoryStats.Used,
								Cached:   metrics.Instance.MemoryStats.Cached,
								Free:     metrics.Instance.MemoryStats.Free,
								Active:   metrics.Instance.MemoryStats.Active,
								Inactive: metrics.Instance.MemoryStats.Inactive,
							},
							LoadUsage: &apostolis_pb.LoadReply{
								Loadavg1:  metrics.Instance.LoadStats.Loadavg1,
								Loadavg5:  metrics.Instance.LoadStats.Loadavg5,
								Loadavg15: metrics.Instance.LoadStats.Loadavg15,
							},
							NetUsage: &apostolis_pb.NetReply{
								Name:    metrics.Instance.NetStats.Name,
								RxBytes: metrics.Instance.NetStats.RxBytes,
								TxBytes: metrics.Instance.NetStats.TxBytes,
							},
							CpuUsage: &apostolis_pb.CPUReply{
								Total:  metrics.Instance.CpuStats.Total,
								User:   metrics.Instance.CpuStats.User,
								System: metrics.Instance.CpuStats.System,
								Idle:   metrics.Instance.CpuStats.Idle,
								Cpus:   metrics.Instance.CpuStats.Cpus,
								Temp:   metrics.Instance.CpuStats.Temp,
								Power:  metrics.Instance.CpuStats.Power,
							},
							GpuUsage: &apostolis_pb.GPUReply{
								DeviceName: metrics.Instance.GpuStats.DeviceName,
								Temp:       metrics.Instance.GpuStats.Temp,
								Percent:    metrics.Instance.GpuStats.Percent,
								Used:       metrics.Instance.GpuStats.Used,
								Total:      metrics.Instance.GpuStats.Total,
								Power:      metrics.Instance.GpuStats.Power,
							},
							DiskUsage: &apostolis_pb.DiskReply{
								Name:            metrics.Instance.DiskStats.Name,
								ReadsCompleted:  metrics.Instance.DiskStats.ReadsCompleted,
								WritesCompleted: metrics.Instance.DiskStats.WritesCompleted,
							},
						},
					},
				},
			)
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}
