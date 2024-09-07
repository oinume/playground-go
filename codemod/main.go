package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fixTests",
	Doc:      "Convert test cases from slice to map and update loop structure.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run: func(pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		// ノードフィルターを定義します。CompositeLitとRangeStmtを追跡します。
		nodeFilter := []ast.Node{
			&ast.CompositeLit{},
			&ast.RangeStmt{},
		}

		inspect.Nodes(nodeFilter, func(n ast.Node, push bool) bool {
			switch n := n.(type) {
			case *ast.CompositeLit:
				// CompositeLitノードを確認し、スライス型かどうかを判定する。
				if compLit, ok := n.Type.(*ast.ArrayType); ok {
					// スライスの要素が構造体型かどうかを確認する。
					if structType, ok := compLit.Elt.(*ast.StructType); ok {
						// 構造体の最初のフィールドの名前を確認する（"name"であることを期待）
						if len(structType.Fields.List) > 0 && structType.Fields.List[0].Names[0].Name == "name" {
							// スライスをマップに変換する修正を提案する。
							mapType := &ast.MapType{
								Key:   ast.NewIdent("string"),
								Value: structType,
							}
							pass.Report(analysis.Diagnostic{
								Pos:     n.Pos(),
								End:     n.End(),
								Message: "Convert test slice to map",
								SuggestedFixes: []analysis.SuggestedFix{
									{
										Message: "Convert to map[string]struct",
										TextEdits: []analysis.TextEdit{
											{
												Pos:     n.Type.Pos(),
												End:     n.Type.End(),
												NewText: []byte(mapTypeString(mapType)),
											},
										},
									},
								},
							})
						}
					}
				}
			case *ast.RangeStmt:
				// RangeStmtノードを確認し、対象がスライス "tests" であるかを確認する。
				if ident, ok := n.X.(*ast.Ident); ok && ident.Name == "tests" {
					// ループ構造をマップに対応するように修正する提案を行う。
					pass.Report(analysis.Diagnostic{
						Pos:     n.Pos(),
						End:     n.End(),
						Message: "Update range loop for map",
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "Use map iteration",
								TextEdits: []analysis.TextEdit{
									{
										Pos:     n.TokPos,
										End:     n.TokPos + token.Pos(len("_, ")),
										NewText: []byte("name, "),
									},
									{
										Pos:     n.Body.Lbrace + 1,
										End:     n.Body.Lbrace + 1,
										NewText: []byte("t.Run(name, "),
									},
									{
										Pos:     n.Body.Rbrace - 1,
										End:     n.Body.Rbrace - 1,
										NewText: []byte("}) "),
									},
								},
							},
						},
					})
				}
			}
			return true
		})

		return nil, nil
	},
}

// サポート関数：MapTypeを文字列に変換する。
func mapTypeString(mt *ast.MapType) string {
	return "map[string]struct {" + "\n" +
		"\targs    args\n" +
		"\twantErr bool\n" +
		"}"
}

func main() {
	singlechecker.Main(Analyzer)
}
