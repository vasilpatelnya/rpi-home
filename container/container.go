package container

// Container контейнер приложения
type Container interface {
	InitApp(filename string) error
	Run()
}
