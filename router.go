package mygin

type Router struct {
	trees methodTrees
}

// 添加路由方法
func (r *Router) addRoute(method, path string, handlers HandlersChain) {
	//根据method获取root
	rootTree := r.trees.get(method)

	//如果root为空
	if rootTree == nil {
		//初始化一个root
		rootTree = &node{part: "/", nType: root}
		//将初始化后的root 加入tree树中
		r.trees = append(r.trees, methodTree{method: method, root: rootTree})

	}

	rootTree.addRoute(path, handlers)

}
