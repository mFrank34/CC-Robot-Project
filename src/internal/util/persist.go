package util

import "os"

// WriteJSONAtomic writes data to path via a temp file + rename,
// so a crash mid-write never leaves a corrupt file.
func WriteJSONAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ReadFileOrEmpty returns the file's contents, or (nil, false) if it
// doesn't exist yet — lets the caller decide what "first run" means.
func ReadFileOrEmpty(path string) ([]byte, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return data, true, nil
}
