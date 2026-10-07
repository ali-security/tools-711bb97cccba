// Command create_mod writes the proxy-format module zip of a source tree,
// using golang.org/x/mod/zip.CreateFromDir (the proxy's own inclusion rules).
// Usage: create_mod <module-path> <version> <dir> <out.zip>
package main

import (
	"log"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("usage: create_mod <module-path> <version> <dir> <out.zip>")
	}
	f, err := os.Create(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	mv := module.Version{Path: os.Args[1], Version: os.Args[2]}
	if err := zip.CreateFromDir(f, mv, os.Args[3]); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
