package dd_cmd

import "strings"

type RouterGroup struct {
	Handlers []HandlerFunc
	basePath string
	engine   *Engine
	root     bool
}

// 编译检查 RouterGroup 是否实现 IRouter 接口
var _ IRouter = &RouterGroup{}

// Use adds middleware to the group, see example code in GitHub.
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

func (group *RouterGroup) Handle(httpMethod Executable, relativePath string, handlers ...HandlerFunc) IRoutes {
	//if matched := regEnLetter.MatchString(httpMethod); !matched {
	//	panic("http method " + httpMethod + " is not valid")
	//}
	return group.handle(httpMethod, relativePath, handlers)
}

//
func (group *RouterGroup) handle(cmdMethod Executable, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(cmdMethod, absolutePath, handlers)
	return group.returnObj()
}

// calculateAbsolutePath 设计指令的绝对路径,并把多余的占位符删除
func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	relativePath = strings.Trim(relativePath, " ")
	return relativePath
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	// assert1(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

// Cmd 将会对所有路由消息中带有 CommandPrefix 前缀信息进行处理并添加到路由组中
func (group *RouterGroup) RegCommand(commandPrefix string, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(&Command{prefix: commandPrefix}, relativePath, handlers)
}

// Route 对消息命令的扩展，只要实现 Executable 接口即可对当前指令进行处理
// 例如: 如果你想实现一个 “>hit 2” , “>”即为 当前指令的prefix, 后期可以支持多字段的prefix 并且emoji亦可
func (group *RouterGroup) Route(cmd Executable, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(cmd, relativePath, handlers)
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
