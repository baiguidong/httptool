// main
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// 地址 角色 ip 底包版本 环境包版本 java包版本

type Foo_v struct {
	name      string
	role      string
	ip        string
	b_version string
	h_version string
	j_version string
}

type FooModel_v struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo_v
}

func NewFooModel_v() *FooModel_v {
	m := new(FooModel_v)
	return m
}

func (m *FooModel_v) RowCount() int {
	return len(m.items)
}

func (m *FooModel_v) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.name

	case 1:
		return item.role

	case 2:
		return item.ip
	case 3:
		return item.b_version
	case 4:
		return item.h_version
	case 5:
		return item.j_version
	}
	return nil
}

// Called by the TableView to sort the model.
func (m *FooModel_v) Sort(col int, order walk.SortOrder) error {
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
			return c(a.name < b.name)

		case 1:
			return c(a.role < b.role)

		case 2:
			return c(a.ip < b.ip)

		case 3:
			return c(a.b_version < b.b_version)
		case 4:
			return c(a.h_version < b.h_version)
		case 5:
			return c(a.j_version < b.j_version)
		default:
			return false
		}
	})

	return m.SorterBase.Sort(col, order)
}

type ver_pack struct {
	Packversion string `json:"packversion"`
	Packtime    string `json:"packtime"`
	Packname    string `json:"packname"`
}

func (my *BaseVer) Export() {
	desfile := fmt.Sprintf("%s.csv", time.Now().Format("版本检测_20060102150405"))
	fp, err := os.OpenFile(desfile, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		walk.MsgBox(my.Mw, "失败", err.Error(), walk.MsgBoxIconQuestion)
		return
	}
	defer fp.Close()
	fp.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(fp)
	w.Write([]string{"地市", "角色", "ip", "基础包版本", "环境包版本", "java包版本"})

	for _, v := range my.model.items {
		data := []string{}
		data = append(data, v.name)
		data = append(data, v.role)
		data = append(data, v.ip)
		data = append(data, v.b_version)
		data = append(data, v.h_version)
		data = append(data, v.j_version)
		w.Write(data)
	}
	w.Flush()
	walk.MsgBox(my.Mw, "成功", fmt.Sprintf("文件:%s", desfile), walk.MsgBoxIconInformation)
}

type Sqlinfo struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Userid     string `json:"userid"`
	Password   string `json:"password"`
	Dbname     string `json:"dbname"`
	PasswordEn string `json:"password_en"`
}
type mysql_get struct {
	Code int     `json:"code"`
	Val  Sqlinfo `json:"val"`
}

func get_constr(str string) string {
	if strings.Contains(str, "@") {
		return str
	}
	get_url := fmt.Sprintf("http://%s:8124/mysql/get", str)
	fmt.Println(get_url)
	res, _ := Http_Get(get_url)
	m := mysql_get{}
	err := json.Unmarshal([]byte(res), &m)
	if err == nil {
		fmt.Println(m.Val)
		return fmt.Sprintf(m.Val.Userid + ":" + m.Val.Password + "@tcp(" + m.Val.Ip + ":" + m.Val.Port + ")")
	}
	return str
}
func (my *BaseVer) Check() {
	go func() {
		f, err := walk.NewFont("宋体", 12, walk.FontBold)
		if err == nil {
			my.tv.SetFont(f)
		}
		//my.Res.SetText("准备检测\n")
		ms, err := sql_query_vec_map_url(1, "", "", "select name,ipstr from base_mysql")
		if err != nil {
			walk.MsgBox(my.Mw, "失败1", err.Error(), walk.MsgBoxIconQuestion)
			//my.Res.AppendText(fmt.Sprintf("%s\n", err.Error()))
			return
		}
		for _, m := range ms {
			//my.Res.AppendText(fmt.Sprintf("%s检测\n", m["name"]))
			ms1, err := sql_query_vec_map_url(0, get_constr(m["ipstr"]), "system", "select GROUP_CONCAT(node_role) role,node_ip ip from nodemanage group by node_ip")
			if err != nil {
				//my.Res.AppendText(fmt.Sprintf("%s\n", err.Error()))
				walk.MsgBox(my.Mw, "失败2", err.Error(), walk.MsgBoxIconQuestion)
				continue
			}
			for _, m1 := range ms1 {
				get_url := fmt.Sprintf("http://%s:8124/base_svr/version_env_base", m1["ip"])
				res, _ := Http_Get(get_url)
				vp := ver_pack{}
				json.Unmarshal([]byte(res), &vp)

				f1 := Foo_v{}
				f1.name = m["name"]
				f1.role = m1["role"]
				f1.ip = m1["ip"]
				f1.b_version = vp.Packversion

				get_url = fmt.Sprintf("http://%s:8124/base_svr/version_bdp", m1["ip"])
				res, _ = Http_Get(get_url)
				json.Unmarshal([]byte(res), &vp)

				f1.h_version = vp.Packversion

				get_url = fmt.Sprintf("http://%s:8124/base_svr/version_java", m1["ip"])
				res, _ = Http_Get(get_url)
				json.Unmarshal([]byte(res), &vp)

				f1.j_version = vp.Packversion

				my.model.items = append(my.model.items, &f1)

				my.model.PublishRowsReset()

				my.model.Sort(my.model.sortColumn, my.model.sortOrder)

				//my.Res.AppendText(fmt.Sprintf("角色:%s\tip:%s\t版本:%s\n", m1["role"], m1["ip"], vp.Packversion))
			}
		}
	}()
}
func run_baseversion(Mw *walk.MainWindow) {
	defer my_recover()
	in := &BaseVer{}
	in.model = NewFooModel_v()
	d := Dialog{
		AssignTo: &in.Mw,
		Title:    "版本检测",
		MinSize:  Size{900, 600},
		MaxSize:  Size{900, 600},
		Layout:   VBox{MarginsZero: true},
		Name:     Getuuid(),
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "检测",
						MaxSize:   Size{100, 30},
						MinSize:   Size{100, 30},
						OnClicked: in.Check,
					},
					PushButton{
						Text:      "导出",
						MaxSize:   Size{100, 30},
						MinSize:   Size{100, 30},
						OnClicked: in.Export,
					},
					Label{
						MinSize: Size{600, 30},
						Text:    "检测各地市基础包 环境包 java包版本,需要tool.db",
					},
				},
			},
			TableView{
				AssignTo:         &in.tv,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: true,
				MultiSelection:   false,
				StretchFactor:    1,

				Columns: []TableViewColumn{
					{Title: "地市", Alignment: AlignFar, Width: 80},
					{Title: "角色", Alignment: AlignFar, Width: 180},
					{Title: "ip", Alignment: AlignFar, Width: 140},
					{Title: "底包版本", Alignment: AlignFar, Width: 160},
					{Title: "环境包版本", Alignment: AlignFar, Width: 160},
					{Title: "java包版本", Alignment: AlignFar, Width: 160},
				},
				Model: in.model,
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", in.tv.SelectedIndexes())
				},
			},
			// TextEdit{
			// 	AssignTo:  &in.Res,
			// 	ReadOnly:  true,
			// 	VScroll:   true,
			// 	MaxLength: 999999999,
			// },
		},
	}
	if err := d.Create(Mw); err != nil {
		return
	}
	//go in.Check()
	(*d.AssignTo).Run()
}

type BaseVer struct {
	Mw    *walk.Dialog
	Res   *walk.TextEdit
	tv    *walk.TableView
	model *FooModel_v
}
