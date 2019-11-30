// Package shred is a golang library to mimic the functionality of the linux shred command
package shred

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

// Conf is a object containing all choices of the user
type Conf struct {
	Times  int
	Zeros  bool
	Remove bool
}

// Path shreds all files in the location of path
// recursively. If remove is set to true files will be deleted
// after shredding. When a file is shredded its content
// is NOT recoverable so !!USE WITH CAUTION!!
func (conf Conf) Path(path string) error {
	stats, err := os.Stat(path)
	if err != nil {
		return err
	} else if stats.IsDir() {
		err := conf.Dir(path)
		if err != nil {
			return err
		}
	} else {
		err := conf.File(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Dir overwrites every File in the location of path and everything in its subdirectories
func (conf Conf) Dir(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		stats, _ := os.Stat(path)

		if !stats.IsDir() {
			wg.Add(1)
			go conf.File(path)
			wg.Wait()
		}
		return nil
	})
	return err
}

// File overwrites a given File in the location of path
func (conf Conf) File(path string) error {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	size := fileinfo.Size()
	err = conf.WriteRandom(path, size)
	if err != nil {
		return err
	}
	err = conf.WriteZeros(path, size)
	if err != nil {
		return err
	}
	if conf.Remove {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	wg.Done()
	return nil
}

// WriteRandom overwrites a File with random stuff.
// conf.Times specifies how many times the file should be overwritten
func (conf Conf) WriteRandom(path string, size int64) error {
	for i := 0; i < conf.Times; i++ {
		file, err := os.OpenFile(path, os.O_RDWR, 0)
		if err != nil {
			return err
		}
		defer file.Close()
		offset, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		buff := make([]byte, size)
		_, err = rand.Read(buff)
		if err != nil {
			return err
		}
		_, err = file.WriteAt(buff, offset)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

// WriteZeros overwrites a File with zeros if conf.Zeros == true
func (conf Conf) WriteZeros(path string, size int64) error {
	if !conf.Zeros {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	defer file.Close()
	offset, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	buff := make([]byte, size)
	_, err = file.WriteAt(buff, offset)
	return err
}
