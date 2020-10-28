package main

import (
	"fmt"
	"git.garena.com/shopee/loan-service/airpay_backend/public/common/log"
	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/debugger"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/go-delve/delve/service/rpccommon"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

/*
	根据pid启动attach server，做attach操作并监听信号
*/
var (
    wg 	*sync.WaitGroup
)
type DelveServer struct{
    server 		service.Server
    //接受停止信号
    disconnectCH	chan struct{}
    //监听端口
    address             string
    attachPid           int
    //attach时长
    duration            time.Duration
    //是否接受多次连接
    acceptMulti         bool
}
func(ds *DelveServer) GetAddress () string{
    return ds.address
}
func(ds *DelveServer) SetAddress (address string) {
    if address == ""{
	address = "localhost:1234"
    }
    ds.address = address
}
func(ds *DelveServer) GetAttachPid() int{
    return ds.attachPid
}
//根据信息初始化server，address为server监听端口，
//acceptMulti表示server是否支持多次client连接，为false表示如果一个client disconnect，那么server也会退出；
//为true表示client disconnect 而server不会退出，如果需要停止attach，需要client手动调用detach
//如果acceptMulti为true，默认设置attach target的时候让target继续运行
func (ds *DelveServer) InitServer(attachPid int , acceptMulti bool, address string, duration time.Duration)error{
    log.Infof("[DelveServer.InitServer]initing delve server with pid : %d , listen address : %s", attachPid , address)

    listener , err := net.Listen("tcp", address)

    if err != nil{
	log.Errorf("[DelveServer.InitServer]Failed to listen to address : %s , error : %s", address , err.Error())
	return err
    }
    ds.SetAddress(address)
    ds.attachPid = attachPid

    ds.disconnectCH = make(chan struct{})

    ds.duration = duration

    ds.acceptMulti = acceptMulti

    workingDir := "."
    config :=  &service.Config{
	Listener:               listener,
	ProcessArgs:            []string{},
	AcceptMulti:            acceptMulti,
	APIVersion:             2,
	CheckLocalConnUser:     false,
	DisconnectChan:         ds.disconnectCH,
	Debugger: debugger.Config{
	    AttachPid:            attachPid,
	    WorkingDir:           workingDir,
	    Backend:              "default",
	    CoreFile:             "",
	    Foreground:           true,
	    Packages:             nil,
	    BuildFlags:           "",
	    ExecuteKind:          debugger.ExecutingOther,
	    DebugInfoDirectories: nil,
	    CheckGoVersion:       true,
	    TTY:                  "",
	    Redirects:            [3]string{},      //可以重定向server的I/O信息
	},
    }
    ds.server = rpccommon.NewServer(config)
    return nil
}
//启动一个serer去attach 目标进程，同时监听停止信号
func (ds *DelveServer) StartServer() error {
    log.Infof("[DelveServer.StartServer]starting delve attach server.")
    if err := ds.server.Run(); err != nil{
	log.Errorf("[DelveServer.StartServer]Failed to run server , error : %s", err.Error())
	return err
    }
    //如果acceptMulti为true，默认在attach时让target继续运行
    if ds.acceptMulti{
	client := rpc2.NewClient(ds.address)
	client.Disconnect(true) //true = continue after close
    }
    log.Infof("[DelveServer.StartServer]delve server attached to pid : %s , listen address : %s", ds.attachPid , ds.address)
    return nil
}
//启动协程，监听停止信号
func (ds *DelveServer) WaitForStopServer() {
    log.Infof("Waiting for Server stop.")
    wg.Add(1)
    go func(){
		//定时任务
		ticker := time.NewTicker(ds.duration)
		//接受系统信号的channel
		ch := make(chan os.Signal, 1)
		signal.Notify(ch , syscall.SIGINT)
		select {
		case <- ticker.C :
			log.Infof("[DelveServer.WaitForStopServer]server stoped by time ticker.")
		case <- ch :
			log.Infof("[DelveServer.WaitForStopServer]server stoped by system signal : SIGINT")
		case <- ds.disconnectCH:
			log.Infof("[DelveServer.WaitForStopServer]server stoped by client.")
		}
		//停止server
		err := ds.server.Stop()
		if err != nil {
			log.Errorf("[DelveServer.WaitForStopServer]failed to stop server : %s", err.Error())
			return
		}
		log.Infof("[DelveServer.WaitForStopServer]server stoped.")
		wg.Done()
    }()
}
func main(){
    wg = new(sync.WaitGroup)
    http.HandleFunc("/", func(w http.ResponseWriter , r *http.Request){
        strpid := r.Header.Get("pid")
        if strpid == ""{
            fmt.Fprintln(w , "[Main]pid can't be empty!")
	    return
	}
        pid  , err := strconv.Atoi( strpid )
        if err != nil{
            fmt.Fprintln(w, "[Main]Failed to conv pid to int")
	    return
	}
        address := r.Header.Get("address")
        myServer := &DelveServer{}
        myServer.SetAddress(address)
        err = myServer.InitServer(pid , true , address , 1500 * time.Second)

        if err != nil{
	    fmt.Fprintln(w, "[Main]Failed to Init Server")
	    return
	}

	err = myServer.StartServer()
	if err != nil{
	    fmt.Fprintln(w, "[Main]Failed to Start Server")
	    return
	}

	myServer.WaitForStopServer()
	fmt.Fprintln(w, "delve done success.")
	wg.Wait()
    })
    http.ListenAndServe("0.0.0.0:3333", nil)
}