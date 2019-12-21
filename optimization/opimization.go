package optimization

import (
	"log"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/mycontext"
	"github.com/wangcong15/escape-from-escape-analysis-of-golang/verification"
)

func Optimize(ctx *mycontext.Context, filePath string) {
	for i, v := range ctx.EscapeCases[filePath] {
		ctx.Counter++
		log.Printf(greenArrow+" Case %v/%v", ctx.Counter, ctx.Number)
		log.Printf(greenArrow+" \033[35;4m%s\033[0m", filePath)
		rewrite(filePath, v.LineNo, fetchCode(filePath, v))
		ctx.EscapeCases[filePath][i].IsCorrect = verification.Verification(filePath, v)
	}
}
