package linter

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/types"
	"log"

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

		// Check if it's a method call
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// Check if the method is "Exec"
		if sel.Sel.Name != "Exec" {
			return
		}

		// Get the type information of the receiver
		recvType := pass.TypesInfo.Types[sel.X].Type
		if recvType == nil {
			return
		}

		// Check if the type is *bun.UpdateQuery
		ptr, ok := recvType.(*types.Pointer)
		if !ok {
			return
		}

		named, ok := ptr.Elem().(*types.Named)
		if !ok {
			return
		}

		if named.Obj().Pkg().Path() != "github.com/uptrace/bun" || named.Obj().Name() != "UpdateQuery" {
			return
		}

		// Report diagnostic with suggested fix
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: "bun.UpdateQuery must wrap with bunwithoptimistic.WithOptimistic",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Wrap with bunwithoptimistic.WithOptimistic",
					TextEdits: []analysis.TextEdit{
						{
							Pos: sel.X.Pos(),
							End: sel.X.End(),
							NewText: func() []byte {
								var buf bytes.Buffer
								x := &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("bunwithoptimistic"),
										Sel: ast.NewIdent("WithOptimistic"),
									},
									Args: []ast.Expr{
										sel.X,
									},
								}
								if err := printer.Fprint(&buf, pass.Fset, x); err != nil {
									log.Fatalf("failed to print AST node: %v", err)
								}
								return buf.Bytes()
							}(),
						},
					},
				},
			},
		})
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
