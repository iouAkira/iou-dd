package dd_cmd

// HandlerPrefix 是单一指令的集合，包含了当前指令的命令名和 当前 root 上存在的指令数组
type HandlerPrefix struct {
	handlerPrfix Executable    //handlerPrfix 是指的命令前缀,如 /cmd /help中的 "/"指令接口
	commands     *CommandNodes //commands   是handlerPrfix指令下的各种命令的集合
}

// commandNode 为指令节点，每个指令节点都对应
type CommandNode struct {
	commandStr string          // 指令
	handlers   HandlerFuncList // 回调处理函数
}

// HandlerPrefixs 是指令组群，通常是需要挂载到 Engine 上来实现多个指令的启用
type HandlerPrefixList []*HandlerPrefix

//get 获取当前前缀下的指令
func (trees HandlerPrefixList) get(handlerPrfix Executable) *HandlerPrefix {
	if len(trees) == 0 {
		return nil
	}
	for _, tree := range trees {
		if tree.handlerPrfix.Prefix() == handlerPrfix.Prefix() {
			return tree
		}
	}
	return nil
}

type CommandNodes []*CommandNode

// addCommandNode 添加指令并封装加入函数
func (nodes *CommandNodes) addCommandNode(path string, handlers HandlerFuncList) {
	cmdPath := path
	var hasPath bool
	for _, s := range *nodes {
		hasPath = false
		if s.commandStr == cmdPath {
			hasPath = true
			break
		}
	}
	if !hasPath {
		*nodes = append(*nodes, &CommandNode{commandStr: cmdPath, handlers: handlers})
	}
}

func (receiver *HandlerPrefix) get(cmdStr string) *CommandNode  {
	if len(*receiver.commands) == 0 {
		return nil
	}
	for _, tree := range *receiver.commands {
		if tree.commandStr == cmdStr {
			return tree
		}
	}
	return nil
}