package i18n

import (
	"github.com/edy4c7/works-uploader/internal/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type Printer interface {
	Print(lang string, key string, params ...interface{}) string
}

type PrinterImpl struct {
	catalog catalog.Catalog
}

func NewPrinter() *PrinterImpl {
	builder := catalog.NewBuilder(catalog.Fallback(language.English))

	// English
	builder.SetString(language.English, errors.WUE00, "An error has occurred")
	builder.SetString(language.English, errors.WUE01, "Works is not found.")
	builder.SetString(language.English, errors.WUE02, "An error has occurred")
	builder.SetString(language.English, errors.WUE99, "A system error has occurred")

	// Japanese
	builder.SetString(language.Japanese, errors.WUE00, "%vの形式が不正です。")
	builder.SetString(language.Japanese, errors.WUE01, "指定された作品は見つかりません")
	builder.SetString(language.Japanese, errors.WUE02, "操作を行う権限がありません。")
	builder.SetString(language.Japanese, errors.WUE99, "システムエラーが発生しました。お手数ですが管理者にお問い合わせ下さい。 ")

	return &PrinterImpl{
		catalog: builder,
	}
}

func (r *PrinterImpl) Print(lang string, key string, params ...interface{}) string {
	tag, _ := language.MatchStrings(r.catalog.Matcher(), lang)
	printer := message.NewPrinter(tag, message.Catalog(r.catalog))
	return printer.Sprintf(key, params...)
}
