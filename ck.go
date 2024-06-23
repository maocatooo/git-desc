package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	inf "github.com/fzdwx/infinite"
	"github.com/fzdwx/infinite/components/selection/singleselect"
)

func checkout(ls *List, args ...string) {
	current, bs := branch(false)
	if len(bs) == 0 {
		return
	}
	ls.options(current, bs)
	options, mp := ls.options(current, bs)
	selectKeymap := singleselect.DefaultSingleKeyMap()
	selectKeymap.Confirm = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	selectKeymap.Choice = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish select"),
	)
	selected, err := inf.NewSingleSelect(
		options,
		singleselect.WithDisableFilter(),
		singleselect.WithDisableOutputResult(),
		singleselect.WithHiddenPaginator(),
		singleselect.WithDisableHelp(),
		singleselect.WithKeyBinding(selectKeymap),
		singleselect.WithPageSize(len(options)),
		// singleselect.WithChoiceTextStyle(),
	).Display("selection branch")
	if err != nil {
		return
	}
	b := mp[options[selected]]
	out, err := execCmd("git checkout " + b)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

}
