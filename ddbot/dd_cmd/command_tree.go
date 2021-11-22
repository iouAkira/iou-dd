package dd_cmd

// commandTrees 是指令组群，通常是需要挂载到 Engine 上来实现多个指令的启用
type commandTrees []*commandTree

// commandTree 是单一指令的集合，包含了当前指令的命令名和 当前 root 上存在的指令数组
type commandTree struct {
	method Executable    //method 是指的命令前缀,如 /cmd /help中的 "/"指令接口
	root   *commandNodes //root   是当前指令下的各种命令的集合
}

//get 获取当前前缀下的指令
func (trees commandTrees) get(method Executable) *commandTree {
	if len(trees) == 0 {
		return nil
	}
	for _, tree := range trees {
		if tree.method.Prefix() == method.Prefix() {
			return tree
		}
	}
	return nil
}

// addRoute 添加指令并封装加入函数
func (nodes *commandNodes) addNode(path string, handlers HandlersChain) {
	cmdPath := path
	var hasPath bool
	//log.Println(len(*nodes))
	for _, s := range *nodes {
		hasPath = false
		if s.path == cmdPath {
			hasPath = true
			break
		}
	}
	if !hasPath {
		*nodes = append(*nodes, &commandNode{path: cmdPath, handlers: handlers})
	}
}
