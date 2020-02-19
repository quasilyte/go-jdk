package loader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quasilyte/GopherJRE/jclass"
	"github.com/quasilyte/GopherJRE/jdeps"
	"github.com/quasilyte/GopherJRE/vmdat"
)

func readClassFiles(name string, cfg *Config) ([]*jclass.File, error) {
	var pkgPath string
	var files []os.FileInfo
	fsname := strings.ReplaceAll(name, ".", string(os.PathSeparator))
	for _, cp := range cfg.ClassPath {
		var err error
		pkgPath = filepath.Join(cp, fsname)
		files, err = dirClassFiles(pkgPath)
		if err == nil && len(files) != 0 {
			break
		}
	}
	if len(files) == 0 {
		return nil, errors.New("none of the class paths contained the specified package")
	}
	out := make([]*jclass.File, len(files))
	for i, f := range files {
		var err error
		out[i], err = decodeClassFile(filepath.Join(pkgPath, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("decode %s: %v", f.Name(), err)
		}
	}

	return out, nil
}

func dirClassFiles(name string) ([]os.FileInfo, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	filesOnly := list[:0]
	for _, f := range list {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".class") {
			filesOnly = append(filesOnly, f)
		}
	}
	f.Close()
	return filesOnly, err
}

func decodeClassFile(filename string) (*jclass.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var dec jclass.Decoder
	return dec.Decode(f)
}

// findDependencies returns non-loaded dependencies for the given package.
func findDependencies(st *vmdat.State, files []*jclass.File) []string {
	var nonLoaded []string
	for _, f := range files {
		deps := jdeps.ClassDependencies(f)
		for _, d := range deps {
			if st.FindPackage(d) == nil {
				nonLoaded = append(nonLoaded, d)
			}
		}
	}
	return nonLoaded
}

// FIXME: duplicated from irgen package.
func splitName(full string) (name, pkg string) {
	delim := strings.LastIndexByte(full, '/')
	if delim == -1 {
		return full, ""
	}
	return full[delim+1:], full[:delim]
}
