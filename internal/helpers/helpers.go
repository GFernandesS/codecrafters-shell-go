package helpers

import "os"

func FileHasExecPermissions(file os.FileInfo) bool {
	if file == nil {
		return false
	}

	return file.Mode().Perm()&0111 != 0
}
