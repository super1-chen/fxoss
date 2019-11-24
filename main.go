// FeiXiang Data OSS system command-line tools
package main

import (
	"github.com/super1-chen/fxoss/cmd"
)

// depend on doc: https://goreleaser.com/customization/
var version string

func main() {
	cmd.Execute(version)
}
