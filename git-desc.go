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
	ls := getDescFilePath(info)
	args := os.Args[1:]
	if len(args) == 0 {
		show_local(ls)
		return
	}
	switch args[0] {
	case "-c", "-l":
		create(ls, args...)
	case "-sw", "--sw", "--ck", "-ck":
		checkout(ls, args[1:]...)
	case "-o":
		open(ls, args[1:]...)
	case "-oc", "-co":
		openCurrent(ls, args[1:]...)
	default:
		help()
	}
}

const HELP = `
The git desc command allows you to generate and manage descriptions for Git branches locally. 
You can create, modify, and open branch descriptions, as well as switch between branches using an interactive branch selection interface

Usage
git desc [-c [<branchname>] <name>] [-l <linkurl>] 
         [-sw | --sw | -ck | --ck [<branchname>]] 
         [-o [<branchname>]]
         [-oc | -co ] 

Generic options
-c [<branchname>] <name>
    Create or modify the description for a branch. If the <branchname> parameter is omitted, the description will be added to the current branch.
    eg: 
        git desc -c "this is current branch"
        git desc -c main "this is the main branch"
-l <linkurl>
    Add a URL to the branch that can be opened in a browser. Use the -o, -co, or -oc command to open this URL.
    eg: 
        git desc -l https://github.com
        git desc -c main "this is the main branch" -l https://github.com
-sw, --sw, -ck, --ck [<branchname>]
    Switch to a different branch using an interactive branch selection interface with arrow keys.
-o [<branchname>]
    Open the URL of the specified branch in a browser. If the <branchname> parameter is omitted, the arrow keys will be used to select a branch.
-oc, -co
    Open the URL of the current branch in a browser.
`

func help() {
	fmt.Print(HELP)
}
