# http 客户端,支持 windows

编译
set GOARCH=386
go build -a -v -ldflags="-w -s -H windowsgui"
