package irequals

import (
	"strings"

	"go/ast"
	"go/types"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	register.Plugin("irequals", New)
}

type plugin struct{}

var _ register.LinterPlugin = (*plugin)(nil)

func New(settings any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

func (f *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{{
		Name: "irequals",
		Doc:  "ensure IR structs in plugin packages implement an Equals method",
		Run:  f.run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}}, nil
}

func (f *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func (f *plugin) run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "internal/kgateway/extensions2/plugins") {
		return nil, nil
	}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder([]ast.Node{(*ast.TypeSpec)(nil)}, func(n ast.Node) {
		ts := n.(*ast.TypeSpec)
		if !strings.HasSuffix(ts.Name.Name, "IR") {
			return
		}
		if _, ok := ts.Type.(*ast.StructType); !ok {
			return
		}
		obj := pass.TypesInfo.Defs[ts.Name]
		named, ok := obj.(*types.TypeName)
		if !ok {
			return
		}
		typ, ok := named.Type().(*types.Named)
		if !ok {
			return
		}
		ms := types.NewMethodSet(types.NewPointer(typ))
		if ms.Lookup(pass.Pkg, "Equals") == nil {
			pass.Reportf(ts.Pos(), "IR struct %s must define an Equals method", ts.Name.Name)
		}
	})
	return nil, nil
}
