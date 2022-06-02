package tgbot

import (
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"log"
	"strings"
)

type Update struct {
	tgbotapi.Update
	chatProv ChatProvider

	chat *ChatInfo
	btn  *Button
	usr  User
}

func WrapUpdate(update tgbotapi.Update, user User, chatProvider ChatProvider) *Update {
	return &Update{Update: update, usr: user, chatProv: chatProvider}
}

func (u *Update) GetUserId() int64 {
	if u.Message != nil && u.Message.Chat != nil {
		return u.Message.From.ID
	}
	if u.EditedMessage != nil && u.EditedMessage.Chat != nil {
		return u.EditedMessage.From.ID
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.From.ID
	}
	if u.InlineQuery != nil {
		return u.InlineQuery.From.ID
	}
	if u.MyChatMember != nil {
		return u.MyChatMember.From.ID
	}
	return 0
}

func (u *Update) GetChatId() int64 {
	if u.Message != nil && u.Message.Chat != nil {
		return u.Message.Chat.ID
	}
	if u.CallbackQuery != nil && u.CallbackQuery.Message != nil {
		return u.CallbackQuery.Message.Chat.ID
	}
	return 0
}

func (u *Update) GetUser() User {
	return u.usr
}

func (u *Update) GetMessageId() int {
	if u.IsButton() && u.CallbackQuery != nil && u.CallbackQuery.Message != nil {
		return u.CallbackQuery.Message.MessageID
	} else if u.Message != nil {
		return u.Message.MessageID
	}
	return 0
}

func (u *Update) GetInlineMessageId() string {
	if u.InlineQuery != nil {
		return u.InlineQuery.ID
	}
	return ""
}

func (u *Update) HasText(text string) bool {
	return u.Update.Message != nil && text == u.Update.Message.Text
}

func (u *Update) StartsWithText(prefix string) bool {
	return u.Update.Message != nil && strings.HasPrefix(u.Update.Message.Text, prefix) ||
		u.InlineQuery != nil && strings.HasPrefix(u.InlineQuery.Query, prefix)
}

func (u *Update) HasCommand(text string) bool {
	return u.IsCommand() && text == u.Update.Message.Text
}

func (u *Update) StartsWithCommand(prefix string) bool {
	return u.IsCommand() && strings.HasPrefix(u.Update.Message.Text, prefix)
}

func (u *Update) IsCommand() bool {
	return u.Update.Message != nil &&
		strings.HasPrefix(u.Update.Message.Text, "/")
}

func (u *Update) IsPrivate() bool {
	return u.Message != nil && u.Message.Chat.Type == "private" ||
		u.CallbackQuery != nil && u.CallbackQuery.Message != nil && u.CallbackQuery.Message.Chat.Type == "private" ||
		u.InlineQuery != nil && u.InlineQuery.ChatType == "sender"
}

//Button

func (u *Update) IsPlainText() bool {
	return !u.IsCommand() && u.Update.Message != nil && u.Update.Message.Text != ""
}

func (u *Update) GetText() string {
	if u.Message == nil {
		return ""
	}
	text := strings.ReplaceAll(u.Message.Text, "_", "\\_")
	text = strings.ReplaceAll(text, "*", "\\*")
	text = strings.ReplaceAll(text, "~", "\\~")
	return text
}

func (u *Update) GetInline() string {
	if u.InlineQuery != nil {
		return u.InlineQuery.Query
	}
	return ""
}

func (u *Update) GetInlineId() string {
	if u.CallbackQuery != nil {
		return u.CallbackQuery.InlineMessageID
	}
	if u.InlineQuery != nil {
		return u.InlineQuery.ID
	}
	return ""
}

func (u *Update) IsButton() bool {
	return u.Update.CallbackData() != ""
}

func (u *Update) GetButton() Button {
	if u.btn == nil {
		button, err := u.chatProv.GetButton(u.CallbackData())
		if err != nil {
			log.Printf("[ERROR] cannot find button %s, %v", u.CallbackData(), err)
		}
		u.btn = &button
	}
	return *u.btn
}

func (u *Update) GetButtonById(btnId string) Button {
	button, err := u.chatProv.GetButton(btnId)
	if err != nil {
		log.Printf("[ERROR] cannot find button %s, %v", u.CallbackData(), err)
	}

	return button
}

func (u *Update) HasAction(action Action) bool {
	return u.IsButton() && u.GetButton().HasAction(action)
}

// ChatInfo

func (u *Update) HasActionOrChain(actionOrChain Action) bool {
	return u.IsButton() && u.GetButton().HasAction(actionOrChain) ||
		u.GetChatInfo().ActiveChain == string(actionOrChain)
}

func (u *Update) HasChain(chain Action) bool {
	return u.GetChatInfo().ActiveChain == string(chain)
}

func (u *Update) GetChatInfo() *ChatInfo {
	if u.chat == nil {
		chat, err := u.chatProv.GetChatInfo(u.GetUserId())
		if err != nil {
			log.Printf("[WARN] cannot find chat info, %v", err)
		}
		u.chat = &chat
	}

	if u.chat.ChatId == 0 {
		u.chat.ChatId = u.GetUserId()
	}
	if u.chat.ChainData == nil {
		u.chat.ChainData = Data{}
	}

	return u.chat
}

func (u *Update) FlushChatInfo() {
	err := u.chatProv.SaveChatInfo(*u.GetChatInfo())
	if err != nil {
		log.Printf("[ERROR] cannot save chat info: %+v", u.GetChatInfo())
	}
}

func (u *Update) StartChain(chain string) *Update {
	u.GetChatInfo().ActiveChain = chain
	return u
}

func (u *Update) StartChainStep(chainStep string) *Update {
	u.GetChatInfo().ActiveChainStep = chainStep
	return u
}

func (u *Update) GetChainStep() string {
	return u.GetChatInfo().ActiveChainStep
}

func (u *Update) GetChain() string {
	return u.GetChatInfo().ActiveChain
}

func (u *Update) AddChainData(key string, value string) *Update {
	u.GetChatInfo().ChainData[key] = value
	return u
}

func (u *Update) GetChainData(key string) string {
	return u.GetChatInfo().ChainData[key]
}

func (u *Update) FinishChain() *Update {
	u.GetChatInfo().ActiveChain = ""
	u.GetChatInfo().ActiveChainStep = ""
	u.GetChatInfo().ChainData = map[string]string{}
	return u
}

func (u *Update) SaveChatInfo(info ChatInfo) error {
	return u.chatProv.SaveChatInfo(info)
}
