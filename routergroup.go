package mygin

import (
	"net/http"
	"path"
	"regexp"
)

// IRoutes 定义了路由组的接口
type IRoutes interface {
	BasePath() string
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes
	Match([]string, string, ...HandlerFunc) IRoutes
}

// anyMethods 包含所有 HTTP 方法的字符串表示
var anyMethods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

// RouterGroup 定义了路由组的结构体
type RouterGroup struct {
	Handlers HandlersChain // 路由组的中间件处理函数链
	basePath string        // 路由组的基础路径
	engine   *Engine       // 路由组所属的引擎
	root     bool          // 是否是根路由组
}

// Use 添加中间件在路由组中
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	assertFunc(len(group.Handlers)+len(middleware) < int(abortIndex), "too many middlewares")
	group.Handlers = append(group.Handlers, middleware...)
	if group.root {
		return group.engine
	}
	return group
}

// Group 创建一个新的路由组
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	assertFunc(len(group.Handlers)+len(handlers) < int(abortIndex), "too many handlers")
	return &RouterGroup{
		Handlers: append(group.Handlers, handlers...),
		basePath: path.Join(group.basePath, relativePath),
		engine:   group.engine,
	}
}

// BasePath 返回路由组的基础路径
func (group *RouterGroup) BasePath() string {
	return group.basePath
}

// handle 处理路由，将路由信息添加到引擎中
func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := path.Join(group.basePath, relativePath)
	assertFunc(len(group.Handlers)+len(handlers) < int(abortIndex), "too many handlers")
	handlers = append(group.Handlers, handlers...)
	group.engine.addRoute(httpMethod, absolutePath, handlers)

	if group.root {
		return group.engine
	}
	return group
}

// Handle 校验 HTTP 方法的有效性，并处理路由
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	// 检查 HTTP 方法的有效性
	if match := regexp.MustCompile("^[A-Z]+$").MatchString(httpMethod); !match {
		panic("http method " + httpMethod + " is not valid")
	}
	// 处理路由
	return group.handle(httpMethod, relativePath, handlers)
}

// GET 注册 GET 方法的路由
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers)
}

// POST 注册 POST 方法的路由
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPost, relativePath, handlers)
}

// DELETE 注册 DELETE 方法的路由
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

// PATCH 注册 PATCH 方法的路由
func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPatch, relativePath, handlers)
}

// PUT 注册 PUT 方法的路由
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPut, relativePath, handlers)
}

// OPTIONS 注册 OPTIONS 方法的路由
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodOptions, relativePath, handlers)
}

// HEAD 注册 HEAD 方法的路由
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodHead, relativePath, handlers)
}

// Match 注册多个方法的路由
func (group *RouterGroup) Match(methods []string, relativePath string, handlers ...HandlerFunc) IRoutes {
	for _, method := range methods {
		group.handle(method, relativePath, handlers)
	}

	if group.root {
		return group.engine
	}
	return group
}

// Any 注册所有方法的路由
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	for _, method := range anyMethods {
		group.handle(method, relativePath, handlers)
	}

	if group.root {
		return group.engine
	}
	return group
}

func assertFunc(guard bool, text string) {
	if !guard {
		panic(text)
	}
}
