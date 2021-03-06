# Week03 作业题目：

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。



# 代码解答



```go
package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	group, egctx := errgroup.WithContext(context.Background())
	var servErr = make(chan error)
	group.Go(func() error {
		handler := http.NewServeMux()
		handler.HandleFunc("/err", func(writer http.ResponseWriter, request *http.Request) {
			// 模拟请求错误
            servErr <- errors.New("http-req err")
		})
		httpServ := &http.Server{
			Addr:              ":8000",
			Handler:           handler,
		}
		go func(s *http.Server) {
            // 服务请求错误
			servErr <- httpServ.ListenAndServe()
		}(httpServ)
		select {
		case <-servErr:
			return errors.New("http server error")
		case <- egctx.Done()://接收其它goroutine的异常通知
			fmt.Println("close the http-server")
			return httpServ.Close()
		}
	})
	group.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-quit:
			return errors.New("syscall signaling")
		case <-egctx.Done():
			return errors.New("")
		}
	})
	err := group.Wait()
	fmt.Printf("app is over by reason[%v]\n", err)
}

```

