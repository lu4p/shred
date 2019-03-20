[![GoDoc](https://godoc.org/github.com/lu4p/shred?status.svg)](https://godoc.org/github.com/lu4p/shred)
[![License](https://img.shields.io/github/license/lu4p/shred.svg)](https://unlicense.org/)
# shred
This is a golang libary to mimic the functionallity of the linux ```shred``` command
## Usage

package main
import (
  "fmt"
  "github.com/lu4p/shred"
)
"golang" ```
func main(){
      conf := shred.Conf{Times: 1, Zeros: true, Remove: false}
			err := shredconf.Path("filename")
}
```
