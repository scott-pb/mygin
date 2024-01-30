package mygin

import (
	"strings"
)

type nodeType uint8

// 路由的类型
const (
	static nodeType = iota
	root
	param
	catchAll
)

// 不同的method 对应不同的节点树 定义
type methodTree struct {
	method string
	root   *node
}

// Param 参数的类型key=> value
type Param struct {
	Key   string
	Value string
}

// Params 切片
type Params []Param

type methodTrees []methodTree

type node struct {
	children  []*node
	part      string
	wildChild bool
	handlers  HandlersChain
	nType     nodeType
	fullPath  string
}

// Get 获取 参数中的值
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

// ByName 通过ByName获取参数中的值 会忽略掉错误，默认返回 空字符串
func (ps Params) ByName(name string) (v string) {
	v, _ = ps.Get(name)
	return
}

// 根据method获取root
func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

// 添加路径时
func (n *node) addRoute(path string, handlers HandlersChain) {

	//根据请求路径按照'/'划分
	parts := n.parseFullPath(path)

	//将节点插入路由后，返回最后一个节点
	matchNode := n.insert(parts)

	//最后的节点，绑定执行链
	matchNode.handlers = handlers

	//最后的节点，绑定完全的URL，后续param时有用
	matchNode.fullPath = path

}

// 按照 "/" 拆分字符串
func (n *node) parseFullPath(fullPath string) []string {
	splits := strings.Split(fullPath, "/")
	parts := make([]string, 0)
	for _, part := range splits {
		if part != "" {
			parts = append(parts, part)
			if part == "*" {
				break
			}
		}
	}
	return parts
}

// 根据路径 生成节点树
func (n *node) insert(parts []string) *node {
	part := parts[0]
	//默认的字节类型为静态类型
	nt := static
	//根据前缀判断节点类型
	switch part[0] {
	case ':':
		nt = param
	case '*':
		nt = catchAll
	}

	//插入的节点查找
	var matchNode *node
	for _, childNode := range n.children {
		if childNode.part == part {
			matchNode = childNode
		}
	}

	//如果即将插入的节点没有找到，则新建一个
	if matchNode == nil {
		matchNode = &node{
			part:      part,
			wildChild: part[0] == '*' || part[0] == ':',
			nType:     nt,
		}
		//新子节点追加到当前的子节点中
		n.children = append(n.children, matchNode)
	}

	//当最后插入的节点时，类型赋值，且返回最后的节点
	if len(parts) == 1 {
		matchNode.nType = nt
		return matchNode
	}

	//匹配下一部分
	parts = parts[1:]
	//子节点继续插入剩余字部分
	return matchNode.insert(parts)
}

// 根据路由 查询符合条件的节点
func (n *node) search(parts []string, searchNode *[]*node) {
	part := parts[0] //a

	allChild := n.matchChild(part) //b c :name

	if len(parts) == 1 {
		// 如果到达路径末尾，将所有匹配的节点加入结果
		*searchNode = append(*searchNode, allChild...)
		return
	}

	parts = parts[1:] //b

	for _, n2 := range allChild {
		// 递归查找下一部分
		n2.search(parts, searchNode)
	}

}

// 根据part 返回匹配成功的子节点
func (n *node) matchChild(part string) []*node {

	allChild := make([]*node, 0)
	for _, child := range n.children {
		if child.wildChild || child.part == part {
			allChild = append(allChild, child)
		}
	}

	return allChild
}
