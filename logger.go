package mygin

import (
	"fmt"
	"net/http"
	"time"
)

const (
	green   = "\033[97;42m" // 绿色
	white   = "\033[90;47m" // 白色
	yellow  = "\033[90;43m" // 黄色
	red     = "\033[97;41m" // 红色
	blue    = "\033[97;44m" // 蓝色
	magenta = "\033[97;45m" // 洋红色
	cyan    = "\033[97;46m" // 青色
	reset   = "\033[0m"     // 重置颜色
)

type LogFormatterParams struct {
}

func Logger() HandlerFunc {
	l := &LogFormatterParams{}
	return l.LoggerFunc()
}

// MethodColor 方法颜色获取
func (l *LogFormatterParams) MethodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// StatusCodeColor 状态颜色获取
func (l *LogFormatterParams) StatusCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// LoggerFunc 记录日志的方法
func (l *LogFormatterParams) LoggerFunc() HandlerFunc {
	return func(context *Context) {
		// 启动计时器
		start := time.Now()

		// 处理请求
		context.Next()

		now := time.Now()
		str := fmt.Sprintf("[MyGIN] %v |%s %3d %s| %13v  |%s %-7s %s %#v\n",
			now.Format("2006/01/02 - 15:04:05"),
			l.StatusCodeColor(context.status), context.status, reset,
			now.Sub(start), //耗时
			l.MethodColor(context.Request.Method), context.Request.Method, reset,
			context.Request.URL.Path,
		)
		fmt.Println(str)
	}
}
