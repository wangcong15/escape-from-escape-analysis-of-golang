package capture

import (
	"fmt"
	"strings"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/util"
)

func toAbsPath(path string) (absPath string) {
	var goListScript string
	goListScript = fmt.Sprintf("go list -e -f {{.Dir}} %s", path)
	absPath, _ = util.ExecShell(goListScript)
	return
}

func getGcLog(flagMainPath string) (errStr string) {
	var goBuildScript string // 执行编译并获取GC日志的脚本

	goBuildScript = fmt.Sprintf("cd %s && go1.9.3 build -gcflags=\"-m -m\" *.go", flagMainPath)
	_, errStr = util.ExecShell(goBuildScript)
	return errStr
}

func getDeps(path, proj string) (depsInProj []string) {
	var goListScript, projPrefix, deps string
	var depArr []string

	projPrefix = fmt.Sprintf("%s", proj)
	goListScript = fmt.Sprintf("go list -e -f {{.Deps}} %s", path)

	deps, _ = util.ExecShell(goListScript)
	deps = strings.TrimRight(deps, "]\n")
	deps = strings.TrimLeft(deps, "[")
	depArr = strings.Split(deps, " ")
	for _, v := range depArr {
		if strings.HasPrefix(v, projPrefix) {
			depsInProj = append(depsInProj, v)
		}
	}
	return
}
