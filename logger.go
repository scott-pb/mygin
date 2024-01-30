package mygin

import (
	"fmt"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func LoggerFunc() HandlerFunc {
	return func(context *Context) {
		// Start timer
		start := time.Now()
		path := context.Request.URL.Path

		// Process request
		context.Next()

		latency := time.Now().Sub(start)

		str := fmt.Sprintf("[MYGIN] %v | %3d| %13v |%s \n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			context.GetStatusCode(),
			latency,
			path,
		)
		fmt.Println(str)

	}
}
