package utils

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var blIsInsideDocker bool = false
var once sync.Once

func IsInsideDocker() bool {
	once.Do(setInsideDocker)
	return blIsInsideDocker
}

func setInsideDocker() {
	blIsInsideDocker = isInsideDocker()
}

func isInsideDocker() bool {
	filePath := "/proc/1/cgroup"
	blInDocker := false

	if _, err := os.Stat(filePath); err == nil {
		// path/to/whatever exists

		data, err := ioutil.ReadFile(filePath)
		if err == nil {
			strData := string(data)
			if strings.Index(strData, "docker") > 0 {
				blInDocker = true
			}
		}

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

	}

	return blInDocker

}
