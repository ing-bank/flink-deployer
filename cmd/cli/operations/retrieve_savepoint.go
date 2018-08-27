package operations

import (
	"errors"
	"strings"

	"github.com/spf13/afero"
)

func (o RealOperator) retrieveLatestSavepoint(dir string) (string, error) {
	if strings.HasSuffix(dir, "/") {
		dir = strings.TrimSuffix(dir, "/")
	}

	files, err := afero.ReadDir(o.Filesystem, dir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", errors.New("No savepoints present in directory: " + dir)
	}

	var newestFile string
	var newestTime int64
	for _, f := range files {
		filePath := dir + "/" + f.Name()
		fi, err := o.Filesystem.Stat(filePath)
		if err != nil {
			return "", err
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = filePath
		}
	}

	return newestFile, nil
}
