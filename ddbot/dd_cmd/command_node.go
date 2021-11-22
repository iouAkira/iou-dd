package dd_cmd

type commandNodes []*commandNode

// commandNode 为指令节点，每个指令节点都对应
type commandNode struct {
	path string // 所在指令的路径  /[cmd]
	// children commandNodes  // 子节点挂载位
	handlers HandlersChain // 回调处理函数
}
