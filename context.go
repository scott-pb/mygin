package mygin

import (
	"encoding/json"
	"math"
	"net/http"
)

// 定义 表示最大和上下文应中止时的索引值
const abortIndex int8 = math.MaxInt8 >> 1

// Context 封装了一个HTTP请求的上下文
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Params   Params
	index    int8
	handlers HandlersChain
	status   int
}

// Next 执行链中的剩余处理程序。
func (c *Context) Next() {
	c.index++
	//遍历handlers
	for c.index < int8(len(c.handlers)) {
		//真正调用执行handler方法
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort 中断链中剩余处理程序的执行。
func (c *Context) Abort() {
	c.index = abortIndex
}

// IsAborted 如果当前上下文被中止，则返回true。
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// writeContentType 如果尚未设置，则设置Content-Type标头。
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}

}

// Status 设置HTTP响应状态码。
func (c *Context) Status(code int) {
	c.status = code
	c.Writer.WriteHeader(code)
}

// GetStatusCode  返回HTTP响应状态码。
func (c *Context) GetStatusCode() int {
	return c.status
}

// JSON 将值序列化为JSON并将其写入响应。
func (c *Context) JSON(v interface{}) error {
	writeContentType(c.Writer, []string{"application/json; charset=utf-8"})
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(v)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
	return err
}

// Html 将字符串以HTML形式写入响应。
func (c *Context) Html(v string) error {
	writeContentType(c.Writer, []string{"text/html; charset=utf-8"})
	c.Status(http.StatusOK)
	_, err := c.Writer.Write([]byte(v))
	return err
}

// String 将字符串写入响应
func (c *Context) String(v string) error {
	writeContentType(c.Writer, []string{"text/plain; charset=utf-8"})
	c.Status(http.StatusOK)
	_, err := c.Writer.Write([]byte(v))
	return err
}
