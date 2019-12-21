package report

import (
	"log"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/mycontext"
)

func Dump(ctx *mycontext.Context) {
	if len(ctx.EscapeCases) == 0 {
		return
	}
	log.Println("-------------------------")
	var i int = 1
	for p := range ctx.EscapeCases {
		for _, ec := range ctx.EscapeCases[p] {
			log.Printf("Case.%d: %s, Line.%d, %s, Correctness(%v)", i, p, ec.LineNo, ec.PtrName, ec.IsCorrect)
			i++
		}
	}
	log.Println("-------------------------")
}
