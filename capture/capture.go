package capture

import (
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/mycontext"
	"github.com/wangcong15/escape-from-escape-analysis-of-golang/util"
)

// Compile :
func Compile(ctx *mycontext.Context) {
	log.Printf("Handling %s", ctx.ArgMain)
	// compile the main package
	ctx.MainAbsPath = toAbsPath(ctx.ArgMain)
	ctx.PathToLog[ctx.MainAbsPath] = getGcLog(ctx.MainAbsPath)
	// compile other related package in same project
	depsArr := getDeps(ctx.ArgMain, ctx.ArgPkg)
	var tempAbsPath string
	for _, v := range depsArr {
		if !(strings.Contains(v, "vendor") || strings.Contains(v, "thrift_gen")) {
			tempAbsPath = toAbsPath(v)
			ctx.PathToLog[tempAbsPath] = getGcLog(tempAbsPath)
		}
	}
}

// LogAnalysis :
func LogAnalysis(ctx *mycontext.Context) {
	for k := range ctx.PathToLog {
		parse(ctx.PathToLog[k], k, ctx)
	}
}

func parse(logs, pkgAbsPath string, ctx *mycontext.Context) {
	status := 0
	var ptrName string
	var lineNo int
	var filePath string
	reg1, _ := regexp.Compile("^(\\.\\/[a-zA-Z0-9].*?):([0-9]+):.*? moved to heap: (.*?)$")
	reg2, _ := regexp.Compile("^(\\.\\/[a-zA-Z0-9].*?):([0-9]+):.*?\tfrom (.*?) \\(interface-converted\\) at .*$")
	reg3, _ := regexp.Compile("^(\\.\\/[a-zA-Z0-9].*?):([0-9]+):.*?\tfrom (.*?) \\(passed to call\\[argument escapes\\]\\) at .*?:([0-9]+):.*?$")

	for _, v := range strings.Split(logs, "\n") {
		if !strings.Contains(v, "\t") && !strings.Contains(v, "escapes to heap") {
			status = 0
		}
		if ans := reg1.FindAllStringSubmatch(v, -1); len(ans) > 0 {
			status = 1
		}
		if ans := reg2.FindAllStringSubmatch(v, -1); len(ans) > 0 && status == 1 {
			status = 2
			ptrName = ans[0][3]
		}
		if ans := reg3.FindAllStringSubmatch(v, -1); len(ans) > 0 && status == 2 && ans[0][3] == ptrName {
			status = 3
			lineNo, _ = strconv.Atoi(ans[0][4])
			filePath = path.Join(pkgAbsPath, ans[0][1])
			newData := util.EscapeCase{
				LineNo:  lineNo,
				PtrName: ptrName,
			}
			if _, ok := ctx.EscapeCases[filePath]; ok {
				ctx.EscapeCases[filePath] = append(ctx.EscapeCases[filePath], newData)
			} else {
				ctx.EscapeCases[filePath] = []util.EscapeCase{newData}
			}
			ctx.Number++
		}
	}
}
