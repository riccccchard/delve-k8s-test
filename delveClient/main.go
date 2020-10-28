package main

import (
    "fmt"
    "github.com/go-delve/delve/service/api"
    "github.com/go-delve/delve/service/rpc2"
    "net/http"
    "strconv"
    "strings"
    "time"
)


func main(){
    http.HandleFunc("/" , func(w http.ResponseWriter , r *http.Request){
        address := r.Header.Get("address")
        timestr := r.Header.Get("time")
        waitTime , err := strconv.Atoi(timestr)
        if err != nil{
            waitTime = 30
        }
        client := rpc2.NewClient(address)
        client.SetReturnValuesLoadConfig(&api.LoadConfig{})
        if client.IsMulticlient(){
            //如果接受multiclient，则在attach时候会让目标程序继续运行
            status , err := client.Halt()
            if err != nil{
               msg := fmt.Sprintf("Failed to let process halt , error - %s", err.Error())
               fmt.Fprintln(w, msg)
               return
            }
            fmt.Fprintf(w , "%+v\n", status)
        }

        breakpoints , err := client.ListBreakpoints()
        if err != nil{
            msg := fmt.Sprintf("Failed to list break point , error - %s", err.Error())
            fmt.Fprintf(w, msg)
            return
        }

        for _ , breakpoint := range breakpoints{
            _ , err = client.ClearBreakpoint(breakpoint.ID)
            if err != nil{
                msg := fmt.Sprintf("Failed to clear breakpoint %+v , error - %s", breakpoint, err.Error())
                fmt.Fprintln(w, msg)
                return
            }
        }
        fmt.Fprintf(w , "cleared all breakpoints")
        locs, err := client.FindLocation(api.EvalScope{}, "database/sql.(*DB).Query", true)
        if err != nil {
            msg := fmt.Sprintf("Failed to find database/sql.(*DB).Query, error - %s", err.Error())
            fmt.Fprintln(w, msg)
            return
        }

        instructions , err := client.DisassemblePC(api.EvalScope{GoroutineID: -1}, locs[0].PC, api.GNUFlavour)
        if err != nil{
            msg := fmt.Sprintf("Failed to disassemble pc , error - %s", err.Error())
            fmt.Fprintln(w, msg)
            return
        }
        var pos int
        for i , inst := range instructions{
            fmt.Fprintf(w, "%24.24s  %24.24d  %24.24s\n", inst.Loc.File, inst.Loc.Line, inst.Text)
            if inst.DestLoc == nil {
                continue
            }
            if inst.DestLoc.Function == nil {
                continue
            }
            if strings.Contains(inst.DestLoc.Function.Name(), "QueryContext") {
                pos = i
                break
            }
        }

        b0, err := client.CreateBreakpoint(&api.Breakpoint{
            Addr: instructions[pos+1].Loc.PC,
        })
        if err != nil {
            msg := fmt.Sprintf("Failed to create break point , error - %s", err.Error())
            fmt.Fprintf(w, msg)
            return
        }
        fmt.Fprintf(w, "Breakpoints: %+v", b0)
        timeTicker := time.NewTicker( time.Second * time.Duration(waitTime))

        for{
            select {
            case <- timeTicker.C:{
                fmt.Fprintf(w , "time has come.")
                goto timedone
            }
            default:
                status := <- client.Continue()
                if status.CurrentThread != nil {
                    fmt.Fprintf(w, "File&Line : %s:%d , pc : 0x%s\n", status.CurrentThread.File, status.CurrentThread.Line, strconv.FormatUint(status.CurrentThread.PC, 16))
                }

                regs, err := client.ListScopeRegisters(api.EvalScope{
                    GoroutineID: -1,
                }, true)
                if err != nil {
                    msg := fmt.Sprintf("Failed to List Scope Register , error - %s", err.Error())
                    fmt.Fprintf(w , msg)
                    return
                }
                for _, reg := range regs {
                    if reg.Name == "Rsp" {
                        value, err := strconv.ParseInt(strings.Trim(reg.Value, "\""), 0, 64)
                        if err != nil {
                            msg := fmt.Sprintf("Failed to convert value , error - %s", err.Error())
                            fmt.Fprintf(w , msg)
                            return
                        }
                        expr := fmt.Sprintf("*(*error)(%d) = \"database/sql/driver\".ErrBadConn", value+0x48)
                        state, err := client.Call(-1, expr, false)

                        if err != nil {
                            msg := fmt.Sprintf("Failed to call expr, error - %s", err.Error())
                            fmt.Fprintf(w , msg)
                            return
                        }
                        if state.CurrentThread != nil {
                            fmt.Fprintf(w, "File&Line : %s:%d , pc : 0x%s\n", state.CurrentThread.File, state.CurrentThread.Line, strconv.FormatUint(state.CurrentThread.PC, 16))
                        }
                    }
                }
            }
        }
timedone:
        client.Disconnect(true)
        fmt.Fprintf(w, "client done.\n")
    })
    http.ListenAndServe("0.0.0.0:8888", nil)
}