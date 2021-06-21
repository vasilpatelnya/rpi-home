package translate

import "os"

var Months = map[string]string{
	"January": "января", "February": "февраля", "March": "марта",
	"April": "апреля", "May": "мая", "June": "июня",
	"July": "июля", "August": "августа", "September": "сентября",
	"October": "октября", "November": "ноября", "December": "декабря",
}

const (
	ErrorCreateContainer int = iota
	ErrorParsingEnv
	ErrorLoggerInit
	ErrorRootPath
	ErrorConfigLoad
	ErrorCreateConnectionContainer
	ErrorNotifierInit
	ErrorRepoInit
)

const (
	Ua = "ua"
	Ru = "ru"
	En = "en"
)

// Dictionary represents a dictionary entry of the form: [constant for a message, for example 1] ["language - string
// constant"] "text"
type Dictionary map[int]map[string]string

var d = Dictionary{
	ErrorCreateContainer:           map[string]string{En: "Error on try create application container", Ru: "Ошибка при попытке создания контейнера приложения"},
	ErrorParsingEnv:                map[string]string{En: "Parse environment mode error", Ru: "Ошибка парсинга уровня работы приложения"},
	ErrorLoggerInit:                map[string]string{En: "Logger initialization error", Ru: "Ошибка при инициализации логгера"},
	ErrorRootPath:                  map[string]string{En: "Root path not founded", Ru: "Корневой путь не найден"},
	ErrorConfigLoad:                map[string]string{En: "Error loading configuration file", Ru: "Ошибка при загрузке конфигурационного файла"},
	ErrorCreateConnectionContainer: map[string]string{En: "Create connection container error", Ru: "Ошибка создания контейнера подключений к БД"},
	ErrorNotifierInit:              map[string]string{En: "Error while initializing the module for sending notifications", Ru: "Ошибка при инициализации модуля отправки уведомлений"},
	ErrorRepoInit:                  map[string]string{En: "Repository initialization error", Ru: "Ошибка инициализации репозитория"},
}

func (d Dictionary) Text(name int) string {
	loc := os.Getenv("RPIHOME_LANG")
	if loc == "" {
		loc = En
	}

	return d[name][loc]
}
