package main

import (
	"SecKill/data"
	"SecKill/engine"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

const port = 20080
func main() {
	router := engine.SeckillEngine()  //路由跳转都写在这里
	defer data.Close()

	go func() {//可视化性能测试
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9876", nil))
	}()

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		println("Error when running server. " + err.Error())
	}
}

