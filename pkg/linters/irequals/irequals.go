package irequals

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

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
		Name:     "irequals",
		Doc:      "ensure structs implement a proper Equals method with field-by-field comparison",
		Run:      f.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}}, nil
}

func (f *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

// fieldInfo tracks information about a struct field
type fieldInfo struct {
	name     string
	typ      types.Type
	position token.Pos
}

// methodInfo tracks information about an Equals method
type methodInfo struct {
	receiverType *types.Named
	fields       []fieldInfo
	pos          token.Pos
}

func (f *plugin) run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "internal/kgateway/extensions2/plugins") {
		return nil, nil
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// First pass: collect all struct types and their fields
	structFields := make(map[*types.Named][]fieldInfo)
	inspect.Preorder([]ast.Node{(*ast.TypeSpec)(nil)}, func(n ast.Node) {
		ts := n.(*ast.TypeSpec)
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

		// Collect fields
		st := ts.Type.(*ast.StructType)
		var fields []fieldInfo
		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				fieldObj := pass.TypesInfo.Defs[name]
				if fieldObj == nil {
					continue
				}
				fields = append(fields, fieldInfo{
					name:     name.Name,
					typ:      fieldObj.Type(),
					position: name.Pos(),
				})
			}
		}
		structFields[typ] = fields
	})

	// Second pass: analyze Equals methods
	inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		fd := n.(*ast.FuncDecl)
		if fd.Name.Name != "Equals" {
			return
		}

		// Check method signature
		if len(fd.Type.Params.List) != 1 {
			pass.Reportf(fd.Pos(), "Equals method must take exactly one parameter")
			return
		}
		if len(fd.Type.Results.List) != 1 {
			pass.Reportf(fd.Pos(), "Equals method must return exactly one result")
			return
		}
		if !isBoolType(pass.TypesInfo.TypeOf(fd.Type.Results.List[0].Type)) {
			pass.Reportf(fd.Pos(), "Equals method must return bool")
			return
		}

		// Get receiver type
		if fd.Recv == nil || len(fd.Recv.List) != 1 {
			pass.Reportf(fd.Pos(), "Equals method must have exactly one receiver")
			return
		}
		recvType := pass.TypesInfo.TypeOf(fd.Recv.List[0].Type)
		named, ok := recvType.(*types.Named)
		if !ok {
			return
		}

		// Get fields that should be compared
		fields, ok := structFields[named]
		if !ok {
			return
		}

		// Validate that the input parameter is 'any'
		var isAny bool
		paramType := pass.TypesInfo.TypeOf(fd.Type.Params.List[0].Type)
		if named, ok := paramType.(*types.Named); ok {
			// Debug: Print the type name and package
			pass.Reportf(fd.Type.Params.List[0].Pos(), "DEBUG: Found named type %v from package %v", named.Obj().Name(), named.Obj().Pkg().Path())
			if named.Obj().Name() == "any" {
				isAny = true
			}
		} else {
			// Debug: Print the actual type we got
			pass.Reportf(fd.Type.Params.List[0].Pos(), "DEBUG: Got non-named type %T", paramType)
		}

		if !isAny {
			pass.Reportf(fd.Type.Params.List[0].Pos(), "Equals method must take 'any' as parameter type")
			return
		}

		// Check for type assertion since parameter is 'any'
		hasTypeAssert := false
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			if ta, ok := n.(*ast.TypeAssertExpr); ok {
				if ident, ok := ta.X.(*ast.Ident); ok {
					if ident.Name == "in" {
						hasTypeAssert = true
					}
				}
			}
			return true
		})
		if !hasTypeAssert {
			pass.Reportf(fd.Pos(), "Equals method should include type assertion for interface parameter")
		}

		// Analyze method body for field comparisons
		comparedFields := make(map[string]bool)
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			// Look for binary expressions (==, !=)
			be, ok := n.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			// Check if this is a field comparison
			if be.Op == token.EQL || be.Op == token.NEQ {
				// Try to find field names in the comparison
				if sel, ok := be.X.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok {
						if ident.Name == "d" || ident.Name == "this" || ident.Name == "r" {
							comparedFields[sel.Sel.Name] = true
						}
					}
				}
			}
			return true
		})

		// Report missing field comparisons
		for _, field := range fields {
			if !comparedFields[field.name] {
				pass.Reportf(field.position, "field %s is not compared in Equals method", field.name)
			}
		}
	})

	return nil, nil
}

func isBoolType(t types.Type) bool {
	basic, ok := t.(*types.Basic)
	return ok && basic.Kind() == types.Bool
}
