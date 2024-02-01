package mygin

import (
	"net/http"
	"sync"
)

// HandlerFunc 定义处理函数类型
type HandlerFunc func(*Context)

// HandlersChain 定义处理函数链类型
type HandlersChain []HandlerFunc

// Engine 定义引擎结构，包含路由器
type Engine struct {
	Router
	RouterGroup
	pool sync.Pool
}

func (e *Engine) Use(middleware ...HandlerFunc) IRoutes {
	e.RouterGroup.Use(middleware...)
	return e
}

// ServeHTTP 实现http.Handler接口的方法，用于处理HTTP请求
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取对应HTTP方法的路由树的根节点
	root := e.trees.get(r.Method)
	// 解析请求路径
	parts := root.parseFullPath(r.URL.Path)

	// 查找符合条件的节点
	searchNode := make([]*node, 0)
	root.search(parts, &searchNode)

	// 没有匹配到路由
	if len(searchNode) == 0 {
		w.Write([]byte("404 Not found!\n"))
		return
	}

	/**
	 *searchNode 可能返回多个，这里只取返回第一个
	 *如果一个url多个路由都能匹配成功，那么就该改写router了
	 */
	// 参数赋值
	params := make([]Param, 0)
	searchPath := root.parseFullPath(searchNode[0].fullPath)
	for i, sp := range searchPath {
		if sp[0] == ':' {
			params = append(params, Param{
				Key:   sp[1:],
				Value: parts[i],
			})
		}
	}

	// 获取处理函数链
	handlers := searchNode[0].handlers
	if handlers == nil {
		w.Write([]byte("404 Not found!\n"))
		return
	}

	//从pool中取
	c := e.pool.Get().(*Context)
	c.Request = r
	c.Writer = w
	c.Params = params
	c.handlers = handlers
	c.index = -1

	// 执行处理函数链
	c.Next()

	//归还到pool中
	e.pool.Put(c)

}

// Default 返回一个默认的引擎实例
func Default() *Engine {
	engine := New()

	engine.Use(Logger(), Recovery())

	// Group 保存 engine 的指针
	engine.RouterGroup.engine = engine

	return engine
}

// New 返回一个引擎实例
func New() *Engine {
	engine := &Engine{
		Router: Router{
			trees: make(methodTrees, 0, 9),
		},
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
	}

	engine.pool.New = func() any {
		return &Context{
			Request:  nil,
			Writer:   nil,
			Params:   make(Params, 0),
			handlers: nil,
			index:    -1,
		}
	}
	// Group 保存 engine 的指针
	engine.RouterGroup.engine = engine

	return engine
}

// Run 启动HTTP服务器的方法
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
