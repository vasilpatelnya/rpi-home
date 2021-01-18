package main

import (
	"log"
	"os"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	/*
		todo гуй будет иметь конфигурацию запуска в бд и возможность запускать приложение с этими настройками
		также должна быть возможность запуска в терминальном режиме
	*/
	a := app.New()
	w := a.NewWindow("Raspberry Pi Home")

	refreshIntervalEntry := widget.NewEntry()

	w.SetContent(&widget.Form{
		Items: []*widget.FormItem{widget.NewFormItem("refresh interval", refreshIntervalEntry)},
		OnSubmit: func() {
			log.Println(refreshIntervalEntry.Text)
		},
		OnCancel: func() {
			log.Println("close app...")
			os.Exit(0)
		},
		SubmitText: "Запуск",
		CancelText: "Отменить",
	})

	w.ShowAndRun()
}
