package verification

import (
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/util"
)

var greenArrow string = "\033[32;1m==>\033[0m"

// irHandler : scan llvm ir, and achieve the verification
func irHandler(irPath string, filePath string, ec util.EscapeCase) bool {
	var result bool = false
	// read file contents and split into lines
	b, err := ioutil.ReadFile(irPath)
	if err != nil {
		return false
	}
	funcBodylineStr := string(b)

	// fetch file id
	fileName := path.Base(filePath)
	var fileID string
	dIFileRegexp := regexp.MustCompile("(!.*?) = !DIFile\\(filename: \\\"" + fileName + "\\\", directory: \\\"\\.\\\"\\)")
	dIFileParams := dIFileRegexp.FindAllStringSubmatch(funcBodylineStr, -1)
	if len(dIFileParams) > 0 {
		fileID = dIFileParams[0][1]
	}

	var varID string
	dILocalVariableRegexp := regexp.MustCompile("(!.*?) = !DILocalVariable\\(name: \\\"" + ec.PtrName + "\\\", .*?, file: " + fileID + ",")
	dILocalVariableParams := dILocalVariableRegexp.FindAllStringSubmatch(funcBodylineStr, -1)
	if len(dILocalVariableParams) > 0 {
		varID = dILocalVariableParams[0][1]
	}

	// divide ir into function bodies
	funcID2Body := make(map[string]string)
	defineFuncRegexp := regexp.MustCompile("define .*? (@[^ ]+)\\(.* (\\{[\\s\\S]*?\n\\})")
	defineParams := defineFuncRegexp.FindAllStringSubmatch(funcBodylineStr, -1)
	callValueRegexp := regexp.MustCompile("call void @llvm.dbg.value\\(.*?([^ ]+?), .*?" + varID + ",")
	var newName string
	var storeName string
	var srcFuncName string
	var goFuncName string
	for _, p := range defineParams {
		funcID := p[1]
		funcBody := p[2]
		if goFuncName != "" {
			if goFuncName == funcID {
				if strings.Contains(funcBody, storeName) {
					log.Printf("From Block.%s may ends before Block.%s, thus %s should be saved in heap memory.", srcFuncName, goFuncName, ec.PtrName)
					return false
				} else {
					return true
				}
			} else {
				continue
			}
		}
		funcID2Body[funcID] = funcBody
		if !strings.Contains(funcBody, varID+",") {
			continue
		}
		srcFuncName = p[1]
		lineArr := strings.Split(funcBody, "\n")
		for _, line := range lineArr {
			line = strings.TrimSpace(line)
			callValueParams := callValueRegexp.FindAllStringSubmatch(line, -1)
			if len(callValueParams) > 0 {
				newName = callValueParams[0][1]
				break
			}
		}
		if newName == "" {
			continue
		}
		storeRegexp := regexp.MustCompile("store .*? " + newName + ", .*? ([^ ]+?),")
		goFuncRegexp := regexp.MustCompile("call void @__go_go\\(.*?(@[^ ]+) ")
		for _, line := range lineArr {
			line = strings.TrimSpace(line)
			if storeName == "" {
				storeParams := storeRegexp.FindAllStringSubmatch(line, -1)
				if len(storeParams) > 0 {
					storeName = storeParams[0][1]
					continue
				}
			} else {
				goFuncParams := goFuncRegexp.FindAllStringSubmatch(line, -1)
				if len(goFuncParams) > 0 {
					goFuncName = goFuncParams[0][1]
					if v, ok := funcID2Body[goFuncName]; ok {
						if strings.Contains(v, storeName) {
							log.Printf("From Block.%s may ends before Block.%s, thus %s should be saved in heap memory.", srcFuncName, goFuncName, ec.PtrName)
							return false
						} else {
							return true
						}
					}
					break
				}
			}
		}

	}
	return result
}
