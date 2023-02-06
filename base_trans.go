// main
package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Base_trans struct {
	Mw          *walk.Dialog
	Unix_time   *walk.LineEdit
	Str_time    *walk.LineEdit
	Unix_time_1 *walk.LineEdit
	Str_time_1  *walk.LineEdit
	Unix_ip     *walk.LineEdit
	Str_ip      *walk.LineEdit
	Base64_des  *walk.LineEdit
	Base64_src  *walk.LineEdit
	Base64_src1 *walk.LineEdit
	Base64_des1 *walk.LineEdit
}

func (my *Base_trans) tran_ip() {
	u, err := strconv.ParseInt(my.Unix_ip.Text(), 10, 64)
	if err != nil {
		my.Str_ip.SetText(err.Error())
	} else {
		my.Str_ip.SetText(inet_ntoa(u))
	}
}
func (my *Base_trans) base64_decode() {
	d, err := base64.StdEncoding.DecodeString(my.Base64_des.Text())
	if err != nil {
		my.Base64_src.SetText(err.Error())
	} else {
		my.Base64_src.SetText(string(d))
	}
}
func (my *Base_trans) base64_encode() {
	d := base64.StdEncoding.EncodeToString([]byte(my.Base64_src1.Text()))
	my.Base64_des1.SetText(d)
}

func (my *Base_trans) tm_trans() {
	u, err := strconv.ParseInt(my.Unix_time.Text(), 10, 64)
	if err != nil {
		my.Str_time.SetText(err.Error())
	} else {
		my.Str_time.SetText(time.Unix(u, 0).Format("2006-01-02 15:04:05"))
	}
}
func (my *Base_trans) tm_trans_str_unix() {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		my.Str_time.SetText(err.Error())
	} else {
		tm, err := time.ParseInLocation("2006-01-02 15:04:05", my.Str_time_1.Text(), loc)
		if err != nil {
			my.Unix_time_1.SetText(err.Error())
		} else {
			my.Unix_time_1.SetText(fmt.Sprintf("%d", tm.Unix()))
		}
	}
}

func run_trans(Mw *walk.MainWindow) {
	defer my_recover()
	in := &Base_trans{}
	f := Font{PointSize: 18}
	fb := Font{PointSize: 14}
	Dialog{
		AssignTo:  &in.Mw,
		Title:     "格式转换",
		MinSize:   Size{800, 500},
		MaxSize:   Size{800, 500},
		Layout:    VBox{MarginsZero: true},
		Name:      Getuuid(),
		FixedSize: true,

		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "时间戳",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Unix_time,
						Text:     fmt.Sprintf("%d", time.Now().Unix()),
						Font:     f,
					},
					PushButton{
						Text:      "转换",
						Font:      fb,
						OnClicked: in.tm_trans,
					},
					Label{
						Text: "字符串",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Str_time,
						Text:     "",
						Font:     f,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "字符串",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Str_time_1,
						Text:     time.Now().Format("2006-01-02 15:04:05"),
						Font:     f,
					},
					PushButton{
						Text:      "转换",
						Font:      fb,
						OnClicked: in.tm_trans_str_unix,
					},
					Label{
						Text: "时间戳",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Unix_time_1,
						Text:     "",
						Font:     f,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "整型ip",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Unix_ip,
						Text:     "",
						Font:     f,
					},
					PushButton{
						Text:      "转换",
						Font:      fb,
						OnClicked: in.tran_ip,
					},
					Label{
						Text: "字符串ip",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Str_ip,
						Text:     "",
						Font:     f,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "base64",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Base64_des,
						Text:     "",
						Font:     f,
					},
					PushButton{
						Text:      "base64解码",
						Font:      fb,
						OnClicked: in.base64_decode,
					},
					Label{
						Text: "原串",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Base64_src,
						Text:     "",
						Font:     f,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "原串",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Base64_src1,
						Text:     "",
						Font:     f,
					},
					PushButton{
						Text:      "base64编码",
						Font:      fb,
						OnClicked: in.base64_encode,
					},
					Label{
						Text: "base64",
						Font: f,
					},
					LineEdit{
						AssignTo: &in.Base64_des1,
						Text:     "",
						Font:     f,
					},
				},
			},
		},
	}.Run(Mw)
}
