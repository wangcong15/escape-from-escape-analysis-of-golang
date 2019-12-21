package main

import (
	"github.com/wangcong15/escape-from-escape-analysis-of-golang/argparser"
	ec "github.com/wangcong15/escape-from-escape-analysis-of-golang/capture"
	mycontext "github.com/wangcong15/escape-from-escape-analysis-of-golang/mycontext"
	co "github.com/wangcong15/escape-from-escape-analysis-of-golang/optimization"
	"github.com/wangcong15/escape-from-escape-analysis-of-golang/report"
)

var ctx mycontext.Context

func main() {
	// STEP.0 parsing execution params
	argparser.Parse(&ctx)
	argparser.Checks(&ctx)

	// STEP.1 Escape Capture
	ec.Compile(&ctx)
	ec.LogAnalysis(&ctx)

	for k := range ctx.EscapeCases {
		// STEP.2 Code Optimization & STEP.3 Correctness Verification
		co.Optimize(&ctx, k)
	}

	// STEP.4 Report
	report.Dump(&ctx)
}
