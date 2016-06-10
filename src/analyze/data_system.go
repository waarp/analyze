package analyze

import (
	"logging"
	"regexp"
	"strings"
)

type SystemData struct {
	System       string
	Distribution string
	Kernel       string
	Arch         string
	Load         string
	Disks        string

	CPUs             int
	CPUInfos         string
	RAM              int
	RAMInfos         string
	RunningProcesses string
}

func (sd *SystemData) Collect() {
	logging.Info.Print("Collecting system data...")

	setters := []func(*SystemData){
		setCPUs, setCPUInfos, setRAMInfos, setKernelInfo, setDistribution,
		setLoad, setDisks, setRunningProcesses,
	}
	for _, setter := range setters {
		setter(sd)
	}
}

//
// Collect Utils
//

func setCPUs(sd *SystemData) {
	out, ok := runCmd("nproc")
	if !ok {
		return
	}

	sd.CPUs = toInt(out)
}

func setCPUInfos(sd *SystemData) {
	sd.CPUInfos, _ = runCmd("lscpu")
}

func setRAMInfos(sd *SystemData) {
	out, ok := readFile("/proc/meminfo")
	sd.RAMInfos = out
	if !ok {
		return
	}

	var rxMemTotal = regexp.MustCompile("(?m)^MemTotal[^0-9]+([0-9]+).*$")
	res := rxMemTotal.FindStringSubmatch(out)
	if len(res) <= 1 {
		logging.Warning.Print("MemTotal not found in output")
		return
	}

	sd.RAM = toInt(res[1])
}

func setKernelInfo(sd *SystemData) {
	sd.Arch, _ = runCmd("uname", "-m")
	sd.Kernel, _ = runCmd("uname", "-srv")
	sd.System, _ = runCmd("uname", "-o")
}

func setDistribution(sd *SystemData) {
	out, ok := readFile("/etc/os-release")
	if !ok {
		sd.Distribution = out
		return
	}

	var (
		rxPretty  = regexp.MustCompile(`PRETTY_NAME="([^"]*?)"`)
		rxName    = regexp.MustCompile(`NAME="([^"]*?)"`)
		rxVersion = regexp.MustCompile(`VERSION="([^"]*?)"`)
	)

	switch {
	case strings.Contains(out, `PRETTY_NAME=`):
		res := rxPretty.FindStringSubmatch(out)
		if len(res) > 1 {
			sd.Distribution = res[1]
		}

	case strings.Contains(out, `NAME=`) && strings.Contains(out, `VERSION=`):
		res := rxName.FindStringSubmatch(out)
		if len(res) > 1 {
			sd.Distribution = res[1]
		}

		res = rxVersion.FindStringSubmatch(out)
		if len(res) > 1 {
			sd.Distribution += " " + res[1]
		}

	case strings.Contains(out, `NAME=`):
		res := rxName.FindStringSubmatch(out)
		if len(res) > 1 {
			sd.Distribution = res[1]
		}

	default:
		sd.Distribution = out
	}
}

func setLoad(sd *SystemData) {
	sd.Load, _ = readFile("/proc/loadavg")
}

func setDisks(sd *SystemData) {
	sd.Disks, _ = runCmd("df", "-h")
}

func setRunningProcesses(sd *SystemData) {
	sd.RunningProcesses, _ = runCmd("ps", "auxf")
}
