package commands

import (
	"chat-hex/business"
	"chat-hex/business/messages"
	"encoding/csv"
	"errors"
	"net/http"
	"strings"
	"time"
)

type CommandSpec struct {
	Content  string `validate:"required"`
	Sender   string `validate:"required"`
	Chatroom string `validate:"required"`
}

type service struct{
	messagesService	messages.Service
}

func NewService(messagesService	messages.Service) Service {
	return &service{
		messagesService: messagesService,
	}
}

func (s *service) ProcessCommand(commandSpec CommandSpec) error {
	pieces := strings.Split(commandSpec.Content, "=")

	possibleCommand := strings.ToLower(pieces[0])
	if possibleCommand == CommandStock {
		stockCode := strings.ToLower(pieces[1])
		if len(stockCode) <= 0 {
			return business.ErrInvalidCommand
		}

		err := s.AsyncStockCommand(stockCode, commandSpec.Chatroom)
		if err != nil {
			return business.ErrInvalidCommand
		}

		return nil
	}

	return business.ErrInvalidCommand
}

func (s *service) AsyncStockCommand(stockCode string, chatroom string) error {
	httpRequest, err := http.NewRequest("GET", "https://stooq.com/q/l/?s="+stockCode+"&f=sd2t2ohlcv&h&e=csv", nil) //getURL is presignedURL which returns csv file.
	if err != nil {
		return err
	}

	client := http.Client{Timeout: time.Second * 10}
	response, err := client.Do(httpRequest)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		content, err := csv.NewReader(response.Body).ReadAll()
		if err != nil {
			return err
		}

		headers := content[0]
		var symbolValue string
		var closeValue string
		for index, header := range headers {
			if strings.ToLower(header) == "symbol" {
				symbolValue = content[1][index]
			}
			if strings.ToLower(header) == "close" {
				closeValue = content[1][index]
			}
		}

		if symbolValue == "" || closeValue == "" || strings.ToLower(closeValue) == "n/d" {
			return errors.New("unavailable values")
		}

		stockMessageSpec  := messages.InsertMessageSpec{
			Content: symbolValue + " quote is $" + closeValue + " per share",
			Sender: "Stockbot",
			Chatroom: chatroom,
		}

		 err = s.messagesService.InsertMessage(stockMessageSpec)
		 if err != nil {
			return err
		 }
	}

	return nil
}