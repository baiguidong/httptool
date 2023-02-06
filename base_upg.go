// main
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func run_upg(Mw *walk.MainWindow) {
	defer my_recover()
	in := &Install{}
	Dialog{
		AssignTo: &in.Mw,
		Title:    "安装包安装",
		MinSize:  Size{900, 600},
		Layout:   VBox{MarginsZero: true},
		Name:     Getuuid(),
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "ip",
					},
					LineEdit{
						AssignTo: &in.Ip,
						Text:     "7.7.11.160,7.7.16.24,7.7.16.25,7.7.16.26,7.7.16.28,7.7.16.29,7.7.16.11,7.7.16.12,7.7.16.13,7.7.16.14,7.7.16.15,7.7.16.16,7.7.16.17,7.7.16.90,7.7.16.92,3.3.3.2,3.3.3.3,3.3.3.4,3.3.3.5,3.3.3.8,3.3.3.9",
					},
					Label{
						Text: "端口",
					},
					LineEdit{
						Text:     "22",
						AssignTo: &in.Port,
						MaxSize:  Size{100, 100},
					},
					Label{
						Text: "密码",
					},
					LineEdit{
						AssignTo: &in.passwd,
						Text:     "inet-eyed",
						MaxSize:  Size{100, 100},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "安装包",
					},
					LineEdit{
						AssignTo: &in.Path,
					},
					PushButton{
						Text:      "浏览",
						OnClicked: in.select_file,
					},
					PushButton{
						Text:      "安装",
						OnClicked: in.install_run,
					},
				},
			},
			TextEdit{
				AssignTo:  &in.Res,
				ReadOnly:  true,
				VScroll:   true,
				MaxLength: 999999999,
			},
		},
	}.Run(Mw)
}

type Install struct {
	Mw     *walk.Dialog
	Ip     *walk.LineEdit
	passwd *walk.LineEdit
	Port   *walk.LineEdit
	Path   *walk.LineEdit
	Res    *walk.TextEdit
}

func (my *Install) install_run() {
	go my.install()
}
func (my *Install) install() {
	buf := []string{}
	f, err := walk.NewFont("宋体", 12, walk.FontBold)
	if err == nil {
		my.Res.SetFont(f)
	}
	if Isexist(my.Path.Text()) {
		ips := strings.Split(my.Ip.Text(), ",")
		for _, ip := range ips {
			my.Res.SetText("准备安装\n")
			my.Res.AppendText(fmt.Sprintf("安装包(%s)存在\n", my.Path.Text()))
			my.Res.AppendText(fmt.Sprintf("连接服务器(ip:%s)(port:%s)(用户:root)\n", ip, my.Port.Text()))
			client, err := Connect_ssh(ip, my.Port.Text(), "root", my.passwd.Text())
			if err != nil {
				my.Res.AppendText(err.Error() + "\n")
				continue
			}
			defer client.Close()
			my.Res.AppendText(fmt.Sprintf("服务器连接成功\n"))
			my.Res.AppendText(fmt.Sprintf("准备发送安装包\n"))

			_, err = Sendfile_sftp(client, my.Path.Text(), fmt.Sprintf("/opt/base.tar.gz"))
			if err != nil {
				my.Res.AppendText(err.Error() + "\n")
				continue
			} else {
				my.Res.AppendText(fmt.Sprintf("安装包发送成功\n"))
			}
			cmd := []string{}
			cmd = append(cmd, "cd /opt/")
			cmd = append(cmd, "rm -rf /opt/updateservice")
			cmd = append(cmd, "tar -zxvf base.tar.gz")
			cmd = append(cmd, "cd /opt/updateservice/bin")
			cmd = append(cmd, "chmod +x /opt/updateservice/bin/*")
			cmd = append(cmd, "./setup.sh")
			d, err := Runcmd_ssh_mul(client, cmd)
			if err != nil {
				my.Res.AppendText(fmt.Sprintf("安装失败(%s)\n", err.Error()))
				continue
			}
			fmt.Println(string(d))
			v := strings.Split(string(d), "\n")
			for _, v1 := range v {
				my.Res.AppendText(v1 + "\n")
			}
			my.Res.AppendText(fmt.Sprintf("安装成功\n"))
			time.Sleep(1000 * 1000 * 1000)
			buf = append(buf, fmt.Sprintf("(ip:%s)(port:%s)(用户:root)(安装包:%s)\n", ip, my.Port.Text(), my.Path.Text()))
		}
	} else {
		my.Res.AppendText(fmt.Sprintf("文件不存在,请重新选择\n"))
	}
	if len(buf) > 1 {
		my.Res.SetText("<------------------>\n")
		for _, line := range buf {
			my.Res.AppendText(line)
		}
		my.Res.AppendText("全部安装完成")
	}
}

func (my *Install) select_file() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = ""
	dlg.Title = "选择安装包"
	dlg.Filter = "安装包(*.*)"
	ok, err := dlg.ShowOpen(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !ok {
		return
	}
	my.Path.SetText(dlg.FilePath)
}
