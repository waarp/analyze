package analyze

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

var (
	reportTpl *template.Template
)

const reportTplContent = `{{title 1 "Waarp Analyze Report"}}

Generation date:
  {{.GenerationDate}}
Generation command:
  {{.Invocation}}


{{title 2 "System Data"}}

System:
  {{.System.System}}
Distribution:
  {{.System.Distribution}}
Kernel:
  {{.System.Kernel}}
Arch:
  {{.System.Arch}}
CPUs:
  {{.System.CPUs}}
RAM:
  {{.System.RAM}}
System Load:
  {{.System.Load}}

{{title 3 "CPU Infos"}}

::

{{indent "  " .System.CPUInfos}}

{{title 3 "RAM Infos"}}

::

{{indent "  "  .System.RAMInfos}}

{{title 3 "Disks"}}

::

{{indent "  "  .System.Disks}}

{{title 3 "Running processes"}}

::

{{indent "  " .System.RunningProcesses}}

{{title 2 "Waarp Instances"}}

{{range $i, $inst := .Instances}}
{{- printf "Instance %s" $inst.HostID | title 3}}

Type:
  {{$inst.Type}}
Id:
  {{$inst.HostID}} ({{$inst.HostIDSSL}})
Command line:
  {{join " " $inst.Cmd}}
Working directory:
  {{$inst.Cwd}}
Pid:
  {{$inst.Pid}}
User:
  {{$inst.User}} ({{$inst.Uid}})
Group:
  {{$inst.Group}} ({{$inst.Gid}})
Jvm:
  {{$inst.JVMPath}}
{{indent "  " $inst.JVMInfos}}
RSS:
  {{$inst.RSS}} ({{$inst.PRAM}}%)
CPU:
  {{$inst.PCPU}} %
VSZ:
  {{$inst.VSZ}}
Thread count:
  {{$inst.NThreads}}
File descriptor count:
  {{$inst.NFD}}

{{title 4 "Jars used"}}

::

{{indent "  " $inst.Jars}}

{{title 4 "Configuration"}}

::

{{indent "  " $inst.Config}}

{{title 4 "Opened Sockets"}}

::

{{indent "  " $inst.Sockets}}

{{title 4 "Status"}}

::

{{indent "  " $inst.Status}}

{{title 4 "Environment"}}

::

{{indent "  " $inst.Environ}}

{{title 4 "System Limits"}}

::

{{indent "  " $inst.Limits}}

{{title 4 "Logs"}}

From: {{$inst.LogFile}}

::

{{indent "  " $inst.Logs}}
{{end}}
`

func init() {
	funcMap := template.FuncMap{
		"title":  makeTitle,
		"indent": indent,
		"join":   join,
	}
	reportTpl = template.New("system report").Funcs(funcMap)
	reportTpl = template.Must(reportTpl.Parse(reportTplContent))
}

func WriteReport(w io.Writer, rd *ReportData) error {
	return reportTpl.Execute(w, rd)
}

func makeTitle(level uint, title string) string {
	levelChar := []string{"#", "=", "-", "~"}
	return fmt.Sprintf("\n%s\n%s", title, strings.Repeat(levelChar[level-1], len(title)))
}

func indent(mark string, text string) string {
	lines := strings.Split(text, "\n")
	indentedText := make([]string, len(lines))

	for i, line := range lines {
		indentedText[i] = fmt.Sprintf("%s%s", mark, line)
	}

	return strings.Join(indentedText, "\n")
}

func join(sep string, items []string) string {
	return strings.Join(items, sep)
}
