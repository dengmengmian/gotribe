package util

import (
	"gotribe/internal/core/constant"
	"time"
)

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(constant.TIME_FORMAT)
}
