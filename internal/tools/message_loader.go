package tools

type MessageLoader interface {
	LoadMessage(key string, locale string, params ...interface{}) string
}
