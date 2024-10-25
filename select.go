package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fzdwx/infinite/components"
	"github.com/fzdwx/infinite/components/selection/multiselect"
	"github.com/fzdwx/infinite/components/selection/singleselect"
	"reflect"
	"unsafe"
)

func NewStartUp(c components.Components, ops ...tea.ProgramOption) *components.StartUp {
	program := tea.NewProgram(c, ops...)
	c.SetProgram(program)
	return &components.StartUp{
		P: program,
	}
}

func Select(c []string, opss ...singleselect.Option) *singleselect.Select {
	inner := components.NewSelection(c)
	startUp := NewStartUp(inner, tea.WithoutCatchPanics())
	var ops []singleselect.Option

	// replace row render
	ops = append(ops, singleselect.WithRowRender(func(cursorSymbol string, hintSymbol string, choice string) string {
		return fmt.Sprintf("%s %s", cursorSymbol, choice)
	}))

	// replace prompt
	ops = append(ops, singleselect.WithPrompt("Please selection your option:"))

	// replace key binding
	ops = append(ops, singleselect.WithKeyBinding(singleselect.DefaultSingleKeyMap()))

	ms := &multiselect.Select{}

	setUnexportedField(ms, "startUp", startUp)
	setUnexportedField(ms, "inner", inner)

	var ss = &singleselect.Select{}
	setUnexportedField(ss, "inner", ms)
	ops = append(ops, opss...)
	ss.Apply(ops...)
	return ss
}

func setUnexportedField(structValue interface{}, fieldName string, newValue interface{}) {
	v := reflect.ValueOf(structValue).Elem() // 获取指针指向的值
	field := v.FieldByName(fieldName)        // 获取字段

	// 使用 unsafe.Pointer 修改字段值
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(newValue))
}
