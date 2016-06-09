package analyze

import (
	"io"
	"logging"
	"os"
)

// Run launches the collect of data about Waarp instances.
func Run(output, hostid string) {
	var writer io.Writer

	switch output {
	case "", "-":
		writer = os.Stdout
	default:
		fd, err := os.Create(output)
		if err != nil {
			logging.Error.Printf("Cannot open output file: %s", err.Error())
			return
		}
		defer fd.Close()
		writer = fd
	}

	report := NewReport()
	report.Instances.DetectRunning()
	report.Collect()

	if err := WriteReport(writer, report); err != nil {
		logging.Error.Print(err.Error())
	}
}
