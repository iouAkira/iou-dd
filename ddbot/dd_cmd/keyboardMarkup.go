package dd_cmd

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type InlineKeyboardMarkup tgbotapi.InlineKeyboardMarkup

// Markupable 按钮包装接口
type Markupable interface {
	IkType() string                           //当前指令的返回
	MakeKeyboardMarkup() InlineKeyboardMarkup //对返回数据的按钮做包装
}

//RdcMarkup 读取cookies类
type RdcMarkup struct {
	Cmd      string
	FilePath string
	RowBtns  int
	Suffix   string
}

func (markup *RdcMarkup) IkType() string {
	return markup.Cmd
}

// MakeKeyboardMarkup todo 业务层待实现
func (markup *RdcMarkup) MakeKeyboardMarkup() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	return keyboardMarkup
}

//WrapCancelBtn 添加取消按钮
func WrapCancelBtn(markup *InlineKeyboardMarkup) InlineKeyboardMarkup {
	var cancelBTN = Command{Cmd: "cancel"}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("取消", cancelBTN.Run()),
	))
	return *markup

}

//WrapCancelWithExampleBtn 添加取消按钮并带有范例显示
func WrapCancelWithExampleBtn(markup *InlineKeyboardMarkup) InlineKeyboardMarkup {
	var cancelBTN = Command{Cmd: "cancel"}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("取消", cancelBTN.Run()),
		tgbotapi.NewInlineKeyboardButtonURL("查看配置示例🔗", "https://github.com/iouAkira/someDockerfile/blob/master/dd_scripts/genCodeConf.list"),
	))
	return *markup
}

//SpnodeMarkup 读取cookies类
type SpnodeMarkup struct {
	Cmd      string
	FilePath string
	RowBtns  int
	Suffix   string
}

func (markup *SpnodeMarkup) IkType() string {
	return markup.Cmd
}

// MakeKeyboardMarkup todo 业务层待实现
func (markup *SpnodeMarkup) MakeKeyboardMarkup() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	return keyboardMarkup
}

//LogsMarkup 读取cookies类
type LogsMarkup struct {
	Cmd      string
	FilePath string
	RowBtns  int
	Suffix   string
}

func (markup *LogsMarkup) IkType() string {
	return markup.Cmd
}

// MakeKeyboardMarkup todo 业务层待实现
func (markup *LogsMarkup) MakeKeyboardMarkup() InlineKeyboardMarkup {
	var keyboardMarkup InlineKeyboardMarkup
	return keyboardMarkup
}

func MakeKeyboard(markup Markupable) InlineKeyboardMarkup {
	return markup.MakeKeyboardMarkup()
}

func (markup InlineKeyboardMarkup) WithCancel() InlineKeyboardMarkup {
	return WrapCancelBtn(&markup)
}
func (markup InlineKeyboardMarkup) WithExampleBtn() InlineKeyboardMarkup {
	return WrapCancelWithExampleBtn(&markup)
}
