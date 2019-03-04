package shred

import (
	"crypto/rand"
	"os"
	"path/filepath"
)

// Conf is a object containing all choices of the user
type Conf struct {
	Filepath string
	Times    int
	Zeros    bool
	Remove   bool
}

// Path shreds all files in the location of Conf.Filepath
// recursively. If remove is set to true files will be deleted
// after shredding. When a file is shredded its content
// is NOT recoverable so !!USE WITH CAUTION!!
func (conf Conf) Path() error {
	stats, err := os.Stat(conf.Filepath)
	if err != nil {
		return err
	} else if stats.IsDir() {
		err := conf.Dir()
		if err != nil {
			return err
		}
	} else {
		err := conf.File()
		if err != nil {
			return err
		}
	}
	return nil
}

// File overwrites a given File in the location of Conf.Filepath
func (conf Conf) File() error {
	fileinfo, err := os.Stat(conf.Filepath)
	if err != nil {
		return err
	}
	size := fileinfo.Size()
	err = conf.WriteRandom(size)
	if err != nil {
		return err
	}
	err = conf.WriteZeros(size)
	if err != nil {
		return err
	}
	if conf.Remove {
		err := os.Remove(conf.Filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Dir overwrites every File in the location of Conf.Filepath and everything in its subdirectories
func (conf Conf) Dir() error {
	err := filepath.Walk(conf.Filepath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		stats, _ := os.Stat(path)

		if stats.IsDir() == false {
			conf.File()
		}
		return nil
	})
	return err
}

// WriteRandom overwrites a File with random stuff.
// conf.Times specifies how many times the file should be overwritten
func (conf Conf) WriteRandom(size int64) error {
	for i := 0; i < conf.Times; i++ {
		file, err := os.OpenFile(conf.Filepath, os.O_RDWR, 0)
		defer file.Close()
		if err != nil {
			return err
		}
		offset, err := file.Seek(0, 0)
		if err != nil {
			return err
		}
		buff := make([]byte, size)
		rand.Read(buff)
		file.WriteAt(buff, offset)
		file.Close()
	}
	return nil
}

// WriteZeros overwrites a File with zeros if conf.Zeros == true
func (conf Conf) WriteZeros(size int64) error {
	if conf.Zeros == false {
		return nil
	}
	file, err := os.OpenFile(conf.Filepath, os.O_RDWR, 0)
	defer file.Close()
	if err != nil {
		return err
	}

	offset, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	buff := make([]byte, size)
	file.WriteAt(buff, offset)
	return nil
}
