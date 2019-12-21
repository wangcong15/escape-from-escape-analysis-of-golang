package verification

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/wangcong15/escape-from-escape-analysis-of-golang/util"
)

// Verification : validate the code optimization, whether make damage to memory safety
func Verification(filePath string, ec util.EscapeCase) bool {
	if checkAST(filePath, ec) {
		return true
	} else if checkLLVMIR(filePath, ec) {
		return true
	}
	return false
}

// Check sync or asyn call of the function, if sync, reply YES!!
func checkAST(filePath string, ec util.EscapeCase) bool {
	// generate abstract syntax tree
	b, _ := ioutil.ReadFile(filePath)
	rawCode := string(b)
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "", rawCode, parser.ParseComments)
	var result bool = true

	// search
	ast.Inspect(f, func(n1 ast.Node) bool {
		// Try to find a block statement
		if ret, ok := n1.(*ast.BlockStmt); ok {
			// filter the block
			if fset.Position(ret.Lbrace).Line <= ec.LineNo && fset.Position(ret.Rbrace).Line >= ec.LineNo {
				ast.Inspect(f, func(n2 ast.Node) bool {
					if ret2, ok := n2.(*ast.GoStmt); ok {
						ret3 := ret2.Call
						if fset.Position(ret3.Lparen).Line <= ec.LineNo && fset.Position(ret3.Rparen).Line >= ec.LineNo {
							result = false
							return false
						}
					}
					return true
				})
			}
			return true
		}
		return true
	})
	if result {
		log.Printf("AST-based verification: \033[32;1mPASSED\033[0m")
	} else {
		log.Printf("AST-based verification: \033[31;1mFAILED\033[0m. Starting LLVMIR-based verification")
	}
	return result
}

func checkLLVMIR(filePath string, ec util.EscapeCase) bool {
	var result bool = false
	// Run Go2IR to generate LLVM IR
	log.Println("[1] Cleaning caches ......")
	script := "rm -rf ~/.cache/go-build"
	util.ExecShell(script)
	log.Println("[2] Generating LLVM IR ......")
	script = fmt.Sprintf("Go2IR -p %s -o /tmp/goescape-cache/", filepath.Dir(filePath))
	util.ExecShell(script)
	log.Println("[3] Analyzing ......")
	result = irHandler("/tmp/goescape-cache/1.ll", filePath, ec)

	// return result
	if result {
		log.Printf("LLVMIR-based verification: \033[32;1mPASSED\033[0m.")
		log.Printf(greenArrow + " Optimization in this case is \033[32;1mVALID\033[0m.")
	} else {
		log.Printf("LLVMIR-based verification: \033[31;1mFAILED\033[0m.")
		log.Printf(greenArrow + " Optimization in this case is \033[31;1mINVALID\033[0m.")
	}
	return false
}
