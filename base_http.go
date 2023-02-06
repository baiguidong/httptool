// main
package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Base_http struct {
	Mw          *walk.MainWindow
	StartTm     *walk.LineEdit
	EndTm       *walk.LineEdit
	UseTm       *walk.LineEdit
	Status      *walk.LineEdit
	Recvdata    *walk.LineEdit
	Url_str     *walk.LineEdit
	Send_num    *walk.LineEdit
	Yanchi      *walk.LineEdit
	Bingfa      *walk.LineEdit
	Task_name   *walk.LineEdit
	Task_list   *walk.ComboBox
	Sfile       *walk.CheckBox
	Base64      *walk.CheckBox
	Method      *walk.ComboBox
	Datatype    *walk.ComboBox
	ContentType *walk.ComboBox
	Content     *walk.TextEdit
	Header      *walk.TextEdit
	Res         *walk.TextEdit
	Ts          *walk.Action
	M_task      map[string]http_res
	Run         bool
	Rs          chan string
	StartTm_c   chan string
	EndTm_c     chan string
	UseTm_c     chan string
	Recvdata_c  chan string
	Status_c    chan string
	Run_num     chan int
}

type http_res struct {
	Url_str     string `json:"url,attr"`
	Contenttype int    `json:"contenttype,attr"`
	Method      int    `json:"method,attr"`
	Content     string `json:"content,attr"`
	Datatype    int    `json:"datatype,attr"`
	Taskname    string `json:"taskname,attr"`
	Res         string `json:"res,attr"`
}

//发送一次
func (my *Base_http) send() {
	go my.send_run_r()
}

func (my *Base_http) upg() {
	go run_upg(my.Mw)
}
func (my *Base_http) logrun() {
	go run_log(my.Mw)
}
func (my *Base_http) baseversion() {
	go run_baseversion(my.Mw)
}
func (my *Base_http) getlog() {
	go run_getlog(my.Mw)
}

func (my *Base_http) basemysql() {
	go run_basemysql(my.Mw)
}
func (my *Base_http) basemachine() {
	go run_basemachine(my.Mw)
}

func (my *Base_http) trans() {
	go run_trans(my.Mw)
}
func (my *Base_http) format_json() {
	go run_json(my.Mw)
}

//初始化  channel
func (my *Base_http) init() {
	my.Rs = make(chan string, 128)
	my.StartTm_c = make(chan string, 2)
	my.EndTm_c = make(chan string, 2)
	my.UseTm_c = make(chan string, 2)
	my.Recvdata_c = make(chan string, 2)
	my.Status_c = make(chan string, 2)
	my.Run_num = make(chan int)
	go func() {
		for {
			select {
			case a := <-my.Rs:
				my.Res.AppendText(a)
			case a := <-my.StartTm_c:
				my.StartTm.SetText(a)
			case a := <-my.EndTm_c:
				my.EndTm.SetText(a)
			case a := <-my.UseTm_c:
				my.UseTm.SetText(a)
			case a := <-my.Recvdata_c:
				my.Recvdata.SetText(a)
			case a := <-my.Status_c:
				my.Status.SetText(a)
			}
		}
	}()
}

//批量发送
func (my *Base_http) send_all() {
	num := my.Send_num.Text()
	n, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		my.Res.SetText(err.Error())
	} else {
		bf := my.Bingfa.Text()
		bfn, err := strconv.ParseUint(bf, 10, 64)
		if err != nil {
			my.Res.SetText(err.Error())
		} else {
			my.Res.SetText("开始批量发送:\n")
			my.Status_c <- ""
			my.StartTm_c <- ""
			my.EndTm_c <- ""
			my.Recvdata_c <- ""
			my.UseTm_c <- ""
			go my.send_run_num_bf(n, bfn)
		}
	}
}
func (my *Base_http) json_format_request() {
	my.Content.SetText(json_format(my.Content.Text()))
}
func (my *Base_http) json_format() {
	my.Res.SetText(json_format(my.Res.Text()))
}
func (my *Base_http) save() {
	h := http_res{}
	h.Content = my.Content.Text()
	h.Url_str = my.Url_str.Text()
	h.Contenttype = my.ContentType.CurrentIndex()
	h.Datatype = my.Datatype.CurrentIndex()
	h.Method = my.Method.CurrentIndex()
	h.Taskname = my.Task_name.Text()
	h.Res = my.Res.Text()
	d, err := json.Marshal(&h)
	if err != nil {
		walk.MsgBox(my.Mw, "错误", err.Error(), walk.MsgBoxIconQuestion)
	} else {
		fd := fmt.Sprintf("%s/http_task", get_home_dir())
		os.MkdirAll(fd, 0777)
		fname := fmt.Sprintf("%s/%s.json", fd, h.Taskname)
		ioutil.WriteFile(fname, d, 0777)
		walk.MsgBox(my.Mw, "成功", "保存成功", walk.MsgBoxIconQuestion)
	}
}

func (my *Base_http) load_task() []string {
	my.M_task = make(map[string]http_res)
	fd := fmt.Sprintf("%s/http_task", get_home_dir())
	v, err := List_file(fd, "", 2000)
	ta := []string{"测试任务"}
	if err == nil {
		for _, vf := range v {
			data, err := ioutil.ReadFile(vf)
			if err != nil {
				continue
			}
			u := http_res{}
			err = json.Unmarshal(data, &u)
			if err != nil {
				continue
			}
			my.M_task[u.Taskname] = u
			ta = append(ta, u.Taskname)
		}
	}
	return ta
}

//get_home_dir
func (my *Base_http) DatatypeChange() {
	if my.Datatype.Text() == "BODY" {
		my.ContentType.SetCurrentIndex(2)
	}
	if my.Datatype.Text() == "KEY-VAL" {
		my.ContentType.SetCurrentIndex(0)
	}
}

func (my *Base_http) taskChange() {
	taskname := my.Task_list.Text()
	if _, ok := my.M_task[taskname]; ok {
		v := my.M_task[taskname]
		my.Url_str.SetText(v.Url_str)
		my.Content.SetText(v.Content)
		my.ContentType.SetCurrentIndex(v.Contenttype)
		my.Datatype.SetCurrentIndex(v.Datatype)
		my.Method.SetCurrentIndex(v.Method)
		my.Task_name.SetText(v.Taskname)
		my.Res.SetText(v.Res)
	}
}

func (my *Base_http) send_run_num_bf(n, b uint64) {
	for i := uint64(0); i < b; i++ {
		go my.send_run_num(n, i)
	}
	for i := uint64(0); i < b; i++ {
		<-my.Run_num
	}
	my.Rs <- fmt.Sprintf("批量发送完成")
}

//发送 n 次 带 延迟
func (my *Base_http) send_run_num(n, x uint64) {
	y := my.Yanchi.Text()
	yn, err := strconv.ParseInt(y, 10, 64)
	if err != nil {
		yn = 0
	}
	sennum := 0
	usetm := ""
	s := time.Now().UnixNano() / 1e6
	for i := uint64(0); i < n; i++ {
		sennum += len(my.send_run())
		if i > 0 && i%100 == 0 {
			e := time.Now().UnixNano() / 1e6
			usetm = fmt.Sprintf("%dms", e-s)
			my.Rs <- fmt.Sprintf("线程(%d)发送%d次,接收字节数(%d)耗时(%s)\n", x+1, i, sennum, usetm)
		}
		if yn > 0 {
			time.Sleep(time.Duration(yn) * time.Millisecond)
		}
	}
	e := time.Now().UnixNano() / 1e6
	usetm = fmt.Sprintf("%dms", e-s)
	my.Rs <- fmt.Sprintf("线程(%d)发送完成,共%d次,接收字节数(%d) 耗时(%s)\n", x+1, n, sennum, usetm)
	my.Run_num <- 1
}

//发送  并且返回给客户端
func (my *Base_http) send_run_r() {
	my.Status_c <- "运行中"
	s := time.Now().UnixNano() / 1e6
	my.StartTm_c <- time.Now().Format("15:04:05.999")
	my.Res.SetText("")
	res := my.send_run()
	my.Rs <- res
	my.EndTm_c <- time.Now().Format("15:04:05.999")
	e := time.Now().UnixNano() / 1e6
	my.UseTm_c <- fmt.Sprintf("%dms", e-s)
	my.Recvdata_c <- fmt.Sprintf("%d", len(res))
	my.Status_c <- "完成"
}

//真正发送,不返回给客户端
func (my *Base_http) send_run() string {
	res := ""
	if my.Method.Text() == "GET" {
		var resp *http.Response
		var err error
		hd := my.Header.Text()
		hd = strings.Replace(hd, "\r", "", -1)
		vhd := strings.Split(hd, "\n")
		if strings.HasPrefix(my.Url_str.Text(), "https://") {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					ClientAuth:         tls.NoClientCert,
					InsecureSkipVerify: true,
				},
			}
			client := &http.Client{Transport: tr}

			//resp, err = client.Get(my.Url_str.Text())

			req, err := http.NewRequest("GET", my.Url_str.Text(), nil)
			if err != nil {
				return err.Error()
			}
			//header
			for _, h1 := range vhd {
				h2 := strings.Split(h1, "=")
				if len(h2) == 2 {
					req.Header.Set(h2[0], h2[1])
				}
			}
			resp, err = client.Do(req)

			defer client.CloseIdleConnections()
		} else {
			req, err := http.NewRequest("GET", my.Url_str.Text(), nil)
			if err != nil {
				return err.Error()
			}
			//header
			for _, h1 := range vhd {
				h2 := strings.Split(h1, "=")
				if len(h2) == 2 {
					req.Header.Set(h2[0], h2[1])
				}
			}
			resp, err = http.DefaultClient.Do(req)
			//resp, err = http.Get(my.Url_str.Text())
		}

		if err != nil {
			res = err.Error()
		} else {
			if resp != nil {
				if resp.Body != nil {
					defer resp.Body.Close()
					d, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						res = err.Error()
					} else {
						res = string(d)
					}
				}
			}
		}
	} else if my.Method.Text() == "POST" {
		val := my.Content.Text()
		hd := my.Header.Text()
		hd = strings.Replace(hd, "\r", "", -1)
		vhd := strings.Split(hd, "\n")

		if my.Datatype.Text() == "BODY" {
			res = PostBody(my.Url_str.Text(), my.ContentType.Text(), val, vhd)
		} else {
			val := strings.Replace(val, "\r", "", -1)
			v := strings.Split(val, "\n")
			data := url.Values{}
			for _, v1 := range v {
				v2 := strings.Split(v1, "=")
				if len(v2) == 2 {
					data[v2[0]] = []string{v2[1]}
				}
			}
			res = Postdata(my.Url_str.Text(), data, vhd)
		}
	}
	if my.Base64.Checked() {
		res = base64.StdEncoding.EncodeToString([]byte(res))
	}
	if my.Sfile.Checked() {
		fd := fmt.Sprintf("%s/http_save", get_home_dir())
		os.MkdirAll(fd, 0777)
		fname := fmt.Sprintf("%s/%s.bcp", fd, Getuuid())
		ioutil.WriteFile(fname, []byte(res), 0777)
		res = fmt.Sprintf("结果保存到:%s", fname)
	}
	return res
}

//程序入口
func run_http() {
	in := &Base_http{}
	in.init()
	task := in.load_task()
	f := Font{PointSize: 16}
	fb := Font{PointSize: 14}

	MainWindow{
		AssignTo: &in.Mw,
		Title:    "小小工具",
		Layout:   VBox{MarginsZero: true},
		Font:     f,
		MaxSize:  Size{1080, 960},
		MinSize:  Size{1080, 960},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 10},
				Children: []Widget{
					PushButton{
						Text:          "版本检测",
						Font:          f,
						OnClicked:     in.baseversion,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "拿日志",
						Font:          f,
						OnClicked:     in.getlog,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "mysql检测",
						Font:          f,
						OnClicked:     in.basemysql,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "机器检测",
						Font:          f,
						OnClicked:     in.basemachine,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "每日总结",
						Font:          f,
						OnClicked:     in.logrun,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "常用转换",
						Font:          f,
						OnClicked:     in.trans,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "程序升级",
						Font:          f,
						OnClicked:     in.upg,
						StretchFactor: 1,
					},
					PushButton{
						Text:          "JSON格式化",
						Font:          f,
						OnClicked:     in.format_json,
						StretchFactor: 1,
					},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "地址",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.Url_str,
								Text:     "http://7.7.16.15:8123/mysql/get",
								Font:     f,
							},
							Label{
								Text: "方式",
								Font: f,
							},
							ComboBox{
								AssignTo:     &in.Method,
								Editable:     false,
								CurrentIndex: 0,
								Model:        []string{"GET", "POST"},
								Font:         f,
							},
							PushButton{
								Text:      "批量发送",
								Font:      fb,
								OnClicked: in.send_all,
								MaxSize:   Size{100, 20},
							},
							LineEdit{
								AssignTo: &in.Send_num,
								Text:     "10",
								Font:     f,
								MaxSize:  Size{100, 20},
							},
							Label{
								Text: "次",
								Font: f,
							},
							Label{
								Text: "延迟",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.Yanchi,
								Text:     "0",
								Font:     f,
								MaxSize:  Size{50, 20},
							},
							Label{
								Text: "ms",
								Font: f,
							},
							Label{
								Text: "并发",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.Bingfa,
								Text:     "1",
								Font:     f,
								MaxSize:  Size{50, 20},
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							ComboBox{
								OnCurrentIndexChanged: in.DatatypeChange,
								AssignTo:              &in.Datatype,
								Editable:              false,
								CurrentIndex:          0,
								Model:                 []string{"KEY-VAL", "BODY"},
								Font:                  f,
							},
							ComboBox{
								AssignTo:     &in.ContentType,
								Editable:     false,
								CurrentIndex: 0,
								Model: []string{
									"application/x-www-form-urlencoded",
									"multipart/form-data",
									"text/plain",
									"application/json",
									"application/xml",
									"application/soap+xml",
									"text/xml",
									"application/javascript",
								},
								Font: f,
							},
							CheckBox{
								Text:     "结果保存文件",
								Checked:  false,
								AssignTo: &in.Sfile,
								Font:     f,
							},
							CheckBox{
								Text:     "结果base64",
								Checked:  false,
								AssignTo: &in.Base64,
								Font:     f,
							},
							Label{
								Text: "任务选择",
								Font: f,
							},
							ComboBox{
								AssignTo:              &in.Task_list,
								Editable:              false,
								CurrentIndex:          0,
								Model:                 task,
								Font:                  f,
								OnCurrentIndexChanged: in.taskChange,
							},
							PushButton{
								Text:      "发送",
								Font:      fb,
								OnClicked: in.send,
							},
						},
					},
					Composite{
						Layout:  HBox{},
						MaxSize: Size{1080, 80},
						Children: []Widget{
							Label{
								Text: "头部",
								Font: f,
							},
							TextEdit{
								AssignTo:  &in.Header,
								VScroll:   true,
								MaxLength: 999999999,
								Font:      f,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "内容",
								Font: f,
							},
							TextEdit{
								AssignTo:  &in.Content,
								VScroll:   true,
								MaxLength: 999999999,
								Font:      f,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "结果",
								Font: f,
							},
							TextEdit{
								AssignTo:  &in.Res,
								ReadOnly:  true,
								VScroll:   true,
								MaxLength: 999999999,
								Font:      f,
							},
						},
					},

					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "开始时间",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.StartTm,
								Font:     f,
								ReadOnly: true,
							},
							Label{
								Text: "结束时间",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.EndTm,
								Font:     f,
								ReadOnly: true,
							},
							Label{
								Text: "耗时",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.UseTm,
								Font:     f,
								ReadOnly: true,
							},
							Label{
								Text: "状态",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.Status,
								Font:     f,
								ReadOnly: true,
							},
							Label{
								Text: "接收字节数",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.Recvdata,
								Font:     f,
								ReadOnly: true,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "任务名称",
								Font: f,
							},
							LineEdit{
								AssignTo: &in.Task_name,
								Text:     "测试任务",
								Font:     f,
							},
							PushButton{
								Text:      "任务保存",
								Font:      fb,
								OnClicked: in.save,
							},
							PushButton{
								Text:      "请求JSON格式化",
								Font:      fb,
								OnClicked: in.json_format_request,
							},
							PushButton{
								Text:      "结果JSON格式化",
								Font:      fb,
								OnClicked: in.json_format,
							},
						},
					},
				},
			},
		},
	}.Run()
}
