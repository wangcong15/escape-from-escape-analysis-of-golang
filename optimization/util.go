package optimization

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/util"
)

var redBar string = "\033[31;1m---\033[0m"
var greenBar string = "\033[32;1m+++\033[0m"
var greenArrow string = "\033[32;1m==>\033[0m"

func fetchCode(filePath string, v util.EscapeCase) string {
	b, _ := ioutil.ReadFile(filePath)
	rawCode := string(b)
	codeSlice := strings.Split(rawCode, "\n")
	var finalCodeSlice1 string
	vLineNo := v.LineNo - 1

	// concat strings
	if vLineNo-3 < 0 {
		finalCodeSlice1 = strings.Join(codeSlice[:vLineNo], "\n") + "\n"
	} else {
		finalCodeSlice1 = strings.Join(codeSlice[vLineNo-3:vLineNo], "\n") + "\n"
	}
	finalCodeSlice2 := redBar + strconv.Itoa(vLineNo) + codeSlice[vLineNo] + "\n"
	cut1 := fmt.Sprintf("\t%sAddr := uintptr(*unsafe.Pointer(%s))", v.PtrName, v.PtrName)
	cut2 := fmt.Sprintf("(*OBJECTTYPE)(*unsafe.Pointer(%sAddr))", v.PtrName)
	cut3 := strings.ReplaceAll(codeSlice[vLineNo], v.PtrName, cut2)
	finalCodeSlice3 := greenBar + strconv.Itoa(vLineNo) + cut1 + "\n" + greenBar + strconv.Itoa(vLineNo+1) + cut3 + "\n"
	var finalCodeSlice4 string
	if vLineNo+4 > len(codeSlice) {
		finalCodeSlice4 = strings.Join(codeSlice[vLineNo+1:], "\n")
	} else {
		finalCodeSlice4 = strings.Join(codeSlice[vLineNo+1:vLineNo+4], "\n")
	}
	finalCodeSlice := finalCodeSlice1 + finalCodeSlice2 + finalCodeSlice3 + finalCodeSlice4
	fmt.Println("--------------------")
	fmt.Println(finalCodeSlice)
	fmt.Println("--------------------")
	return cut1 + "\n" + cut3
}

func rewrite(filePath string, lineNo int, newCode string) {
	lineNo--
	b, _ := ioutil.ReadFile(filePath)
	dumpFile := filePath + ".escape"
	rawCode := string(b)
	codeSlice := strings.Split(rawCode, "\n")
	codeSlice[lineNo] = newCode
	stringToWrite := strings.Join(codeSlice, "\n")
	ioutil.WriteFile(dumpFile, []byte(stringToWrite), 0644)
	log.Printf(greenArrow+" Optimized code is written to \033[35;4m%v\033[0m", dumpFile)
}
