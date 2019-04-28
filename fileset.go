package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func NewFileSet(sysRoot string) *FileSet {
	return &FileSet{
		data:    make(map[string][]byte),
		sysRoot: sysRoot,
	}
}

type FileSet struct {
	// path to contents
	data    map[string][]byte
	sysRoot string
}

func (fs *FileSet) AddFile(path string, contents []byte) error {
	if _, ok := fs.data[path]; ok {
		return fmt.Errorf("duplicate file in FileSet: %v", path)
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path is not absolute: %v", path)
	}
	fs.data[path] = contents
	return nil
}

func (fs *FileSet) Flush() error {
	for path, contents := range fs.data {
		log.Printf("would write to path %q: %s", filepath.Join(fs.sysRoot, path), string(contents))
		if err := ioutil.WriteFile(filepath.Join(fs.sysRoot, path), contents, 0666); err != nil {
			return fmt.Errorf("failed writing file: %v", err)
		}
	}
	return nil
}
