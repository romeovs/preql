package preql

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/tools/go/packages"
)

// Query holds information about methods that have an associated
// sql query.
type Query struct {
	// Name is the name of the function or method.
	Name string

	// QueryString is the associated sql query.
	QueryString string

	// BindArgs are the arguments to the function query.
	BindArgs map[string]interface{}

	// Args are the arguments of the function.
	Args []Arg

	// Type is either "query" or "exec"
	Type string
}

// parseQueries parses all functions that have an associated sql query.
func parseQueries(pkg *packages.Package) []*Query {
	res := make([]*Query, 0)
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			if q := parseFuncQuery(decl, pkg.TypesInfo); q != nil {
				res = append(res, q)
			}
		}
	}

	return res
}

// parseFuncQuery parses a declaration and returns the Query associated with it.
// Returns nil if there is no associated query.
func parseFuncQuery(decl ast.Decl, info *types.Info) *Query {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	sql, typ, ok := parseQuery(fn.Doc.Text())
	if !ok {
		return nil
	}

	return &Query{
		Name:        fn.Name.Name,
		QueryString: sql,
		Args:        parseArgs(fn, info),
		BindArgs:    parseBindArgs(fn, info),
		Type:        typ,
	}
}

var matchSQLQuery = regexp.MustCompile(`(?s)sql:\s+(.*)`)
var matchSQLExec = regexp.MustCompile(`(?s)exec:\s+(.*)`)

// parseQuery parses the func docstring to get the sql query.
func parseQuery(doc string) (string, string, bool) {
	matches := matchSQLQuery.FindStringSubmatch(doc)
	if len(matches) == 2 {
		return matches[1], "query", true
	}

	matches = matchSQLExec.FindStringSubmatch(doc)
	if len(matches) == 2 {
		return matches[1], "exec", true
	}

	return "", "", false
}

// parseBindArgs parses the function args and returns a map of named query arg to variable names.
func parseBindArgs(fn *ast.FuncDecl, info *types.Info) map[string]interface{} {
	res := make(map[string]interface{})

	for _, p := range fn.Type.Params.List {
		if len(p.Names) < 1 {
			continue
		}

		for _, n := range p.Names {
			name := n.Name
			res[name] = name

			fields := fields(info.TypeOf(p.Type))
			for k, v := range fields {
				res[name+"."+k] = name + "." + v
			}
		}
	}

	return res
}

// Arg is a description of a function argument.
type Arg struct {
	// Name is the name of the argument.
	Name string

	// TypeName is the name of the argument type.
	TypeName string
}

// parseArgs parses the arguments of the function.
func parseArgs(fn *ast.FuncDecl, info *types.Info) []Arg {
	res := make([]Arg, 0)
	for _, param := range fn.Type.Params.List {
		if len(param.Names) < 1 {
			continue
		}

		t := typeName(param.Type)
		for _, n := range param.Names {
			res = append(res, Arg{
				Name:     n.Name,
				TypeName: t,
			})
		}
	}

	return res
}

// typeName returns the name of the type is defined in the signature.
func typeName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.StarExpr:
		return "*" + typeName(v.X)
	case *ast.SelectorExpr:
		return typeName(v.X) + "." + v.Sel.Name
	case *ast.Ident:
		return v.Name
	default:
		fmt.Printf("%T\n", v)
		return ""
	}
}

func (q Query) Params() (string, error) {
	bound, params, err := sqlx.BindNamed(sqlx.DOLLAR, q.QueryString, q.BindArgs)
	if err != nil {
		return "", fmt.Errorf("could not create bound query: %s", err)
	}

	strs, err := toStrings(params)
	if err != nil {
		return "", err
	}

	ctx := ctxName(q.Args)
	args := append([]string{ctx, `"` + simplifySQL(bound) + `"`}, strs...)

	return strings.Join(args, ", "), nil
}

func toStrings(x []interface{}) ([]string, error) {
	res := make([]string, len(x))
	for i, v := range x {
		if s, ok := v.(string); ok {
			res[i] = s
			continue
		}
		return nil, fmt.Errorf("Got non-string arg %s", v)
	}

	return res, nil
}

func simplifySQL(query string) string {
	return strings.Join(strings.Fields(query), " ")
}

func ctxName(args []Arg) string {
	for _, arg := range args {
		if arg.TypeName == "context.Context" {
			return arg.Name
		}
	}

	return "context.Background()"
}

func (q Query) Arguments() string {
	args := make([]string, 0, len(q.Args))
	for _, arg := range q.Args {
		args = append(args, arg.Name+" "+arg.TypeName)
	}

	return strings.Join(args, ", ")
}
