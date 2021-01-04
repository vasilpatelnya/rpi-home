package notification

type Notifier interface {
	SendText(t string) error
	SendFile(fp string, m string) error
}
