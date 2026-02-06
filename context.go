package godashboard

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

// RequestContext holds the context for a dashboard request and exposes system information
// to the template engine. It wraps the standard context.Context and provides helper methods
// to retrieve system stats like CPU, memory, and disk usage.
//
// This struct is the primary data source passed to the templates when rendering blocks.
type RequestContext struct {
	context      context.Context
	sizeBaseline int
}

// NewRequestContext creates a new RequestContext with the given background context.
// It initializes the sizeBaseline to 100.
func NewRequestContext(ctx context.Context) *RequestContext {
	return &RequestContext{
		context:      ctx,
		sizeBaseline: 100,
	}
}

// Hostname retrieves the kernel's hostname (e.g. from /proc/sys/kernel/hostname).
func (r *RequestContext) Hostname() (string, error) {
	return os.Hostname()
}

// Platform identifies the OS distribution (e.g. "ubuntu", "alpine") via /etc/os-release or LSB.
func (r *RequestContext) Platform() (string, error) {
	platform, _, _, err := host.PlatformInformationWithContext(r.context)
	return platform, err
}

// PlatformFamily identifies the OS family (e.g. "debian", "rhel").
func (r *RequestContext) PlatformFamily() (string, error) {
	_, family, _, err := host.PlatformInformationWithContext(r.context)
	return family, err
}

// PlatformVersion returns the OS distribution version (e.g. "20.04").
func (r *RequestContext) PlatformVersion() (string, error) {
	_, _, version, err := host.PlatformInformationWithContext(r.context)
	return version, err
}

// KernelVersion returns the kernel release version string (e.g. "5.4.0-generic").
func (r *RequestContext) KernelVersion() (string, error) {
	return host.KernelVersionWithContext(r.context)
}

// KernelArch returns the machine hardware name (e.g. "x86_64").
func (r *RequestContext) KernelArch() (string, error) {
	return host.KernelArch()
}

// BootTime returns the system boot timestamp in seconds since the epoch.
func (r *RequestContext) BootTime() (uint64, error) {
	return host.BootTimeWithContext(r.context)
}

// Uptime retrieves the system uptime in seconds (excluding suspend time).
func (r *RequestContext) Uptime() (uint64, error) {
	return host.UptimeWithContext(r.context)
}

// HostID returns the unique host ID, typically read from /etc/machine-id.
func (r *RequestContext) HostID() (string, error) {
	return host.HostIDWithContext(r.context)
}

// Users lists users currently logged in to the system (from utmp/wtmp).
func (r *RequestContext) Users() ([]host.UserStat, error) {
	return host.UsersWithContext(r.context)
}

// Temperatures retrieves readings from available hardware sensors (e.g. coretemp).
func (r *RequestContext) Temperatures() ([]host.TemperatureStat, error) {
	return host.SensorsTemperaturesWithContext(r.context)
}

// CPUPhysicalCoreNumber counts the number of physical CPU cores (excluding HyperThreading).
func (r *RequestContext) CPUPhysicalCoreNumber() (int, error) {
	return cpu.CountsWithContext(r.context, false)
}

// CPULogicalCoreNumber counts the number of logical CPU threads (including HyperThreading).
func (r *RequestContext) CPULogicalCoreNumber() (int, error) {
	return cpu.CountsWithContext(r.context, true)
}

// CPUUsagePerCPU calculates the usage percentage for each individual CPU core.
func (r *RequestContext) CPUUsagePerCPU() ([]float64, error) {
	return cpu.PercentWithContext(r.context, 0, true)
}

// CPUUsage calculates the total CPU usage percentage across all cores.
func (r *RequestContext) CPUUsage() (float64, error) {
	vals, err := cpu.PercentWithContext(r.context, 0, false)
	if err != nil {
		return 0, err
	}
	return vals[0], nil
}

// DiskUsage returns storage statistics (free/used/total) for the mounted path.
func (r *RequestContext) DiskUsage(path string) (*disk.UsageStat, error) {
	return disk.UsageWithContext(r.context, path)
}

// AvgLoad returns the system load averages (1, 5, 15 minutes).
func (r *RequestContext) AvgLoad() (*load.AvgStat, error) {
	return load.AvgWithContext(r.context)
}

// ProcsRunning returns the number of processes currently in the 'Running' state.
func (r *RequestContext) ProcsRunning() (int, error) {
	misc, err := load.MiscWithContext(r.context)
	if err != nil {
		return 0, err
	}
	return misc.ProcsRunning, nil
}

// ProcsTotal returns the total number of processes existing on the system.
func (r *RequestContext) ProcsTotal() (int, error) {
	misc, err := load.MiscWithContext(r.context)
	if err != nil {
		return 0, err
	}
	return misc.ProcsTotal, nil
}

// ProcsCreated returns the total number of processes created since boot (forks).
func (r *RequestContext) ProcsCreated() (int, error) {
	misc, err := load.MiscWithContext(r.context)
	if err != nil {
		return 0, err
	}
	return misc.ProcsCreated, nil
}

// ProcsBlocked returns the number of processes currently blocked on I/O.
func (r *RequestContext) ProcsBlocked() (int, error) {
	misc, err := load.MiscWithContext(r.context)
	if err != nil {
		return 0, err
	}
	return misc.ProcsBlocked, nil
}

// Memory returns virtual memory statistics (RAM usage, available, etc.).
func (r *RequestContext) Memory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemoryWithContext(r.context)
}

// Swap returns swap memory statistics (used/free swap).
func (r *RequestContext) Swap() (*mem.SwapMemoryStat, error) {
	return mem.SwapMemoryWithContext(r.context)
}

// SwapDevices returns details about available swap devices/partitions.
func (r *RequestContext) SwapDevices() ([]*mem.SwapDevice, error) {
	return mem.SwapDevicesWithContext(r.context)
}

// Processes retrieves a list of all currently running processes with their details.
func (r *RequestContext) Processes() ([]*process.Process, error) {
	return process.ProcessesWithContext(r.context)
}

// ProcessPID creates a handle to inspect a specific process by its ID.
func (r *RequestContext) ProcessPID(pid int32) (*process.Process, error) {
	return process.NewProcessWithContext(r.context, pid)
}

// ProcessExistsPID checks if a process with the given PID currently exists.
func (r *RequestContext) ProcessExistsPID(pid int32) (bool, error) {
	return process.PidExistsWithContext(r.context, pid)
}

// Now returns the current server-side timestamp (Unix epoch), useful for template logic.
func (r *RequestContext) Now() int64 {
	return time.Now().Unix()
}
