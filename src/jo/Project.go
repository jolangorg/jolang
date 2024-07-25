package jo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jolangorg/jolang/src/jo/nodetype"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type UnitsMap map[string]*Unit
type UnitsByPkgMap map[string]UnitsMap

type NodesById map[uint]*Node

type Project struct {
	*sitter.Parser

	Units          []*Unit
	UnitsByPkg     UnitsByPkgMap
	UnitsByAbsName UnitsMap

	NodesById
	Declarations NodeListMap
}

func resolvePath(path string) (string, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}

	return filepath.EvalSymlinks(path)
}

const JavaExt = ".java"

func (p *Project) AddSourceDir(dirname string) error {
	dirname, err := resolvePath(dirname)
	if err != nil {
		return err
	}

	if _, err = os.Stat(dirname); err != nil {
		return err
	}

	err = filepath.WalkDir(dirname, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) != JavaExt {
			return nil
		}

		_, err = p.AddSource(path)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) IndexDeclarations() error {
	for _, unit := range p.Units {
		decls := unit.Root.FindNodesByTypeRecursive(
			nodetype.CLASS_DECLARATION,
			nodetype.ENUM_DECLARATION,
			nodetype.INTERFACE_DECLARATION,
			nodetype.CONSTRUCTOR_DECLARATION,
			nodetype.METHOD_DECLARATION,
			nodetype.FIELD_DECLARATION,
		)
		for _, decl := range decls {
			p.Declarations.AddNode(decl.GetAbsName(), decl)
		}
	}
	return nil
}

func (p *Project) AddSource(filename string) (*Unit, error) {
	filename, err := resolvePath(filename)
	if err != nil {
		return nil, err
	}

	unit := &Unit{}
	unit.Project = p

	unit.SourceCode, err = os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	unit.Tree, err = p.ParseCtx(context.Background(), nil, unit.SourceCode)
	if err != nil {
		return nil, err
	}

	p.Units = append(p.Units, unit)

	unit.Root = unit.WrapNode(unit.RootNode())
	mainDecls := unit.Root.FindNodesByType(
		nodetype.CLASS_DECLARATION,
		nodetype.ENUM_DECLARATION,
		nodetype.INTERFACE_DECLARATION,
	)

	if len(mainDecls) > 1 {
		fmt.Println("[WARN] More than two main declarations in file " + filename)
	}

	mainDecl := unit.Root.FindNodeByType(
		nodetype.CLASS_DECLARATION,
		nodetype.ENUM_DECLARATION,
		nodetype.INTERFACE_DECLARATION,
	)
	pkgDecl := unit.Root.FindNodeByType(nodetype.PACKAGE_DECLARATION)

	if pkgDecl != nil && pkgDecl.ChildCount() > 1 {
		unit.Package = pkgDecl.Child(1).Content()
		if _, ok := p.UnitsByPkg[unit.Package]; !ok {
			p.UnitsByPkg[unit.Package] = make(UnitsMap)
		}
	} else {
		return nil, errors.New("PACKAGE_DECLARATION not found in " + filename)
	}

	if mainDecl != nil {
		unit.Name = mainDecl.GetName()
		p.UnitsByAbsName[unit.AbsName()] = unit
		p.UnitsByPkg[unit.Package][unit.Name] = unit
	} else {
		return nil, nil
		//return nil, errors.New("CLASS_DECLARATION not found in " + filename)
	}

	return unit, nil
}

func NewProject() *Project {
	parser := sitter.NewParser()
	parser.SetLanguage(java.GetLanguage())
	return &Project{
		Parser:         parser,
		Units:          []*Unit{},
		UnitsByPkg:     make(UnitsByPkgMap),
		UnitsByAbsName: make(UnitsMap),
		NodesById:      make(NodesById),
		Declarations:   make(NodeListMap),
	}
}
