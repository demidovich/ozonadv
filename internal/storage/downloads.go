package storage

import (
	"fmt"
	"os"
)

type Downloads struct {
	dir string
}

func NewDownloads(dir string) Downloads {
	return Downloads{
		dir: dir,
	}
}

func (d *Downloads) list() (result []string, err error) {
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

func (d *Downloads) Write(fname string, data string) error {
	return os.WriteFile(d.fpath(fname), []byte(data), 0644)
}

func (d *Downloads) Read(fname string) ([]byte, error) {
	return os.ReadFile(d.fpath(fname))
}

func (d *Downloads) ReadString(fname string) (string, error) {
	b, err := d.Read(fname)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (d *Downloads) Remove(fname string) error {
	return os.Remove(d.fpath(fname))
}

func (d *Downloads) RemoveAll() error {
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
func (d *Downloads) fpath(fname string) string {
	return fmt.Sprintf("%s/%s", d.dir, fname)
}
