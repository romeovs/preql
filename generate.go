package preql

import (
	"bytes"
	"io"
	"os"
	"path"

	"golang.org/x/tools/imports"
)

func (pkg Package) Generate() error {
	pth := path.Join(pkg.Dir, "preql.go")
	buf := new(bytes.Buffer)

	err := templ.Execute(buf, pkg)
	if err != nil {
		return err
	}

	clean, err := imports.Process(pth, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	return rewrite(pth, bytes.NewReader(clean))
}

func rewrite(pth string, content io.Reader) error {
	file, err := os.OpenFile(pth, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}
