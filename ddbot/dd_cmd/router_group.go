package dd_cmd

// HandlerFunc 定义函数类型
type HandlerFunc func(*Context)

// HandlersChain 定义 HandlerFunc 函数类型切片
type HandlersChain []HandlerFunc

// Add 操作指令接口
// type Add func(cmd Executable) Executable
