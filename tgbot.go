package tgbot

import (
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type ChatProvider interface {
	GetChatInfo(chatId int64) (ChatInfo, error)
	SaveChatInfo(chat ChatInfo) error
	GetButton(btnId string) (Button, error)
	SaveButton(button Button) error
	SaveUser(user User) error
}

type Bot struct {
	*tgbotapi.BotAPI
	handler  func(update *Update)
	chatProv ChatProvider
	BotSelf  tgbotapi.User
}

func NewBot(token string, chatProv ChatProvider) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create TgBot")
	}
	me, _ := api.GetMe()
	return &Bot{BotAPI: api, chatProv: chatProv, BotSelf: me}, nil
}

func NewBotWithAPIEndpoint(token, apiEndpoint string, chatProv ChatProvider) (*Bot, error) {
	api, err := tgbotapi.NewBotAPIWithAPIEndpoint(token, apiEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create TgBot")
	}
	me, _ := api.GetMe()
	return &Bot{BotAPI: api, chatProv: chatProv, BotSelf: me}, nil
}

func (b *Bot) StartLongPolling(handler func(update *Update)) error {
	if b.handler != nil {
		return errors.New("long polling already started")
	}
	b.handler = handler
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range b.GetUpdatesChan(u) {
		wrappedUpdate, err := b.WrapUpdate(update)
		if err != nil {
			log.Printf("[ERROR]rapped update failed, %v", err)
			continue
		}
		b.handler(wrappedUpdate)
	}
	return nil
}

func (b *Bot) WrapUpdate(update tgbotapi.Update) (*Update, error) {
	user, err := b.SaveUser(&update)
	if err != nil {
		log.Printf("[ERROR] WrapUpdate, %v", err)
		return nil, err
	}
	return WrapUpdate(update, user, b.chatProv), nil
}

func (b *Bot) WrapRequest(req *http.Request) (*Update, error) {
	update, err := b.HandleUpdate(req)
	if err != nil {
		return nil, err
	}
	return b.WrapUpdate(*update)
}

func (b *Bot) SaveUser(update *tgbotapi.Update) (User, error) {
	tgUser, err := getFrom(update)
	if err != nil {
		return User{}, err
	}

	user := User{UserId: tgUser.ID, UserName: tgUser.UserName, DisplayName: tgUser.FirstName, LastName: tgUser.LastName}
	err = b.chatProv.SaveUser(user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getFrom(update *tgbotapi.Update) (tgbotapi.User, error) {
	var user *tgbotapi.User
	if update.CallbackQuery != nil {
		user = update.CallbackQuery.From
	} else if update.Message != nil {
		user = update.Message.From
	} else if update.EditedMessage != nil {
		user = update.EditedMessage.From
	} else if update.InlineQuery != nil {
		user = update.InlineQuery.From
	} else if update.MyChatMember != nil {
		user = &update.MyChatMember.From
	} else if update.PreCheckoutQuery != nil {
		user = update.PreCheckoutQuery.From
	} else {
		return tgbotapi.User{}, errors.Errorf("Not define user, update - %v", update)
	}
	return *user, nil
}
