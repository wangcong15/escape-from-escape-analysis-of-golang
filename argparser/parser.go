package argparser

import (
	"flag"
	"log"
	"os"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/mycontext"
	"github.com/wangcong15/escape-from-escape-analysis-of-golang/util"
)

// Parse : the terminal parameters
func Parse(ctx *mycontext.Context) {
	// parse user input parameters
	flag.StringVar(&ctx.ArgMain, "m", "", "Specify the Main Package")
	flag.StringVar(&ctx.ArgPkg, "p", "", "Specify the Checking Package")
	flag.Parse()
	ctx.PathToLog = make(map[string]string)
	ctx.EscapeCases = make(map[string][]util.EscapeCase)
	ctx.Counter = 0
	ctx.Number = 0
}

// Checks : validate the parameters
func Checks(ctx *mycontext.Context) {
	if ctx.ArgMain == "" || ctx.ArgPkg == "" {
		log.Println("usage: go-escape -m main-package -p project-in-gopath")
		os.Exit(1)
	}
}
