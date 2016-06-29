package explorer

import (
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"runtime"
)

var (
	startupSystemLog bool
)

func formatValues(values ...interface{}) string {
	line := ""
	for i, v := range values {
		if i != 0 {
			line += ": "
		}
		switch v.(type) {
		case int:
			line += fmt.Sprintf("%d", v)
		case uint32:
			line += fmt.Sprintf("%d", v)
		default:
			line += fmt.Sprintf("%s", v)
		}
	}
	return line
}

func formatLine(level string, values ...interface{}) string {
	return formatValues(values...)
}

func FatalShutdown(suggestedWorkaround string, values ...interface{}) {
	seelog.Errorf("Suggested workaround:")
	seelog.Errorf(suggestedWorkaround, values...)
	seelog.Flush()
	os.Exit(1)
}

func logOs() {
	hostname, _ := os.Hostname()
	seelog.Tracef("OS: Hostname: %s", hostname)

	wd, _ := os.Getwd()
	seelog.Tracef("OS: Working directory: %s", wd)
	seelog.Tracef("OS: Executor UID: %d", os.Geteuid())
}

func logRuntime() {
	seelog.Tracef("Runtime: Operating system: %s", runtime.GOOS)
	seelog.Tracef("Runtime: Architecture: %s", runtime.GOARCH)
	seelog.Tracef("Runtime: Go Version: %s", runtime.Version())
	seelog.Tracef("Runtime: Num CPUs: %d", runtime.NumCPU())
}

func logEnv() {
	for i, a := range os.Args[1:] {
		seelog.Tracef("Arg[%d]: [%s]", i, a)
	}
	for _, e := range os.Environ() {
		seelog.Tracef("Env: %s", e)
	}
}

func logSystem() {
	if !startupSystemLog {
		logOs()
		logRuntime()
		logEnv()
	}
	startupSystemLog = true
}

// Start the explorer
func Start() {
	logSystem()
}