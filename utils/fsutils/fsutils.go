package fsutils

import "os"

func DirMustExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(path, 0777)
}
