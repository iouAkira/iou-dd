package dd_cmd

import "log"

type CommandHandler struct {
	Handlers []HandlerFunc
	basePath string
	engine   *Engine
	root     bool
}

// 编译检查 CommandHandler 是否实现 ICommandHandler 接口
var _ ICommandHandler = &CommandHandler{}

// Handle [ICommandHandler]接口定义的方法
func (commandHandler *CommandHandler) Handle(httpMethod Executable, command string, handlers ...HandlerFunc) ICommandHandler {
	return commandHandler.handle(httpMethod, command, handlers)
}

// RegCommandByChar [ICommandHandler]接口定义的方法；启动初始化注册监听指令
func (commandHandler *CommandHandler) RegCommandByChar(commandPrefix string, command string, handlers ...HandlerFunc) ICommandHandler {
	return commandHandler.handle(&Command{prefix: commandPrefix,Cmd: command}, command, handlers)
}

// RegCommand 对消息命令的扩展，只要实现 Executable 接口即可对当前指令进行处理
// 例如: 如果你想实现一个 “>hit 2” , “>”即为 当前指令的prefix, 后期可以支持多字段的prefix 并且emoji亦可
func (commandHandler *CommandHandler) RegCommand(cmd Executable, command string, handlers ...HandlerFunc) ICommandHandler {
	return commandHandler.handle(cmd, command, handlers)
}

// Use adds middleware to the commandHandler, see example code in GitHub.
func (commandHandler *CommandHandler) Use(middleware ...HandlerFunc) ICommandHandler {
	commandHandler.Handlers = append(commandHandler.Handlers, middleware...)
	return commandHandler.returnObj()
}

//
func (commandHandler *CommandHandler) handle(cmdMethod Executable, command string, handlers HandlerFuncList) ICommandHandler {
	handlers = commandHandler.combineHandlers(handlers)
	commandHandler.engine.addCommand(cmdMethod, command, handlers)
	log.Printf("[CommandHandler] 注册指令: %v %v，帮助：%v", cmdMethod.Prefix(), command, cmdMethod.Description())
	return commandHandler.returnObj()
}

// 注册中间件
func (commandHandler *CommandHandler) combineHandlers(handlers HandlerFuncList) HandlerFuncList {
	finalSize := len(commandHandler.Handlers) + len(handlers)
	mergedHandlers := make(HandlerFuncList, finalSize)
	copy(mergedHandlers, commandHandler.Handlers)
	copy(mergedHandlers[len(commandHandler.Handlers):], handlers)
	return mergedHandlers
}

func (commandHandler *CommandHandler) returnObj() ICommandHandler {
	if commandHandler.root {
		return commandHandler.engine
	}
	return commandHandler
}
