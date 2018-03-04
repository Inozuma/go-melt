package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	flag.Parse()

	path := flag.Arg(0)

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	for name, pkg := range pkgs {
		file := ast.MergePackageFiles(pkg, ast.FilterUnassociatedComments|ast.FilterImportDuplicates)
		ast.SortImports(fset, file)

		fmt.Fprintln(buf, "package", name)
		fmt.Fprintln(buf)
		fmt.Fprintln(buf, "import (")
		for _, imp := range file.Imports {
			printer.Fprint(buf, fset, imp)
			fmt.Fprintln(buf)
		}
		fmt.Fprintln(buf, ")")

		fmt.Fprintln(buf)
		for _, decl := range file.Decls {
			if gendecl, ok := decl.(*ast.GenDecl); ok && gendecl.Tok == token.IMPORT {
				// skip imports declarations
				continue
			}
			fmt.Fprintln(buf)
			printer.Fprint(buf, fset, decl)
			fmt.Fprintln(buf)
		}
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(output)
}
