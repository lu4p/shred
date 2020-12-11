// Package shred is a golang library to mimic the functionality of the linux shred command
package shred

import (
	"crypto/rand"
	"os"
	"path/filepath"
)

// Conf is a object containing all choices of the user
type Conf struct {
	Times  int
	Zeros  bool
	Remove bool
}

// Path shreds all files in the location of path
// recursively. If remove is set to true files will be deleted
// after shredding. When a file is shredded its content
// is NOT recoverable so USE WITH CAUTION!!!
func (conf Conf) Path(path string) error {
	stats, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stats.IsDir() {
		return conf.Dir(path)
	}

	return conf.File(path)
}

// Dir overwrites every File in the location of path and everything in its subdirectories
func (conf Conf) Dir(path string) error {
	var chErrors []chan error

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		chErr := make(chan error)
		chErrors = append(chErrors, chErr)
		go func() {
			err := conf.File(path)
			chErr <- err
		}()

		return nil
	}

	if err := filepath.Walk(path, walkFn); err != nil {
		return err
	}

	for _, chErr := range chErrors {
		if err := <-chErr; err != nil {
			return err
		}
	}

	return nil
}

// File overwrites a given File in the location of path
func (conf Conf) File(path string) error {
	for i := 0; i < conf.Times; i++ {
		if err := overwriteFile(path, true); err != nil {
			return err
		}
	}

	if conf.Zeros {
		if err := overwriteFile(path, false); err != nil {
			return err
		}
	}

	if conf.Remove {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

func overwriteFile(path string, random bool) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	buff := make([]byte, info.Size())
	if random {
		if _, err := rand.Read(buff); err != nil {
			return err
		}
	}

	_, err = f.WriteAt(buff, 0)
	return err
}
