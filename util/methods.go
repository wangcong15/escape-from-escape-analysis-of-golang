package util

import (
	"bytes"
	"os/exec"
	"strings"
)

func ExecShell(s string) (string, string) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var stdout, stderr bytes.Buffer
	var outStr, errStr string

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	outStr = strings.TrimRight(outStr, "]\n")
	errStr = strings.TrimRight(errStr, "]\n")
	return outStr, errStr
}
