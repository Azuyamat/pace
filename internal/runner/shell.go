package runner

import (
	"runtime"
	"strings"
)

type Shell struct {
	globals map[string]string
}

func NewShell(globals map[string]string) *Shell {
	return &Shell{
		globals: globals,
	}
}

func (s *Shell) GetShellCommand() (string, []string) {
	if shell, ok := s.globals["SHELL"]; ok {
		if shellArgs, ok := s.globals["SHELL_ARGS"]; ok {
			args := strings.Fields(shellArgs)
			return shell, args
		}
		return shell, []string{"-c"}
	}

	if runtime.GOOS == "windows" {
		return "powershell.exe", []string{"-Command"}
	}
	return "sh", []string{"-c"}
}
