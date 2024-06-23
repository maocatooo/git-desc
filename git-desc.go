package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execCmd(s string) (string, error) {
	cmds := strings.Split(s, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func main() {
	info, err := execCmd("git rev-parse --git-dir")
	if err != nil {
		fmt.Println(info)
		return
	}
	ls := get_desc_file_path(info)
	args := os.Args[1:]
	if len(args) == 0 {
		show_local(ls)
		return
	}
	switch args[0] {
	case "-c":
		create(ls, args[1:]...)
	case "--ck", "ck":
		checkout(ls, args[1:]...)
	}
}
