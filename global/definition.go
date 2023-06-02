package global

/*// #cgo LDFLAGS: -L../lib/ -llibVirtuoso
// #include "../../LibVirtuoso/library.h"*/ // Temporarily Unavailable
import "C"
import (
	"fmt"
	"orchestra/utils"
	"os"
	"runtime/debug"
)

/*
func CoreLibraryVersion() string {
	return C.GoString(C.LibraryVersion())
}*/

var (
	Prompter        = utils.NewPrompter()
	RuntimeSearcher = utils.NewRuntimeSearcher()
)

func CreateDirectoryIfNotExist(dir string) error {
	_, e := os.Stat(dir)
	if e != nil {
		if os.IsNotExist(e) {
			if e := os.Mkdir(dir, os.ModePerm); e != nil {
				fmt.Println(fmt.Sprintf("%v\n%s", e, debug.Stack()))
				return nil
			}
		} else {
			return e
		}
	}
	return nil
}
