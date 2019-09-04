package preql

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

var (
	cfg = &packages.Config{Mode: packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax}
)

// Package contains the information
type Package struct {
	ScannableTypes []*ScannableType
	Queries        []*Query
}

// Load loads the package at pkgpath and returns the relevant types and functions for preql.
func Load(pkgpath string) (*Package, error) {
	pkgs, err := packages.Load(cfg, pkgpath)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("Expected exactly one package at %s but got %v", pkgpath, len(pkgs))
	}

	return &Package{
		ScannableTypes: parseScannableTypes(pkgs[0]),
		Queries:        parseQueries(pkgs[0]),
	}, nil
}
