package preql

import (
	"os"
	"path"
)

func (pkg Package) Generate() error {
	pth := path.Join(pkg.Dir, "preql.go")
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

	err = templ.Execute(file, pkg)
	if err != nil {
		return err
	}

	return nil
}
