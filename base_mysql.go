// main
package main

import (
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func (my *Basemysql) Check() {
	f, err := walk.NewFont("宋体", 12, walk.FontBold)
	if err == nil {
		my.Res.SetFont(f)
	}
	my.Res.SetText("准备安装")
	for {
		my.Res.AppendText("准备安装")
		time.Sleep(1 * time.Second)
	}

}
func run_basemysql(Mw *walk.MainWindow) {
	defer my_recover()
	in := &Basemysql{}

	d := Dialog{
		AssignTo: &in.Mw,
		Title:    "安装包安装",
		MinSize:  Size{900, 600},
		Layout:   VBox{MarginsZero: true},
		Name:     Getuuid(),
		Children: []Widget{
			TextEdit{
				AssignTo:  &in.Res,
				ReadOnly:  true,
				VScroll:   true,
				MaxLength: 999999999,
			},
		},
	}
	if err := d.Create(Mw); err != nil {
		return
	}
	go in.Check()
	(*d.AssignTo).Run()
}

type Basemysql struct {
	Mw  *walk.Dialog
	Res *walk.TextEdit
}
