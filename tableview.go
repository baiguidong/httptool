// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sort"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Foo struct {
	username string
	ywxt     string
	gznr     string
	wczt     string
	wcsj     string
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

func (m *FooModel) RowCount() int {
	return len(m.items)
}

func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.username

	case 1:
		return item.ywxt

	case 2:
		return item.gznr

	case 3:
		return item.wczt
	case 4:
		return item.wcsj
	}
	panic("unexpected col")
}

// Called by the TableView to sort the model.
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.username < b.username)

		case 1:
			return c(a.ywxt < b.ywxt)

		case 2:
			return c(a.gznr < b.gznr)

		case 3:
			return c(a.wczt < b.wczt)
		case 4:
			return c(a.wcsj < b.wcsj)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *FooModel) ResetRows() {
	sqlstr := fmt.Sprintf("select username,ywxt,gznr,wczt,wcsj from zj where ip='%s' and day<='%s' and day>='%s'",
		GetLocalip(), time.Now().Format("20060102"), time.Now().AddDate(0, 0, -7).Format("20060102"))
	ms, err := sql_query_vec_map("zongjie", sqlstr)
	if err == nil {
		for _, m1 := range ms {
			f1 := Foo{}
			f1.username = m1["username"]
			f1.ywxt = m1["ywxt"]
			f1.gznr = m1["gznr"]
			f1.wczt = m1["wczt"]
			f1.wcsj = m1["wcsj"]
			m.items = append(m.items, &f1)
		}
	}
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

func show_tables(mw *walk.Dialog) {
	model := NewFooModel()
	var tv *walk.TableView
	f := Font{PointSize: 16}
	d := Dialog{
		Title:   "总结列表",
		Font:    f,
		Layout:  VBox{MarginsZero: false, SpacingZero: false, Alignment: AlignHFarVFar, Spacing: 30},
		Size:    Size{1000, 600},
		MinSize: Size{1000, 600},
		Children: []Widget{
			TableView{
				AssignTo:         &tv,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: true,
				MultiSelection:   false,
				StretchFactor:    1,

				Columns: []TableViewColumn{
					{Title: "姓名", Alignment: AlignFar, Width: 100},
					{Title: "业务系统", Alignment: AlignFar, Width: 200},
					{Title: "工作内容", Alignment: AlignFar, Width: 360},
					{Title: "完成状态", Alignment: AlignFar, Width: 150},
					{Title: "完成时间", Alignment: AlignFar, Width: 150},
				},
				Model: model,
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
				},
			},
		},
	}
	d.FixedSize = true
	d.Run(mw)
}
