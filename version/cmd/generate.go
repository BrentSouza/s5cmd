// This package will generate version/version.go by running git commands to get current branch and commit hash
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func mustRunGetResult(cmd string, arg ...string) string {
	var buf bytes.Buffer

	c := exec.Command(cmd)
	c.Args = append(c.Args, arg...)
	c.Stdout = &buf
	c.Stderr = os.Stderr
	err := c.Run()

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return strings.Trim(buf.String(), "\n\r ")
}

func commandToConst(name, command string, args []string) string {
	data := mustRunGetResult(command, args...)

	ret := "\n// " + name + " is the output of \"" + command + " " + strings.Join(args, " ") + "\"\n"
	switch name {
	case "GitSummary":
		ret += "// For release builds, manually edit this to reflect the released version tag.\n"
	case "GitBranch":
		ret += "// For release builds this should be left empty.\n"
	}
	ret += "const " + name + ` = "` + data + `"` + "\n"

	return ret
}

const destinationFile = "version/version.go"

func main() {
	summary := commandToConst("GitSummary", "git", strings.Split("describe --tags --dirty --always", " "))
	branch := commandToConst("GitBranch", "git", strings.Split("symbolic-ref -q --short HEAD", " "))

	timestamp := time.Now().Format(time.UnixDate)

	b := bytes.NewBuffer(nil)
	fmt.Fprint(b, `// Package version is auto-generated using version/cmd/generate.go on non-release builds.
package version

// AUTO-GENERATED. DO NOT EDIT
// `+timestamp+"\n"+summary+branch)
	log.Printf("Writing %s...\n", destinationFile)
	if err := ioutil.WriteFile(destinationFile, b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
