package jolang2

import (
	"context"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"io/fs"
	"jolang2/nodetype"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type UnitMap map[string]*Unit

type Project struct {
	*sitter.Parser

	Units       []*Unit
	UnitsByPkg  UnitMap
	UnitsByName UnitMap

	*NameNode
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

	if pkgDecl := unit.Root.FindNodeByType(nodetype.PACKAGE_DECLARATION); pkgDecl != nil && pkgDecl.ChildCount() > 1 {
		unit.Package = pkgDecl.Child(1).Content()
		p.UnitsByPkg[unit.Package] = unit
	}

	if classDecl := unit.Root.FindNodeByType(nodetype.CLASS_DECLARATION); classDecl != nil && classDecl.ChildCount() > 2 {
		unit.Name = classDecl.Child(2).Content()
		p.UnitsByName[unit.AbsName()] = unit
	}

	pkgNamedNode := p.AddChild(unit.Package)
	pkgNamedNode.AddChild(unit.Name)

	return unit, nil
}

func NewProject() *Project {
	parser := sitter.NewParser()
	parser.SetLanguage(java.GetLanguage())
	return &Project{
		Parser:      parser,
		Units:       []*Unit{},
		UnitsByPkg:  UnitMap{},
		UnitsByName: UnitMap{},
		NameNode:    NewNameNode("."),
	}
}
