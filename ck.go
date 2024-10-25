package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	"github.com/fzdwx/infinite/components/selection/singleselect"
	"github.com/fzdwx/infinite/style"
)

func selection(ls *List, args ...string) *info {
	options, mp := ls.options(branch(false))
	if len(options) == 0 {
		return nil
	}
	if len(args) > 0 {
		for _, item := range mp {
			if item.Branch == args[0] {
				return item
			}
		}
	}

	selectKeymap := singleselect.DefaultSingleKeyMap()
	selectKeymap.Confirm = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	selectKeymap.Choice = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	selectKeymap.Choice = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	// page := paginator.New()
	var selected, err = func() (i int, err error) {
		defer func() {
			if e1 := recover(); e1 != nil {
				err = fmt.Errorf("error: pathspec `input value` did not match any file(s) known to git")
			}
		}()

		selected, err := Select(
			options,
			// singleselect.WithDisableFilter(),
			singleselect.WithDisableOutputResult(),
			// singleselect.WithHiddenPaginator(),
			singleselect.WithDisableHelp(),
			// singleselect.WithPaginator(page),
			singleselect.WithPromptStyle(style.New().Fg(usedColor)),
			singleselect.WithKeyBinding(selectKeymap),
			singleselect.WithPageSize(10),
			singleselect.WithChoiceTextStyle(style.New().Fg(usedColor).Bold()),
		).Display("selection branch or input filter key")
		return selected, err
	}()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if _, ok := mp[options[selected]]; !ok {
		return nil
	}
	b := mp[options[selected]]
	return b
}

func checkout(ls *List, args ...string) {
	b := selection(ls, args...)
	if b != nil {
		out, _ := execCmd("git checkout " + b.Branch)
		fmt.Println(out)
	}
}

func open(ls *List, args ...string) {
	b := selection(ls, args...)
	if b != nil && b.Link != "" {
		openBrowser(b.Link)
	}
}

func openCurrent(ls *List, _ ...string) {
	b := ls.Current()
	if b != nil && b.Link != "" {
		openBrowser(b.Link)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = "open"
		args = append(args, url)
	case "linux":
		cmd = "xdg-open"
		args = append(args, url)
	default:
		fmt.Printf("unsupported platform")
		return
	}

	exec.Command(cmd, args...).Start()
}
