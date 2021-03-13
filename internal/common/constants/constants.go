package constants

type ContentType int

const (
	ContentTypeFile = iota
	ContentTypeURL
)

type ActivityType int

const (
	ActivityAdded ActivityType = iota
	ActivityUpdated
)
