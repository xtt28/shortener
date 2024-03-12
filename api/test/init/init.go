// Package init sets the process's working directory to three directories above
// the current one when imported. This is a workaround because tests run from
// the folder they are in and not the project root; therefore, usage of paths
// relative to the project root will cause errors.
package init

import (
	"os"
	"path"
	"runtime"
)

// init sets the process's working directory to three directories above the
// current one.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
