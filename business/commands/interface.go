package commands

type Service interface {
	ProcessCommand(commandSpec CommandSpec) error
	AsyncStockCommand(stockCode string, chatroom string) error
}