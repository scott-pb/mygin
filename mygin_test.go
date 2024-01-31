package mygin

import (
	"fmt"
	"net/http"
	"path"
	"testing"
)

func TestMyGin06(t *testing.T) {
	r := Default()
	r.Use()

	//测试发送了错误
	group := r.Group("/api")

	//这个回调会执行
	group.GET("/hello/:name", func(context *Context) {
		name := context.Params.ByName("name")
		arr := []int{1, 3, 5, 7, 9}
		fmt.Println(arr[9])
		context.String(http.StatusOK, path.Join("hello xxxx", name, "!\n"))
	})

	group.GET("/hello2/:name", func(context *Context) {
		name := context.Params.ByName("name")
		context.String(http.StatusOK, path.Join("hello2 ", name, "!\n"))
	})

	// 启动服务器并监听端口
	r.Run(":8088")
}
