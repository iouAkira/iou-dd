package dd_cmd

import (
	"fmt"
	"log"
	"strings"
	"sync"

	ddutils "ddbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Engine 中注入了Bot所需要的一些，相当与一个大框架，使用时候需要New()来对Engine初始化操作
type Engine struct {
	CommandHandler
	handlerPrfixList HandlerPrefixList
	pool             sync.Pool
	Token            string
	Userid           int64
	bot              *tgbotapi.BotAPI
}

// New 返回一个 Engine 实体,初始化操作并不包含任何路由和中间件
func New() *Engine {
	engine := &Engine{
		CommandHandler: CommandHandler{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
	}
	engine.CommandHandler.engine = engine
	return engine
}

// addCommand 将所有handlerPrfix相同(指令前缀)的指令合并到一颗树组合中。
func (engine *Engine) addCommand(handlerPrfix Executable, command string, handlers HandlerFuncList) {
	//查询 handlerPrfix 节点是否存在  例如 "/" "@" ">"，
	//对为空对象是初始化添加节点组
	commandTree := engine.handlerPrfixList.get(handlerPrfix)
	if commandTree == nil {
		tree := new(HandlerPrefix)
		tree.handlerPrfix = handlerPrfix
		tree.commands = &CommandNodes{}
		engine.handlerPrfixList = append(engine.handlerPrfixList, tree)
		commandTree = engine.handlerPrfixList.get(handlerPrfix)
	}
	//todo 解决获取节点到底是哪里的问题
	commandTree.commands.addCommandNode(command, handlers)
}

// GetCommandPrefixs 获取当前engine的指令前缀集合
func (engine *Engine) GetCommandPrefixs() []string {
	var prefixs []string
	for _, v := range engine.handlerPrfixList {
		prefixs = append(prefixs, v.handlerPrfix.Prefix())
	}
	return prefixs
}

// Run 主入口函数包含两个两个参数，token tg的bot token,userid 用户的id
// 函数在启动时候会直接将 tgbotapi 进行初始化操作,并在结束进程时候停止信息的接受
func (engine *Engine) Run(token string, userid int64) {
	engine.Token = token
	engine.Userid = userid
	bot, err := tgbotapi.NewBotAPI(token)
	bot.Debug = false
	if err != nil {
		log.Panicf("start bot failed with some error %v", err)
	}
	log.Printf("Telegram bot stared，Bot info ==> %s %s[%s]", bot.Self.FirstName, bot.Self.LastName, bot.Self.UserName)
	c := &Context{Request: bot}
	c.reset()
	c.Request = bot
	c.engine = engine
	engine.pool.Put(c)
	engine.bot = bot
	engine.handleHTTPRequest(c)
}

// handleHTTPRequest
func (engine *Engine) handleHTTPRequest(c *Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.Request.GetUpdatesChan(u)
	for update := range updates {
		//log.Println(update.Message)
		//log.Println(update.CallbackQuery)
		if update.Message != nil && update.Message.Document != nil {
			fmt.Printf("%+v\n", "sss")
			continue
		}
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		c.mu.Lock()
		c.Update = &update
		c.mu.Unlock()
		engine.handleRequest(c)
	}
}

// handleRequest 对请求数据处理进不同的路由中,由解析器将用户返回的消息解析加工到 Executable 中,
func (engine *Engine) handleRequest(c *Context) {
	var msg string
	if c.Update.Message != nil {
		log.Printf("来源Message: %s", c.Update.Message.Text)
		msg = c.Update.Message.Text
	} else {
		log.Printf("来源CallbackQuery: %s", c.Update.CallbackQuery.Data)
		msg = c.Update.CallbackQuery.Data
	}

	msg = strings.Trim(msg, " ")
	//对路由集合遍历的查询开头与请求一致的指令
	var hp *HandlerPrefix
	var msgPrefix string
	for _, tree := range engine.handlerPrfixList {
		log.Printf("匹配指令前缀: %v", tree.handlerPrfix.Prefix())
		if strings.HasPrefix(msg, tree.handlerPrfix.Prefix()) {
			hp = tree
			msgPrefix = tree.handlerPrfix.Prefix()
		}
	}

	if hp != nil && len(*hp.commands) > 0 {
		hasCommand := false
		log.Printf("入口命令：%s", msgPrefix)
		for _, command := range *hp.commands {
			log.Printf("before cleanPath: %v", msg)
			log.Printf("cleanPath: %v", ddutils.CleanCommand(msg, 0)[0])
			if fmt.Sprintf("%s%s", hp.handlerPrfix.Prefix(), command.commandStr) == ddutils.CleanCommand(msg, 0)[0] {
				hasCommand = true
				log.Printf("存在命令 %s", command.commandStr)
				log.Printf("等待命令 %s", msg)
				log.Printf("接受命令 %s", msg)
				c.HandlerPrefixStr = ddutils.CleanCommand(msg, 0)[0]
				for _, handler := range command.handlers {
					handler(c)
				}
			}
		}
		if !hasCommand {
			log.Printf("不存在命令 %s，响应[unkonw]", msg)
			unknowMsg := "/unknow"
			hp = engine.handlerPrfixList.get(&Command{prefix: "/"})
			for _, route := range *hp.commands {
				if fmt.Sprintf("%s%s", hp.handlerPrfix.Prefix(), route.commandStr) == ddutils.CleanCommand(unknowMsg, 0)[0] {
					for _, handler := range route.handlers {
						handler(c)
					}
				}
			}
		}

	}
}
