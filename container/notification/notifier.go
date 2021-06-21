package notification

const TypeTelegram = "telegram"

type Notifier interface {
	SendText(t string) error
	SendFile(fp string, m string) error
}
