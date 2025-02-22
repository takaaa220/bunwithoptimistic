package linter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "bunwithoptimistic is a linter for bun with optimistic"

// Analyzer is a linter for bun with optimistic
var Analyzer = &analysis.Analyzer{
	Name: "bunwithoptimistic",
	Doc:  "Checks that bun Update queries are wrapped with WithOptimistic",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Check if it's an Exec(ctx) call
		if !isExecCall(call) {
			return
		}

		// Check if the method chain starts with NewUpdate()
		chain := extractMethodChain(call.Fun)
		if !startsWithNewUpdate(chain) {
			return
		}

		// Check if the query is wrapped with WithOptimistic
		if !isWrappedWithOptimistic(chain) {
			pass.Reportf(call.Pos(), "bun Update query must be wrapped with WithOptimistic")
		}
	})

	return nil, nil
}

// isExecCall checks if the call expression is Exec()
func isExecCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		return sel.Sel.Name == "Exec"
	}
	return false
}

// extractMethodChain extracts the method chain from the expression
func extractMethodChain(expr ast.Expr) []string {
	var chain []string
	for {
		sel, ok := expr.(*ast.SelectorExpr)
		if !ok {
			break
		}
		chain = append([]string{sel.Sel.Name}, chain...)
		expr = sel.X
	}
	return chain
}

// startsWithNewUpdate checks if the method chain starts with NewUpdate()
func startsWithNewUpdate(chain []string) bool {
	for _, method := range chain {
		if method == "NewUpdate" {
			return true
		}
	}
	return false
}

// isWrappedWithOptimistic checks if the query is wrapped with WithOptimistic
func isWrappedWithOptimistic(chain []string) bool {
	for _, method := range chain {
		if method == "WithOptimistic" {
			return true
		}
	}
	return false
}
