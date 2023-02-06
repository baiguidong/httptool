package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime/debug"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/sftp"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
)

func my_recover() {
	if err := recover(); err != nil {
		fmt.Println(string(debug.Stack()))
	}
}
func Postdata(str string, data url.Values, head []string) string {
	var res *http.Response
	var err error
	if strings.HasPrefix(str, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				ClientAuth:         tls.NoClientCert,
				InsecureSkipVerify: true,
			},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("POST", str, strings.NewReader(data.Encode()))
		if err != nil {
			return err.Error()
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//header
		for _, h1 := range head {
			h2 := strings.Split(h1, "=")
			if len(h2) == 2 {
				req.Header.Set(h2[0], h2[1])
			}
		}
		res, err = client.Do(req)
		defer client.CloseIdleConnections()
	} else {
		req, err := http.NewRequest("POST", str, strings.NewReader(data.Encode()))
		if err != nil {
			return err.Error()
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//header
		for _, h1 := range head {
			h2 := strings.Split(h1, "=")
			if len(h2) == 2 {
				req.Header.Set(h2[0], h2[1])
			}
		}
		res, err = http.DefaultClient.Do(req)
	}

	if err != nil {
		return err.Error()
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err.Error()
		}
		return string(body)
	}
	return ""
}
func get_home_dir() string {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir
	} else {
		return ""
	}
}
func PostBody(str string, contentType, data string, head []string) string {
	var res *http.Response
	var err error
	if strings.HasPrefix(str, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				ClientAuth:         tls.NoClientCert,
				InsecureSkipVerify: true,
			},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("POST", str, strings.NewReader(data))
		if err != nil {
			return err.Error()
		}
		req.Header.Set("Content-Type", contentType)
		//header
		for _, h1 := range head {
			h2 := strings.Split(h1, "=")
			if len(h2) == 2 {
				req.Header.Set(h2[0], h2[1])
			}
		}

		res, err = client.Do(req)
		client.CloseIdleConnections()
	} else {
		req, err := http.NewRequest("POST", str, strings.NewReader(data))
		if err != nil {
			return err.Error()
		}
		req.Header.Set("Content-Type", contentType)
		//header
		for _, h1 := range head {
			h2 := strings.Split(h1, "=")
			if len(h2) == 2 {
				req.Header.Set(h2[0], h2[1])
			}
		}
		res, err = http.DefaultClient.Do(req)
	}

	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)
}
func inet_ntoa(ipnr int64) string {
	var bytes_u [4]byte
	bytes_u[0] = byte(ipnr & 0xFF)
	bytes_u[1] = byte((ipnr >> 8) & 0xFF)
	bytes_u[2] = byte((ipnr >> 16) & 0xFF)
	bytes_u[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes_u[3], bytes_u[2], bytes_u[1], bytes_u[0]).String()
}

func json_format(s string) string {
	var a interface{}
	err := json.Unmarshal([]byte(s), &a)
	if err != nil {
		return err.Error()
	} else {
		d, err := json.MarshalIndent(&a, "\r\n", "\t")
		if err != nil {
			return err.Error()
		} else {
			return string(d)
		}
	}
}
func Isexist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}
func Connect_ssh(ip, port, user, passwd string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},

		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func Runcmd_ssh(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("b:" + err.Error())
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}

//	cmd := []string{}
//	cmd = append(cmd, "cd /sdzw/ibp")
//	cmd = append(cmd, "pwd")
//	cmd = append(cmd, "du -sh /sdzw/")
func Runcmd_ssh_mul(client *ssh.Client, cmd []string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		return "", err
	}
	err = session.Shell()
	if err != nil {
		return "", err
	}
	for _, cmd_1 := range cmd {
		stdin.Write([]byte(cmd_1 + "\n"))
	}
	stdin.Write([]byte("exit\n"))
	session.Wait()
	d, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
func Sendfile_sftp(client *ssh.Client, local, remote string) (bool, error) {
	sftpclient, err := sftp.NewClient(client)
	if err != nil {
		fmt.Println("b:" + err.Error())
		return false, err
	}
	defer sftpclient.Close()

	srcf, err := os.Open(local)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	defer srcf.Close()

	desf, err := sftpclient.Create(remote)
	if err != nil {
		fmt.Println(":" + err.Error())
		return false, err
	}
	defer desf.Close()
	_, err = io.Copy(desf, srcf)
	if err != nil {
		fmt.Println(":" + err.Error())
		return false, err
	}
	return true, nil
}
func DownfileSftp(client *ssh.Client, local, remote string) (bool, error) {
	sftpclient, err := sftp.NewClient(client)
	if err != nil {
		fmt.Println("b:" + err.Error())
		return false, err
	}
	defer sftpclient.Close()
	//检查目录
	localdir := filepath.Dir(local)
	os.MkdirAll(localdir, 0777)

	desf, err := os.Create(local)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	defer desf.Close()

	srcf, err := sftpclient.Open(remote)
	if err != nil {
		fmt.Println(":" + err.Error())
		return false, err
	}
	defer srcf.Close()
	_, err = io.Copy(desf, srcf)
	if err != nil {
		fmt.Println(":" + err.Error())
		return false, err
	}
	return true, nil
}
func List_file(path string, reg string, max int) ([]string, error) {

	filelist := []string{}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return filelist, err
	}
	sep := "/"
	for _, fi := range dir {
		if fi.IsDir() {
			s, err := List_file(path+sep+fi.Name(), reg, max)
			if err == nil {
				for _, a := range s {
					filelist = append(filelist, a)
					if len(filelist) >= max {
						return filelist, nil
					}
				}
			}
		} else {

			if strings.HasSuffix(fi.Name(), ".tmp") || strings.HasSuffix(fi.Name(), ".temp") || strings.HasSuffix(fi.Name(), ".dealling") {
			} else {
				filelist = append(filelist, path+sep+fi.Name())
			}

			if len(filelist) >= max {
				return filelist, nil
			}
		}
	}
	return filelist, nil
}

func Getuuid() string {
	u1 := uuid.NewV4()
	return fmt.Sprintf("%s", u1)
}

// 根据url 数据库 获取连接方式
func get_mysql_db(dburl, dbstr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dburl+"/"+dbstr)
	if err == nil {
		return db, nil
	} else {
		return nil, err
	}
}
func get_sqlite_db() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "tool.db")
	if err == nil {
		return db, nil
	} else {
		return nil, err
	}
}
func Sql_exec_mysql_url(t int, dburl, dbstr, sqlstr string) (bool, error) {
	var db *sql.DB
	var err error
	if t == 0 {
		db, err = get_mysql_db(dburl, dbstr)
	} else {
		db, err = get_sqlite_db()
	}
	if err != nil {
		return false, err
	}
	defer db.Close()
	_, err = db.Exec(sqlstr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func sql_query_vec_map_url(t int, dburl, dbstr, sqlstr string) ([]map[string]string, error) {
	defer my_recover()
	var db *sql.DB
	var err error
	if t == 0 {
		db, err = get_mysql_db(dburl, dbstr)
	} else {
		db, err = get_sqlite_db()
	}

	if err != nil {
		return nil, err
	}
	defer db.Close()
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
