package i18n

type MessageLoader interface {
	LoadMessage(key string, locale string, params ...interface{}) string
}
