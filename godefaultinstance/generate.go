package godefaultinstance

import (
	"fmt"
	. "go/ast"
	"go/token"

	"golang.org/x/tools/go/types"
)

func (c Config) GeneratePackage(files []*File, pkg *types.Package, info *types.Info) (*File, error) {
	g := &generator{Config: c, files: files, pkg: pkg, info: info}
	err := g.generate()
	return g.ret, err
}

type generator struct {
	Config
	pkg   *types.Package
	files []*File
	s     *types.Scope
	obj   types.Type
	info  *types.Info
	ret   *File
	pos   token.Pos
}

func (g *generator) generate() error {
	g.s = g.pkg.Scope()

	ty, ok := g.s.Lookup(g.typ).(*types.TypeName)
	if !ok {
		return fmt.Errorf("no type %s in package %s.", g.typ, g.pkg.Name())
	}
	g.obj = ty.Type()
	if g.ptr {
		g.obj = types.NewPointer(g.obj)
	}

	g.ret = &File{
		Name: NewIdent(g.pkg.Name()),
	}
	if g.s.Lookup(g.name) == nil {
		g.generateVar()
	}

	g.generateMethods()

	return nil
}

func (g *generator) generateVar() {
	spec := &TypeSpec{
		Name: NewIdent(g.name),
		Type: NewIdent(g.typ),
	}
	if g.ptr {
		spec.Type = &StarExpr{X: spec.Type}
	}
	g.ret.Decls = append(g.ret.Decls, &GenDecl{
		Doc:   g.comment(fmt.Sprintf("%s is the default instance of %s.", g.name, g.typ)),
		Tok:   token.VAR,
		Specs: []Spec{spec},
	})
}

func (g *generator) generateMethods() {
	ms := g.methodSet(g.obj)
	for i := 0; i < ms.Len(); i++ {
		m := ms.At(i)
		if !m.Obj().Exported() {
			continue
		}
		mname := m.Obj().Name()
		for _, exc := range g.exclude {
			if exc == mname {
				continue
			}
		}

		params := []Expr{}
		mtype := m.Type().(*types.Signature)
		tparams := mtype.Params()
		for j := 0; j < tparams.Len(); j++ {
			tp := tparams.At(j)
			s := tp.Name()
			if j == tparams.Len()-1 && mtype.Variadic() {
				s += "..."
			}
			params = append(params, NewIdent(s))
		}
		g.ret.Decls = append(g.ret.Decls, &FuncDecl{
			Doc:  g.comment(fmt.Sprintf("%[1]s calls the %[1]s method on %[2]s.", mname, g.name)),
			Name: NewIdent(mname),
			Type: g.lookupInAST(mname, m.Indirect()),
			Body: &BlockStmt{List: []Stmt{
				&ReturnStmt{Results: []Expr{&CallExpr{
					Fun: &SelectorExpr{
						X:   NewIdent(g.name),
						Sel: NewIdent(mname),
					},
					Args: params,
				}}},
			}},
		})
	}
}

func (g *generator) methodSet(t types.Type) *types.MethodSet {
	ms := types.NewMethodSet(g.obj)
	return ms
}

func (g *generator) lookupInAST(method string, ptrRecv bool) *FuncType {
	for _, f := range g.files {
		for _, decl := range f.Decls {
			fdecl, ok := decl.(*FuncDecl)
			if !ok || fdecl.Recv == nil {
				continue
			}
			if fdecl.Name.Name != method {
				continue
			}
			typ := fdecl.Recv.List[0].Type
			x, isPtrRecv := typ.(*StarExpr)
			if isPtrRecv {
				if !ptrRecv {
					continue
				}
				typ = x.X
			}
			if typ.(*Ident).Name != g.typ {
				continue
			}
			return fdecl.Type
		}
	}
	return nil
}

func (g *generator) comment(s string) *CommentGroup {
	return &CommentGroup{List: []*Comment{
		{Text: "// " + s},
	}}
}
