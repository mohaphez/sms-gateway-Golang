package utils

import (
	"errors"
	"math"
	"regexp"
)

// Calcualte the message count
func SmsCount(message string) (int16, string, error) {
	var length float64 = 1
	var language string
	characterCount := len([]rune(message))
	langRegex, _ := regexp.Compile(`^[A-Za-z0-9\-\?\!\@\(\)\$\^\%\*\&\#\{\}\[\]\:\.\<\>\s]+$`)
	lang := langRegex.MatchString(message)
	if lang {
		length = 160
		language = "en"
	} else {
		length = 70
		language = "fa"
	}
	if characterCount > int(length) {
		if lang {
			length = length - 7
		} else {
			length = length - 3
		}
	}
	count := math.Ceil(float64(characterCount) / length)
	if lang {
		if count > 8 {
			return int16(count), language, errors.New("message count is big")
		}
	} else {
		if count > 10 {
			return int16(count), language, errors.New("message count is big")
		}
	}
	return int16(count), language, nil
}

func StatusText(id int16) string {
	switch id {
	case 1:
		return "در صف ارسال"
	case 2:
		return "زمان بندی شده"
	case 3:
		return "ارسال شده به مخابرات"
	case 4:
		return "خطا در ارسال پیام (fail)"
	case 5:
		return "نرسیده به گیرنده (Undelivered)"
	case 6:
		return "بلاک شده است"
	case 7:
		return "ارسال ناموفق"
	case 8:
		return "پیامک در مخابرات است"
	case 9:
		return "پیامک توسط سرور دریافت نشده است (ریزش)"
	case 10:
		return "رسیده به گیرنده"
	case 11:
		return "ارسال ناموفق (کد خطا اپراتور را بررسی کنید)"
	case 100:
		return "شناسه پیامک نامعتبر است"
	default:
		return "ارسال شده به مخابرات"
	}
}
