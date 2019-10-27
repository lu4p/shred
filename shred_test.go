package shred_test

import (
	"github.com/lu4p/shred"
)

func Example() {
	shredconf := shred.Conf{Times: 1, Zeros: true, Remove: false}
	shredconf.Path("/path/to/dir_or_file")
}

func ExampleConf_Path() {
	shredconf := shred.Conf{Times: 1, Zeros: true, Remove: false}
	shredconf.Path("/path/to/dir_or_file")
}

func ExampleConf_Dir() {
	shredconf := shred.Conf{Times: 1, Zeros: true, Remove: false}
	shredconf.Dir("/path/to/dir")
}

func ExampleConf_File() {
	shredconf := shred.Conf{Times: 1, Zeros: true, Remove: false}
	shredconf.File("/path/to/file")
}
