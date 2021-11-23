package dd_cmd

// HandlerFunc 定义函数类型
type HandlerFunc func(*Context)

// HandlerFuncList 定义 HandlerFunc 函数类型切片
type HandlerFuncList []HandlerFunc

// Add 操作指令接口
// type Add func(cmd Executable) Executable

// ICommandHandler 定义所有路由 Handler 接口。
type ICommandHandler interface {
	RegCommandByChar(string, string, ...HandlerFunc) ICommandHandler
	RegCommand(Executable, string, ...HandlerFunc) ICommandHandler
	Handle(Executable, string, ...HandlerFunc) ICommandHandler

}

type IPrefixHandler interface {
	GetCommandPrefixs() []string
	GetPrefix(string) string
}