package storage

import (
	"fmt"
	"os"
)

type downloads struct {
	dir string
}

func newDownloads(dir string) downloads {
	return downloads{
		dir: dir,
	}
}

func (d downloads) list() (result []string, err error) {
	entries, err := os.ReadDir(d.dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if e.Name() == "." || e.Name() == ".." {
			continue
		}
		result = append(result, e.Name())
	}

	return
}

func (d downloads) AbsolutePath(fname string) string {
	return d.dir + "/" + fname
}

func (d downloads) Write(fname string, data []byte) error {
	return os.WriteFile(d.fpath(fname), data, 0644)
}

func (d downloads) Read(fname string) ([]byte, error) {
	return os.ReadFile(d.fpath(fname))
}

func (d downloads) ReadString(fname string) (string, error) {
	b, err := d.Read(fname)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d downloads) Remove(fname string) error {
	return os.Remove(d.fpath(fname))
}

func (d downloads) RemoveAll() error {
	list, err := d.list()
	if err != nil {
		return err
	}

	for _, fname := range list {
		if err := d.Remove(fname); err != nil {
			return err
		}
	}

	return nil
}

// Absolute file path
func (d downloads) fpath(fname string) string {
	return fmt.Sprintf("%s/%s", d.dir, fname)
}
