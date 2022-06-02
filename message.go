package tgbot

import (
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
)

type MessageBuilder struct {
	editMessage         bool
	removeReplyKeyboard bool
	chatId              int64
	replyMessageId      int
	messageId           int
	inlineId            string
	text                string
	photoId             string
	keyboard            [][]tgbotapi.InlineKeyboardButton
	replyKeyboard       [][]tgbotapi.KeyboardButton
}

func (b *MessageBuilder) EditMessageTextAndMarkup(chatId int64, messageId int) *MessageBuilder {
	b.chatId = chatId
	b.messageId = messageId
	b.editMessage = true
	return b
}

func (b *MessageBuilder) NewMessage(chatId int64) *MessageBuilder {
	b.chatId = chatId
	b.editMessage = false
	return b
}

func (b *MessageBuilder) Message(chatId int64, messageId int) *MessageBuilder {
	if messageId == 0 {
		return b.NewMessage(chatId)
	} else {
		return b.EditMessageTextAndMarkup(chatId, messageId)
	}
}

func (b *MessageBuilder) Text(text string) *MessageBuilder {
	b.text = text
	return b
}

func (b *MessageBuilder) ChatId(chatId int64) *MessageBuilder {
	b.chatId = chatId
	return b
}

func (b *MessageBuilder) ReplyMessageId(replyMessageId int) *MessageBuilder {
	b.replyMessageId = replyMessageId
	return b
}

func (b *MessageBuilder) RemoveReplyKeyboard() *MessageBuilder {
	b.removeReplyKeyboard = true
	return b
}

func (b *MessageBuilder) PhotoId(fileId string) *MessageBuilder {
	b.photoId = fileId
	return b
}

func (b *MessageBuilder) MessageId(messageId int) *MessageBuilder {
	b.messageId = messageId
	return b
}

func (b *MessageBuilder) InlineId(inlineId string) *MessageBuilder {
	b.inlineId = inlineId
	return b
}

func (b *MessageBuilder) Edit(editMessage bool) *MessageBuilder {
	b.editMessage = editMessage
	return b
}

func (b *MessageBuilder) AddKeyboardRow() *MessageBuilder {
	b.keyboard = append(b.keyboard, []tgbotapi.InlineKeyboardButton{})
	return b
}

func (b *MessageBuilder) AddReplyKeyboardRow() *MessageBuilder {
	b.replyKeyboard = append(b.replyKeyboard, []tgbotapi.KeyboardButton{})
	return b
}

func (b *MessageBuilder) AddButton(text, callbackData string) *MessageBuilder {
	b.keyboard[len(b.keyboard)-1] = append(b.keyboard[len(b.keyboard)-1],
		tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callbackData})
	return b
}

func (b *MessageBuilder) AddReplyButton(text string) *MessageBuilder {
	b.replyKeyboard[len(b.replyKeyboard)-1] = append(b.replyKeyboard[len(b.replyKeyboard)-1], tgbotapi.KeyboardButton{Text: text})
	return b
}

func (b *MessageBuilder) AddReplyWebAppButton(text, url string) *MessageBuilder {
	button := tgbotapi.KeyboardButton{Text: text, WebApp: &tgbotapi.WebAppInfo{URL: url}}
	b.replyKeyboard[len(b.replyKeyboard)-1] = append(b.replyKeyboard[len(b.replyKeyboard)-1], button)
	return b
}
func (b *MessageBuilder) AddReplyRequestContactButton(text string) *MessageBuilder {
	button := tgbotapi.KeyboardButton{Text: text, RequestContact: true}
	b.replyKeyboard[len(b.replyKeyboard)-1] = append(b.replyKeyboard[len(b.replyKeyboard)-1], button)
	return b
}

func (b *MessageBuilder) AddWebAppButton(text, url string) *MessageBuilder {
	button := tgbotapi.InlineKeyboardButton{Text: text, WebApp: &tgbotapi.WebAppInfo{URL: url}}
	b.keyboard[len(b.keyboard)-1] = append(b.keyboard[len(b.keyboard)-1], button)
	return b
}

func (b *MessageBuilder) AddButtonUrl(text, url string) *MessageBuilder {
	b.keyboard[len(b.keyboard)-1] = append(b.keyboard[len(b.keyboard)-1],
		tgbotapi.InlineKeyboardButton{Text: text, URL: &url})
	return b
}

func (b *MessageBuilder) AddButtonSwitch(text, sw string) *MessageBuilder {
	b.keyboard[len(b.keyboard)-1] = append(b.keyboard[len(b.keyboard)-1],
		tgbotapi.NewInlineKeyboardButtonSwitch(text, sw),
	)
	return b
}

func (b *MessageBuilder) Build() tgbotapi.Chattable {
	if b.editMessage {
		kb := b.getKeyboard()
		var msg tgbotapi.Chattable
		if len(kb) > 0 {
			m := tgbotapi.NewEditMessageTextAndMarkup(b.chatId, b.messageId, b.text, tgbotapi.NewInlineKeyboardMarkup(kb...))
			m.ParseMode = tgbotapi.ModeMarkdown
			m.InlineMessageID = b.inlineId
			msg = m
		} else {
			m := tgbotapi.NewEditMessageText(b.chatId, b.messageId, b.text)
			m.ParseMode = tgbotapi.ModeMarkdown
			m.InlineMessageID = b.inlineId
			msg = m
		}
		return msg

	} else if b.photoId != "" {
		msg := tgbotapi.NewPhoto(b.chatId, tgbotapi.FileID(b.photoId))
		msg.Caption = b.text
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg

	} else {
		msg := tgbotapi.NewMessage(b.chatId, b.text)
		keyboard := b.getKeyboard()
		replyKeyboard := b.getReplyKeyboard()
		msg.ReplyToMessageID = b.replyMessageId

		if len(keyboard) > 0 {
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
		} else if len(replyKeyboard) > 0 {
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(replyKeyboard...)
		} else if b.removeReplyKeyboard {
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		}
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}
}

func (b *MessageBuilder) getKeyboard() [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, buttons := range b.keyboard {
		if len(buttons) > 0 {
			keyboard = append(keyboard, buttons)
		}
	}
	return keyboard
}
func (b *MessageBuilder) getReplyKeyboard() [][]tgbotapi.KeyboardButton {
	var keyboard [][]tgbotapi.KeyboardButton

	for _, buttons := range b.replyKeyboard {
		if len(buttons) > 0 {
			keyboard = append(keyboard, buttons)
		}
	}
	return keyboard
}

type inlineMessageBuilder struct {
	inlineQueryId string
	articles      []*tgbotapi.InlineQueryResultArticle
}

func NewInlineRequest(inlineQueryId string) *inlineMessageBuilder {
	return &inlineMessageBuilder{inlineQueryId: inlineQueryId}
}

func (b *inlineMessageBuilder) AddArticle(id, title, descr, text string) *inlineMessageBuilder {
	article := tgbotapi.NewInlineQueryResultArticleMarkdown(id, title, text)
	article.Description = descr
	b.articles = append(b.articles, &article)
	return b
}

func (b *inlineMessageBuilder) getLastArticleMarkup() *tgbotapi.InlineKeyboardMarkup {
	article := b.articles[len(b.articles)-1]
	if article.ReplyMarkup != nil {
		return article.ReplyMarkup
	} else {
		markup := tgbotapi.NewInlineKeyboardMarkup()
		article.ReplyMarkup = &markup
		return article.ReplyMarkup
	}
}

func (b *inlineMessageBuilder) AddKeyboardRow() *inlineMessageBuilder {
	markup := b.getLastArticleMarkup()
	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{})
	return b
}

func (b *inlineMessageBuilder) AddButton(text, callbackData string) *inlineMessageBuilder {
	markup := b.getLastArticleMarkup()

	markup.InlineKeyboard[len(markup.InlineKeyboard)-1] = append(markup.InlineKeyboard[len(markup.InlineKeyboard)-1],
		tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callbackData})
	return b
}

func (b *inlineMessageBuilder) AddButtonSwitch(text, sw string) *inlineMessageBuilder {
	markup := b.getLastArticleMarkup()

	markup.InlineKeyboard[len(markup.InlineKeyboard)-1] = append(markup.InlineKeyboard[len(markup.InlineKeyboard)-1],
		tgbotapi.NewInlineKeyboardButtonSwitch(text, sw),
	)
	return b
}

func (b *inlineMessageBuilder) AddButtonUrl(text, url string) *inlineMessageBuilder {
	markup := b.getLastArticleMarkup()

	markup.InlineKeyboard[len(markup.InlineKeyboard)-1] = append(markup.InlineKeyboard[len(markup.InlineKeyboard)-1],
		tgbotapi.InlineKeyboardButton{Text: text, URL: &url},
	)
	return b
}

func (b *inlineMessageBuilder) Build() tgbotapi.Chattable {

	var articles []interface{}
	for _, article := range b.articles {
		articles = append(articles, *article)
	}

	return tgbotapi.InlineConfig{
		InlineQueryID: b.inlineQueryId,
		IsPersonal:    true,
		Results:       articles,
	}

}
