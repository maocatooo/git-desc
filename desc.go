package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fzdwx/infinite/color"
	"github.com/fzdwx/infinite/style"
)

const descPath = ".git_desc"

var usedColor = color.Pink

func branch(all bool) (string, int, []string) {
	cmd := "git branch"
	if all {
		cmd = "git branch -a"
	}
	b, _ := execCmd(cmd)
	if b == "" {
		return "", 0, nil
	}
	allBranch := strings.Split(b, "\n")
	currentBranch := ""
	currentIndex := 0
	maxLen := 0
	for index, item := range allBranch {
		branch := strings.TrimSpace(item)
		if strings.HasPrefix(item, "*") {
			branch = strings.Trim(branch, "* ")
			currentBranch = branch
			currentIndex = index
		}
		allBranch[index] = branch
		maxLen = max(maxLen, len(branch))
	}
	// currentBranch into top
	allBranch[0], allBranch[currentIndex] = allBranch[currentIndex], allBranch[0]
	return currentBranch, maxLen, allBranch
}

func current() string {
	b, err := execCmd("git branch --show-current")
	if err != nil {
		panic(err)
	}
	return b
}

func toInfo(args []string) *info {
	var i = &info{}
	c := []string{}
	l := []string{}
	m := 0
	for _, arg := range args {
		switch arg {
		case "-c":
			m = 1
			i.setD = true
		case "-l":
			i.setL = true
			m = 2
		default:
			if m == 1 {
				c = append(c, arg)
			}
			if m == 2 {
				l = append(l, arg)
			}
		}
	}

	if len(c) > 1 {
		i.Branch = c[0]
		i.Desc = strings.Join(c[1:], " ")
	} else {
		i.Branch = current()
		i.Desc = strings.Join(c, " ")
	}
	i.Link = strings.Join(l, " ")
	return i
}

func create(ls *List, args ...string) {
	ls.Append(toInfo(args))
	ls.Save()
}

func show_local(ls *List) {
	ls.LS(branch(false))
}

func getDescFilePath(dir string) *List {
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
			Data:     nil,
		}
	} else if err != nil {
		panic(err)
	}
	bs, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var infos []*info

	_ = json.Unmarshal(bs, &infos)

	return &List{
		filePath: filePath,
		Data:     infos,
	}
}

type info struct {
	Branch string `json:"b"`

	setD bool
	Desc string `json:"d"`

	setL bool
	Link string `json:"l"`
}
type List struct {
	filePath string
	Data     []*info
}

func (l *List) ShowInfo(s string) string {
	for _, item := range l.Data {
		if item.Branch == s {
			fmt.Println(item.Branch, "\t", item.Desc)
		}
	}
	return ""
}
func (l *List) Current() *info {
	c := current()
	for _, item := range l.Data {
		if item.Branch == c {
			return item
		}
	}
	return nil
}

func (l *List) Append(s *info) {
	for _, item := range l.Data {
		if item.Branch == s.Branch {
			if s.setD {
				item.Desc = s.Desc
			}
			if s.setL {
				item.Link = s.Link
			}
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

	bs, _ := json.Marshal(l.Data)
	file.Write(bs)
}

func (l *List) roll(current string, maxLen int, all []string, f func(bool, int, string, *info)) {
	mp := map[string]*info{}
	for _, item := range l.Data {
		mp[item.Branch] = item
	}
	maxLen += 3
	for _, item := range all {
		inf := mp[item]
		if inf == nil {
			inf = &info{
				Branch: item,
			}
		}
		f(item == current, maxLen, item, inf)
	}
}

func (l *List) LS(current string, maxLen int, all []string) {
	maxLen = maxLen + 3
	l.roll(current, maxLen, all, func(isCurrent bool, maxLen int, item string, inf *info) {
		i := "   " + item
		if isCurrent {
			item = " * " + item
			i = style.New().Fg(usedColor).Render(fmt.Sprintf("%-*s", maxLen, item))
		} else {
			i = fmt.Sprintf("%-*s", maxLen, i)
		}
		desc := inf.Desc
		fmt.Printf("%s  %s\n", i, desc)
	})
}

func (l *List) options(current string, maxLen int, all []string) ([]string, map[string]*info) {

	var (
		opts []string
		mp   = map[string]*info{}
	)
	maxLen = maxLen + 3
	l.roll(current, maxLen, all, func(isCurrent bool, maxLen int, item string, inf *info) {
		opt := "   "
		if isCurrent {
			opt = " * "
		}
		desc := inf.Desc
		opt = fmt.Sprintf("%-*s  %s", maxLen, opt+item, desc)
		opts = append(opts, opt)
		mp[opt] = inf
	})
	return opts, mp
}
