package static

import (
	"errors"
	"os"
	"path"
)

var ErrFileNotFound = errors.New("file not found")

func (s *Files) getLocalPath(relativePath string) (string, error) {
	if isDefaultPath(relativePath) {
		for _, defaultFile := range s.DefaultFiles {
			defaultExits, err := s.localFileExists(relativePath)
			if err != nil {
				return "", err
			}

			if defaultExits {
				return path.Join(s.StaticDir, defaultFile), nil
			}
		}

		return "", nil
	}

	exits, err := s.localFileExists(relativePath)
	if err != nil {
		return "", err
	}

	if exits {
		return path.Join(s.StaticDir, relativePath), nil
	}

	return "", ErrFileNotFound
}

func (s *Files) localFileExists(fileName string) (bool, error) {
	filePath := path.Join(s.StaticDir, fileName)

	var err error

	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
