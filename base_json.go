// main
package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Base_json struct {
	Mw   *walk.Dialog
	Src  *walk.TextEdit
	Des  *walk.TextEdit
	Push *walk.PushButton
}

func (my *Base_json) format() {
	my.Des.SetText(json_format(my.Src.Text()))
}

func run_json(Mw *walk.MainWindow) {
	defer my_recover()
	in := &Base_json{}
	f := Font{PointSize: 18}
	fb := Font{PointSize: 14}
	Dialog{
		AssignTo: &in.Mw,
		Title:    "json格式化",
		MinSize:  Size{1200, 800},
		MaxSize:  Size{1200, 800},
		Layout:   VBox{MarginsZero: true},
		Name:     Getuuid(),

		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text:    "源字符串",
						MaxSize: Size{200, 100},
						Font:    f,
					},
					PushButton{
						Text:      "格式化",
						Font:      fb,
						OnClicked: in.format,
						MaxSize:   Size{200, 100},
					},
					Label{
						Text:    "格式化后字符串",
						MaxSize: Size{200, 100},
						Font:    f,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					TextEdit{
						AssignTo:  &in.Src,
						Font:      f,
						MaxLength: 999999999,
						VScroll:   true,
					},
					TextEdit{
						AssignTo:  &in.Des,
						Font:      f,
						ReadOnly:  true,
						MaxLength: 999999999,
						VScroll:   true,
					},
				},
			},
		},
	}.Run(Mw)
}
