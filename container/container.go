package container

// Container ...
type Container interface {
	InitApp(filename string) error
}