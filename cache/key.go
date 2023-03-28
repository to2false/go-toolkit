package cache

import "fmt"

type (
	Key string
)

var cacheNamespace = ""

func SetCacheNamespace(namespace string) {
	cacheNamespace = namespace
}

func (c Key) String() string {
	if cacheNamespace != "" {
		return cacheNamespace + ":" + string(c)
	}

	return string(c)
}

func (c Key) Format(v ...any) Key {
	return Key(fmt.Sprintf(string(c), v...))
}
