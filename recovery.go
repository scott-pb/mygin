package mygin

import (
	"fmt"
	"net/http"
)

// Recovery 发生错误时，恢复函数，且返回相应错误信息。
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err.(error).Error())
				c.String(http.StatusInternalServerError, "Internal Server Error\n")
				c.Abort()
			}
		}()
		c.Next()
	}
}
