package main

import (
	"context"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type RequestContext struct {
    context context.Context
}

func NewRequestContext(ctx context.Context) *RequestContext {
    return &RequestContext{ctx}
}

func (r *RequestContext) Hostname() (string, error) {
    return os.Hostname()
}

func (r *RequestContext) Platform() (string, error) {
    platform, _, _, err := host.PlatformInformationWithContext(r.context)
    return platform, err
}

func (r *RequestContext) PlatformFamily() (string, error) {
    _, family, _, err := host.PlatformInformationWithContext(r.context)
    return family, err
}


func (r *RequestContext) PlatformVersion() (string, error) {
    _, _, version, err := host.PlatformInformationWithContext(r.context)
    return version, err
}

func (r *RequestContext) KernelVersion() (string, error) {
    return host.KernelVersionWithContext(r.context)
}

func (r *RequestContext) KernelArch() (string, error) {
    return host.KernelArch()
}

func (r *RequestContext) BootTime() (uint64, error) {
    return host.BootTimeWithContext(r.context)
}

func (r *RequestContext) Uptime() (uint64, error) {
    return host.UptimeWithContext(r.context)
}

func (r *RequestContext) HostID() (string, error) {
    return host.HostIDWithContext(r.context)
}

func (r *RequestContext) Users() ([]host.UserStat, error) {
    return host.UsersWithContext(r.context)
}

func (r *RequestContext) Temperatures() ([]host.TemperatureStat, error) {
    return host.SensorsTemperaturesWithContext(r.context)
}

func (r *RequestContext) CPUPhysicalCoreNumber() (int, error) {
    return cpu.CountsWithContext(r.context, false)
}

func (r *RequestContext) CPULogicalCoreNumber() (int, error) {
    return cpu.CountsWithContext(r.context, true)
}

func (r *RequestContext) CPUUsagePerCPU() ([]float64, error) {
    return cpu.PercentWithContext(r.context, 0, true)
}

func (r *RequestContext) CPUUsage() (float64, error) {
    vals, err := cpu.PercentWithContext(r.context, 0, false)
    if err != nil {
        return 0, err
    }
    return vals[0], nil
}

func (r *RequestContext) DiskUsage(path string) (*disk.UsageStat, error) {
    return disk.UsageWithContext(r.context, path)
}

func (r *RequestContext) AvgLoad() (*load.AvgStat, error) {
    return load.AvgWithContext(r.context)
}

func (r *RequestContext) ProcsRunning() (int, error) {
    misc, err := load.MiscWithContext(r.context)
    if err != nil {
        return 0, err
    }
    return misc.ProcsRunning, nil
}

func (r *RequestContext) ProcsTotal() (int, error) {
    misc, err := load.MiscWithContext(r.context)
    if err != nil {
        return 0, err
    }
    return misc.ProcsTotal, nil
}

func (r *RequestContext) ProcsCreated() (int, error) {
    misc, err := load.MiscWithContext(r.context)
    if err != nil {
        return 0, err
    }
    return misc.ProcsCreated, nil
}

func (r *RequestContext) ProcsBlocked() (int, error) {
    misc, err := load.MiscWithContext(r.context)
    if err != nil {
        return 0, err
    }
    return misc.ProcsBlocked, nil
}

func (r *RequestContext) Memory() (*mem.VirtualMemoryStat, error) {
    return mem.VirtualMemoryWithContext(r.context)
}

func (r *RequestContext) Swap() (*mem.SwapMemoryStat, error) {
    return mem.SwapMemoryWithContext(r.context)
}

func (r *RequestContext) SwapDevices() ([]*mem.SwapDevice, error) {
    return mem.SwapDevicesWithContext(r.context)
}

func (r *RequestContext) Processes() ([]*process.Process, error) {
    return process.ProcessesWithContext(r.context)
}

func (r *RequestContext) ProcessPID(pid int32) (*process.Process, error) {
    return process.NewProcessWithContext(r.context, pid)
}

func (r *RequestContext) ProcessExistsPID(pid int32) (bool, error) {
    return process.PidExistsWithContext(r.context, pid)
}

func (r *RequestContext) Now() (int64) {
    return time.Now().Unix()
}
