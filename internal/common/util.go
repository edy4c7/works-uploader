package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func MapToStruct(src interface{}, dest interface{}) error {
	tmp, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(tmp, dest)
	if err != nil {
		return err
	}

	return nil
}

// offsetとlimitがクエリパラメータで指定されている場合、その値を抽出する。
// 1番目の戻り値がoffset、2番目の戻り値がlimitを表す。
// 指定されていない場合は、offset=0、limit=-1 を返す。
func ExtractOffsetAndLimit(r *http.Request) (int, int, error) {
	offset, limit := 0, -1

	strOffset := r.URL.Query().Get("offset")
	if strOffset != "" {
		o, err := strconv.Atoi(strOffset)
		if err != nil {
			return 0, 0, err
		}
		offset = o
	}

	strLimit := r.URL.Query().Get("limit")
	if strLimit != "" {
		l, err := strconv.Atoi(strLimit)
		if err != nil {
			return 0, 0, err
		}
		limit = l
	}

	return offset, limit, nil
}

func GetScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
