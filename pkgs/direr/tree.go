package direr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Tree is the data struct for files and directories in a defined root.
// FileName, Path, and IsDir are exactly what they sound like. Prev and Next
// are meant only for going from file end node to file end node, so no dir will
// have those filled in. And obviously, only dirs can have Children.
type Tree struct {
	FileName string
	IsDir    bool
	URLPath  string
	Prev     *Tree
	Next     *Tree
	Children []*Tree
}

// MakeTree makes a pointer to a tree based on the relative path to the root
// dir given in pathToRoot and recursively fills out FileName, and IsDir for all nodes.
// (joining pathPadding to pathToRoot to make absolute paths)
func MakeTree(pathToRoot string, extFilter string) (t *Tree, err error) {
	info, err := os.Lstat(pathToRoot)
	t = &Tree{
		FileName: info.Name(),
		IsDir:    info.IsDir(),
	}

	if !t.IsDir && !strings.HasSuffix(t.FileName, extFilter) {
		return nil, nil
	}

	if t.IsDir {
		contents, err := os.ReadDir(pathToRoot)
		if err != nil {
			return nil, err
		}

		for _, c := range contents {
			child, err := MakeTree(filepath.Join(pathToRoot, c.Name()), extFilter)
			if err != nil {
				return nil, err
			}

			if child != nil {
				t.Children = append(t.Children, child)
			}
		}
	}

	if t.IsDir && len(t.Children) < 1 {
		return nil, nil
	}

	return t, nil
}

// TODO fill out prev and next and URLPath
