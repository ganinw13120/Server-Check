package service

import (
	"fmt"
	"regexp"
	"server-health/model"
	"server-health/repository"
	"strings"
	"sync"
	"time"

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
			err := h.removeWishList(v)
			if err != nil {
				return err
			}
		} else if match, _ := regexp.MatchString("http", v.Message.Text); match {
			err := h.addWishList(v)
			if err != nil {
				return err
			}
		} else if match, _ := regexp.MatchString("all", v.Message.Text); match {
			err := h.showWishList(v.Source.UserID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h healthService) removeWishList(hook model.WebhookEvent) error {
	path := strings.Replace(hook.Message.Text, "delete ", "", 1)
	path = strings.Replace(path, " ", "", 0)
	err := h.wishListRepository.RemoveWishList(hook.Source.UserID, path)
	if err != nil {
		return nil
	}
	replyMessage := linebot.NewTextMessage("ลบสำเร็จ")
	_, err = h.linebot.ReplyMessage(hook.ReplyToken, replyMessage).Do()
	return err
}

func (h healthService) addWishList(hook model.WebhookEvent) error {
	path := strings.Replace(hook.Message.Text, " ", "", 0)
	err := h.wishListRepository.AddWishList(hook.Source.UserID, path)
	if err != nil {
		return nil
	}
	replyMessage := linebot.NewTextMessage("เพิ่มสำเร็จ")
	_, err = h.linebot.ReplyMessage(hook.ReplyToken, replyMessage).Do()
	return err
}

func (h healthService) showWishList(line_id string) error {
	lists, err := h.wishListRepository.GetPersonWishList(line_id)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	healthList := make([]model.Health, 0)
	for _, v := range lists {
		wg.Add(1)
		go func(v model.WishList) {
			health := h.healthRepository.CheckHealth(v.Path)
			healthList = append(healthList, *health)
			defer wg.Done()
		}(v)
		wg.Wait()
	}
	err = h.generateFlexMessage(line_id, healthList)
	if err != nil {
		return err
	}
	return nil
}

func (h healthService) generateFlexMessage(line_id string, lists []model.Health) error {
	list := ""
	for k, v := range lists {
		comma := ","
		if k == 0 {
			comma = ""
		}
		var isPassed string
		var statusColor string
		if v.IsAlive {
			isPassed = "Passed"
			statusColor = "#21BF65"
		} else {
			isPassed = "Failed"
			statusColor = "#fc5f51"
		}
		var responseTimeColor string
		if v.ResponseTime < 300*time.Millisecond {
			responseTimeColor = "#21BF65"
		} else if v.ResponseTime < 300*time.Millisecond {
			responseTimeColor = "#d9b638"
		} else {
			responseTimeColor = "#fc5f51"
		}
		list = fmt.Sprintf(`%s 
		%s{
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "%s",
			  "color": "#404040"
			},
			{
			  "type": "text",
			  "text": "hello, world",
			  "contents": [
				{
				  "type": "span",
				  "text": "Status : ",
				  "size": "sm"
				},
				{
				  "type": "span",
				  "text": "%s",
				  "color": "%s",
				  "weight": "bold"
				}
			  ]
			},
			{
			  "type": "text",
			  "text": "hello, world",
			  "contents": [
				{
				  "type": "span",
				  "text": "Response Time : ",
				  "size": "sm"
				},
				{
				  "type": "span",
				  "text": "%v",
				  "color": "%s",
				  "weight": "bold"
				}
			  ]
			},
			{
				"type": "button",
				"action": {
					"type": "message",
					"label": "Delete",
					"text": "delete %s"
				},
				"margin": "lg",
				"height": "sm",
				"color": "#eb6a50",
				"style": "primary",
				"gravity": "center",
				"adjustMode": "shrink-to-fit"
			},
			{
			  "type": "separator",
			  "margin": "10px"
			}
		  ],
		  "paddingTop": "10px",
		  "paddingBottom": "10px"
		}`, list, comma, v.Path, isPassed, statusColor, v.ResponseTime, responseTimeColor, v.Path)
	}
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(fmt.Sprintf(`{
		"type": "bubble",
		"size": "mega",
		"header": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "box",
			  "layout": "vertical",
			  "contents": [
				{
				  "type": "text",
				  "text": "Check",
				  "color": "#FAFAFA",
				  "size": "sm"
				},
				{
				  "type": "text",
				  "text": "Server Status",
				  "color": "#FAFAFA",
				  "size": "xl",
				  "flex": 4,
				  "weight": "bold"
				}
			  ]
			}
		  ],
		  "paddingAll": "20px",
		  "backgroundColor": "#21BF65",
		  "spacing": "md",
		  "height": "100px",
		  "paddingTop": "22px"
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			  %s
		  ],
		  "backgroundColor": "#F2F2F2"
		}
	  }`, list)))
	if err != nil {
		return err
	}
	flexMessage := linebot.NewFlexMessage("Status Result", flexContainer)
	_, err = h.linebot.PushMessage(line_id, flexMessage).Do()
	return err
}
