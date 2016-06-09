package analyze

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"logging"

	wxml "alm.waarp.fr/go/go-utils/waarp/xml"
)

type WaarpInstanceData struct {
	Type              string
	HostID, HostIDSSL string
	ConfFile, Config  string

	LogFile, Logs string

	Pid           int
	Cmd           []string
	Cwd           string
	RSS, VSZ      int
	PCPU, PRAM    float64
	User, Group   string
	Uid, Gid      string
	NThreads, NFD int
	Environ       string
	Limits        string
	Status        string
	fds           []string

	Sockets string

	JVMPath, JVMInfos string
	Jars              string
}

func (w WaarpInstanceData) Running() bool {
	return w.Pid != 0
}

func (w *WaarpInstanceData) Collect() {
	logging.Info.Printf("Collecting data for instance with PID %d...", w.Pid)

	setters := []func(*WaarpInstanceData){
		setCmd, setCwd, setPs, setEnviron, setLimits, setCounts,
		setStatus, setFds, setJars, setSockets, setJvm, setConf,
		setHostID, setLogs,
	}
	for _, setter := range setters {
		setter(w)
	}

}

func setCmd(w *WaarpInstanceData) {
	out, _ := readProcFile(w.Pid, "cmdline")
	w.Cmd = strings.Split(out, "\000")
}

func setCwd(w *WaarpInstanceData) {
	w.Cwd, _ = readLink(fmt.Sprintf("/proc/%d/cwd", w.Pid))
}

func setEnviron(w *WaarpInstanceData) {
	out, _ := readProcFile(w.Pid, "environ")
	w.Environ = strings.Replace(out, "\000", "\n", -1)
}

func setLimits(w *WaarpInstanceData) {
	w.Limits, _ = readProcFile(w.Pid, "limits")
}

func setPs(w *WaarpInstanceData) {
	out, ok := runCmd("ps", "h", "p", strconv.Itoa(w.Pid),
		"o", "user,uid,group,gid,rss,%mem,vsz,%cpu")
	if !ok {
		return
	}

	fields := strings.Fields(out)

	w.User = fields[0]
	w.Uid = fields[1]
	w.Group = fields[2]
	w.Gid = fields[3]
	w.RSS = toInt(fields[4])
	w.VSZ = toInt(fields[6])
	w.PRAM = toFloat(fields[5])
	w.PCPU = toFloat(fields[7])
}

func setCounts(w *WaarpInstanceData) {
	w.NFD = countFilesInDir(fmt.Sprintf("/proc/%d/fd", w.Pid))
	w.NThreads = countFilesInDir(fmt.Sprintf("/proc/%d/task", w.Pid))
}

func setStatus(w *WaarpInstanceData) {
	out, _ := readProcFile(w.Pid, "status")
	w.Status = strings.Replace(out, "\000", "\n", -1)
}

func setFds(w *WaarpInstanceData) {
	out, _ := runCmd("lsof", "-nP", "-p", strconv.Itoa(w.Pid))
	lines := strings.Split(out, "\n")
	w.fds = lines[1:len(lines)]
}

func setJars(w *WaarpInstanceData) {
	jars := []string{}
	knownJars := map[string]bool{}

	for _, file := range w.fds {
		if !strings.Contains(file, ".jar") {
			continue
		}
		jarpath := file[strings.LastIndex(file, " "):len(file)]
		if knownJars[jarpath] {
			continue
		}

		jars = append(jars, jarpath)
		knownJars[jarpath] = true
	}
	w.Jars = strings.Join(jars, "\n")
}

func setSockets(w *WaarpInstanceData) {
	socks := []string{}

	for _, line := range w.fds {
		logging.Debug.Printf("line: %s", line)
		fields := strings.Fields(line)
		switch fields[4] {
		case "IPv4", "IPv6", "unix":
			socks = append(socks, strings.Join(fields[7:len(fields)], " "))
		}
	}
	w.Sockets = strings.Join(socks, "\n")
}

func setJvm(w *WaarpInstanceData) {

	if out, _ := readProcFile(w.Pid, "comm"); out == "java" {
		jvmpath, ok := readLink(fmt.Sprintf("/proc/%d/exe", w.Pid))
		w.JVMPath = jvmpath
		if ok {
			w.JVMInfos, _ = runCmd(jvmpath, "-version")
		}
	}
}

func setConf(w *WaarpInstanceData) {
	w.ConfFile = "not found"

	for i, f := range w.Cmd {
		if f == "org.waarp.openr66.server.R66Server" && i+1 < len(w.Cmd) {
			w.ConfFile = w.Cmd[i+1]
			break
		}
	}

	if w.ConfFile == "not found" {
		return
	}

	w.Config, _ = readFile(w.ConfFile)
}

func setHostID(w *WaarpInstanceData) {
	var instConfig wxml.Server
	if err := xml.Unmarshal([]byte(w.Config), &instConfig); err != nil {
		return
	}

	w.HostID = instConfig.Identity.HostId
	w.HostIDSSL = instConfig.Identity.SslHostId
}

func setLogs(w *WaarpInstanceData) {
	w.LogFile = "not found"
	var logConf string
	for _, f := range w.Cmd {
		if strings.Contains(f, "-Dlogback.configurationFile=") {
			logging.Debug.Print(f)

			logConf = strings.Split(f, "=")[1]
			break
		}
	}

	if logConf == "" {
		logging.Error.Print("Log configuration not found")
		return
	}

	content, err := os.Open(logConf)
	if err != nil {
		logging.Error.Printf("Error while opening '%s' (%s)",
			logConf, err.Error())
		return
	}

	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "<file>") {
			continue
		}

		w.LogFile = line[strings.Index(line, ">")+1 : strings.LastIndex(line, "<")]
	}

	if w.LogFile == "not found" {
		return
	}

	w.Logs, _ = readFile(w.LogFile)
}
