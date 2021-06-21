package translate

import "sync"

var trans Dictionary = nil
var once sync.Once

type Translator interface {
	Text(name int) string
}

func T() Translator {
	once.Do(func() {
		trans = d
	})

	return trans
}
