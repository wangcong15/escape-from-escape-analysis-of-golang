package mycontext

import "github.com/wangcong15/escape-from-escape-analysis-of-golang/util"

// Context : global information
type Context struct {
	// argument variables
	ArgMain string
	ArgPkg  string

	// inner data
	Counter     int
	Number      int
	MainAbsPath string
	PathToLog   map[string]string
	EscapeCases map[string][]util.EscapeCase // key: file path
}
