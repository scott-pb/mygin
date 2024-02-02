package mygin

import (
	"fmt"
	"net/http"
)

// Recovery 发生错误时，恢复函数，且返回相应错误信息。
func Recovery() HandlerFunc {
	return func(c *Context) {
		// 使用defer延迟执行，以便在函数退出时进行recover
		defer func() {
			if err := recover(); err != nil {
				// 如果发生panic，打印错误信息并返回500 Internal Server Error响应
				fmt.Println(err.(error).Error())
				c.Writer.Write([]byte("Internal Server Error\n"))
				c.status = http.StatusInternalServerError
				c.Abort() // 终止后续中间件的执行
			}
		}()

		c.Next() // 调用下一个中间件或处理函数
	}
}
