// main
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func GetLocalip() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
	} else {
		for _, address := range addrs {
			ipnet, ok := address.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.String() != "127.0.0.1" && len(ipnet.IP.String()) > 6 && len(ipnet.IP.String()) < 16 {
					return ipnet.IP.String()
				}
			}
		}
	}
	return "127.0.0.1"
}
func init() {
	m_mysql_map = make(map[string]*sql.DB)
}

func Get_mysql() string {
	return fmt.Sprintf("root:root@tcp(8.8.11.225:3401)")
}

var m_mysql_map map[string]*sql.DB
var m_mysql_map_lock sync.RWMutex

func get_base_db(dbstr string) (*sql.DB, error) {
	m_mysql_map_lock.RLock()
	key := Get_mysql() + dbstr
	if _, ok := m_mysql_map[key]; ok {
		m_mysql_map_lock.RUnlock()
		return m_mysql_map[key], nil
	} else {
		m_mysql_map_lock.RUnlock()
		db, err := sql.Open("mysql", Get_mysql()+"/"+dbstr)
		if err == nil {

			db.SetMaxOpenConns(30)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(time.Second * 30)
			m_mysql_map_lock.Lock()
			m_mysql_map[key] = db
			m_mysql_map_lock.Unlock()
			return db, nil

		} else {
			return nil, err
		}
	}
}
func Sql_exec_base(dbstr, sqlstr string) (bool, error) {
	db, err := get_base_db(dbstr)

	if err != nil {
		return false, err
	}
	_, err = db.Exec(sqlstr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func sql_query_vec_map(dbstr, sqlstr string) ([]map[string]string, error) {
	defer my_recover()

	db, err := get_base_db(dbstr)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	scan := make([]interface{}, len(cols))
	vals := make([]sql.RawBytes, len(cols))
	for i := range vals {
		scan[i] = &vals[i]
	}
	var res []map[string]string
	for rows.Next() {
		rows.Scan(scan...)
		m_p := make(map[string]string)
		for b, v := range vals {
			m_p[cols[b]] = string(v)
		}
		res = append(res, m_p)
	}
	return res, nil
}

var server = "8.8.11.225:3401"

type Base_log struct {
	Mw   *walk.Dialog
	Ywxt *walk.LineEdit
	Name *walk.LineEdit
	Gznr *walk.TextEdit
	Wczt *walk.LineEdit
	Wcsj *walk.LineEdit
	Push *walk.PushButton
}

func (my *Base_log) upload() {
	sqlstr := fmt.Sprintf("insert into zj(day,ip,username,ywxt,gznr,wczt,wcsj) values('%s','%s','%s','%s','%s','%s','%s')",
		time.Now().Format("20060102"), GetLocalip(),
		my.Name.Text(), my.Ywxt.Text(), my.Gznr.Text(), my.Wczt.Text(), my.Wcsj.Text())

	_, err := Sql_exec_base("zongjie", sqlstr)
	if err != nil {
		walk.MsgBox(my.Mw, "失败", err.Error(), walk.MsgBoxIconInformation)
		return
	}
	walk.MsgBox(my.Mw, "成功", "成功", walk.MsgBoxIconInformation)

	fd := fmt.Sprintf("%s/http_task", get_home_dir())
	os.MkdirAll(fd, 0777)
	fname := fmt.Sprintf("%s/name.json", fd)
	ioutil.WriteFile(fname, []byte(my.Name.Text()), 0777)
}
func (my *Base_log) view() {
	show_tables(my.Mw)
}

func Http_Get(str string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*15)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 15))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 15,
		},
	}
	resp, err := client.Get(str)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		} else {
			return string(body), nil
		}
	}
}
func (my *Base_log) down() {
	res, err := Http_Get("http://8.8.11.225:10001/export")
	if err != nil {
		walk.MsgBox(my.Mw, "错误", err.Error(), walk.MsgBoxIconError)
	} else {
		str := fmt.Sprintf("http://8.8.11.225:10001/%s", strings.Replace(res, "/sdzw/", "", -1))
		uri, err := url.ParseRequestURI(str)
		if err != nil {
			walk.MsgBox(my.Mw, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		fname := path.Base(uri.Path)
		resp, err := http.Get(str)
		if err != nil {
			walk.MsgBox(my.Mw, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			walk.MsgBox(my.Mw, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}

		err = ioutil.WriteFile(fname, body, 0777)
		if err != nil {
			walk.MsgBox(my.Mw, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		open_file(fname)
	}
}

func Runcmd(cmd string) ([]byte, error) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd)
	} else {
		c = exec.Command("/bin/sh", "-c", cmd)
	}
	return c.Output()
}
func open_file(name string) {
	cmd := fmt.Sprintf("start %s", name)
	Runcmd(cmd)
	return
}

func get_name() string {
	fd := fmt.Sprintf("%s/http_task", get_home_dir())
	os.MkdirAll(fd, 0777)
	fname := fmt.Sprintf("%s/name.json", fd)
	d, _ := ioutil.ReadFile(fname)
	return string(d)
}
func run_log(Mw *walk.MainWindow) {
	defer my_recover()
	in := &Base_log{}
	f := Font{PointSize: 18}
	Dialog{
		AssignTo: &in.Mw,
		Title:    "每日总结",
		MinSize:  Size{800, 600},
		// MaxSize:  Size{1200, 800},
		Layout: VBox{MarginsZero: true},
		Name:   Getuuid(),
		Font:   f,
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "姓    名",
					},
					LineEdit{
						Text:     get_name(),
						AssignTo: &in.Name,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "业务系统",
					},
					LineEdit{
						Text:     "DAAS平台",
						AssignTo: &in.Ywxt,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "工作内容",
					},
					TextEdit{
						Text:     "",
						AssignTo: &in.Gznr,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "完成状态",
					},
					LineEdit{
						Text:     "已完成",
						AssignTo: &in.Wczt,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "完成时间",
					},
					LineEdit{
						Text:     time.Now().Format("20060102"),
						AssignTo: &in.Wcsj,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "查看",
						OnClicked: in.view,
					},
					PushButton{
						Text:      "下载",
						OnClicked: in.down,
					},
					PushButton{
						Text:      "上报",
						OnClicked: in.upload,
					},
				},
			},
		},
	}.Run(Mw)
}
