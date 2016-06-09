package analyze

import (
	"logging"
	"os"
	"strings"
	"time"
)

type ReportData struct {
	GenerationDate time.Time
	Invocation     string
	System         SystemData
	Instances      waarpInstancesData
}

func NewReport() *ReportData {
	return &ReportData{
		GenerationDate: time.Now().UTC(),
		Invocation:     strings.Join(os.Args, " "),
		System:         SystemData{},
		Instances:      waarpInstancesData{},
	}
}

func (rd *ReportData) Collect() {
	rd.System.Collect()
	rd.Instances.Collect()
}

type waarpInstancesData []WaarpInstanceData

func (w *waarpInstancesData) DetectRunning() {
	logging.Info.Print("Looking for instance...")

	out, ok := runCmd("ps", "wwh", "-C", "java", "o", "pid,cmd")
	if !ok {
		return
	}

	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "org.waarp.openr66.server.R66Server") {
			continue
		}

		wid := WaarpInstanceData{
			Pid:       toInt(strings.Fields(line)[0]),
			Type:      "Waarp R66 Server",
			HostID:    "unknown",
			HostIDSSL: "unknown",
		}
		*w = append(*w, wid)
		logging.Info.Printf("Found running instance with PID %d", wid.Pid)
	}
}

func (w waarpInstancesData) Collect() {
	for i, inst := range w {
		inst.Collect()
		w[i] = inst
	}
}
