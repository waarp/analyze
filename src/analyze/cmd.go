package analyze

import (
	"fmt"
	"io/ioutil"
	"logging"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCmd(cmd string, args ...string) (string, bool) {
	logging.Debug.Printf(">>> Running command: %s %s", cmd, strings.Join(args, " "))
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		logging.Debug.Printf("<<< Error: %s", out)
		fullCmd := []string{cmd}
		fullCmd = append(fullCmd, args...)
		msg := fmt.Sprintf("Command '%s' exited with errors (%s)",
			strings.Join(fullCmd, " "), err.Error())
		logging.Warning.Printf(msg)
		return msg, false
	}
	logging.Debug.Printf("<<< Output: %s", out)
	return strings.TrimSpace(string(out)), true
}

func readFile(path string) (string, bool) {
	logging.Debug.Printf(">>> Reading file: %s", path)
	out, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf(">>> Error: Cannot read file '%s' (%s)",
			path, err.Error())
		logging.Warning.Printf(msg)
		return msg, false
	}

	return strings.TrimSpace(string(out)), true
}

func readLink(path string) (string, bool) {
	logging.Debug.Printf(">>> Resolving link: %s", path)
	out, err := os.Readlink(path)
	if err != nil {
		msg := fmt.Sprintf(">>> Error: Cannot read link '%s' (%s)",
			path, err.Error())
		logging.Warning.Printf(msg)
		return msg, false
	}

	return out, true
}

func readProcFile(pid int, file string) (string, bool) {
	logging.Debug.Printf("in readprocfile(%d, %#v)", pid, file)
	out, ok := readFile(fmt.Sprintf("/proc/%d/%s", pid, file))
	logging.Debug.Printf("output: %s", out)
	return out, ok
}

func toInt(str string) int {
	i, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		logging.Warning.Printf("Cannot get number from '%s' (%s)",
			str, err.Error())
		return 0
	}

	return i
}

func toFloat(str string) float64 {
	i, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		logging.Warning.Printf("Cannot get float from '%s' (%s)",
			str, err.Error())
		return 0
	}

	return i
}

func countFilesInDir(path string) int {
	logging.Debug.Printf(">>> Listing files in: %s", path)
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		logging.Error.Printf(">>> Error: Cannot list files of '%s' (%s)",
			path, err.Error())
		return 0
	}
	return len(fis)
}
