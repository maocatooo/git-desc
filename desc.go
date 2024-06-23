package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const descPath = ".git_desc"

func branch(all bool) (string, []string) {
	cmd := "git branch"
	if all {
		cmd = "git branch -a"
	}
	b, _ := execCmd(cmd)
	if b == "" {
		return "", nil
	}
	bs := strings.Split(b, "\n")
	current := ""
	for index, item := range bs {
		branch := strings.TrimSpace(item)
		if strings.HasPrefix(item, "*") {
			branch = strings.Trim(branch, "* ")
			current = branch
		}
		bs[index] = branch
	}
	return current, bs
}

func current() string {
	b, err := execCmd("git branch --show-current")
	if err != nil {
		panic(err)
	}
	return b
}

func create(ls *List, args ...string) {
	if len(args) == 1 {
		b := current()
		ls.Append([]string{b, args[0]})
		ls.Save()
	} else if len(args) > 1 {
		ls.Append([]string{args[0], args[1]})
		ls.Save()
	}
}

func show_local(ls *List) {
	current, localBs := branch(false)
	ls.LS(current, localBs)
}

func get_desc_file_path(dir string) *List {
	filePath := path.Join(dir, descPath)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		return &List{
			filePath: filePath,
			Data:     [][]string{},
		}
	} else if err != nil {
		panic(err)
	}
	bs, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(bs), "\n")
	data := [][]string{}

	for _, line := range lines {
		s := strings.Split(line, " ")
		if len(s) != 2 {
			continue
		}
		data = append(data, s)
	}
	return &List{
		filePath: filePath,
		Data:     data,
	}
}

type List struct {
	filePath string
	Data     [][]string
}

func (l *List) ShowInfo(s string) string {
	for _, item := range l.Data {
		if item[0] == s {
			fmt.Println(item[0], "\t", item[1])
		}
	}
	return ""
}

func (l *List) Append(s []string) {
	for _, item := range l.Data {
		if item[0] == s[0] {
			item[1] = s[1]
			return
		}
	}
	l.Data = append(l.Data, s)
}

func (l *List) Save() {
	file, err := os.OpenFile(l.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)

	}
	defer file.Close()
	for _, item := range l.Data {
		_, err = file.WriteString(fmt.Sprint(item[0], " ", item[1], "\n"))
		if err != nil {
			panic(err)
		}
	}
}

func (l *List) roll(current string, all []string, f func(bool, int, string, string)) {
	maxLen := 0
	mp := map[string]string{}
	for _, item := range l.Data {
		mp[item[0]] = item[1]
	}
	for _, item := range all {
		desc := mp[item]
		if maxLen < len(item) {
			maxLen = len(item)
		}
		defer func(item string) {
			f(item == current, maxLen, item, desc)
		}(item)
	}
	maxLen += 3
}

func (l *List) LS(current string, all []string) {
	l.roll(current, all, func(isCurrent bool, maxLen int, item, desc string) {

		if isCurrent {
			item = " * " + item
			fmt.Printf("\033[32m%-*s\033[0m  %s\n", maxLen, item, desc)
		} else {
			item = "   " + item
			fmt.Printf("%-*s  %s\n", maxLen, item, desc)
		}

	})
}

func (l *List) options(current string, all []string) ([]string, map[string]string) {

	var (
		opts []string
		mp   = map[string]string{}
	)

	l.roll(current, all, func(isCurrent bool, maxLen int, item, desc string) {
		opt := "   " + item + "  " + desc
		if isCurrent {
			opt = " * " + item + "  " + desc
		}
		opts = append(opts, opt)
		mp[opt] = item
	})
	return opts, mp
}
