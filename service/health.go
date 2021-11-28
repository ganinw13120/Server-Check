package service

import (
	"regexp"
	"server-health/model"
	"server-health/repository"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type healthService struct {
	healthRepository   repository.IHealthRepository
	wishListRepository repository.IWishListRepository
	linebot            *linebot.Client
}

type IHealthService interface {
	WebhookEnter(hook model.LineWebhook) error
}

func NewHealthService(
	healthRepository repository.IHealthRepository,
	wishListRepository repository.IWishListRepository,
	linebot *linebot.Client,
) healthService {
	return healthService{
		healthRepository:   healthRepository,
		wishListRepository: wishListRepository,
		linebot:            linebot,
	}
}

func (h healthService) WebhookEnter(hook model.LineWebhook) error {
	for _, v := range hook.Events {
		if match, _ := regexp.MatchString("delete ", v.Message.Text); match {
			err := h.removeWishList(v.Source.UserID, v.Message.Text)
			if err != nil {
				return err
			}
		} else if match, _ := regexp.MatchString("http", v.Message.Text); match {
			err := h.addWishList(v.Source.UserID, v.Message.Text)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h healthService) removeWishList(line_id string, message string) error {
	path := strings.Replace(message, "delete ", "", 1)
	path = strings.Replace(path, " ", "", 0)
	err := h.wishListRepository.RemoveWishList(line_id, path)
	return err
}

func (h healthService) addWishList(line_id string, message string) error {
	path := strings.Replace(message, " ", "", 0)
	err := h.wishListRepository.AddWishList(line_id, path)
	return err
}

func (h healthService) generateFlexMessage() error {
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(`PASTE_JSON_TO_HERE`))
	if err != nil {
		return err
	}
	flexMessage := linebot.NewFlexMessage("FlexWithJSON", flexContainer)
	_, err = h.linebot.ReplyMessage("", flexMessage).Do()
	return err
}
