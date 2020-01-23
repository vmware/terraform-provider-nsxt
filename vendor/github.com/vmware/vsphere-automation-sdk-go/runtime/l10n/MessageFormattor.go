/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package l10n

import (
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
	"time"
)

// Localizes values and strings with particular localization settings including
// message strings locale, format locale and timezone
type MessageFormatter struct {
	// Language of message text - Default English
	acceptLocale language.Tag
	// Printer configured with a catalog and accept language of this formatter
	messagePrinter *message.Printer
	// To decide date format to use - Defaults to accept locale
	formatLocale language.Tag
	// For formatting datetime parameter - Default UTC
	timezone time.Location
	// For formatting floatPrecision in float and double - Default "2"
	floatPrecision int64
	// For formatting floatPrecision in float and double - Default "%.2f"
	floatFormat string
	// Format of string representation of datetime parameter
	dateTimeFormat string
}

// Construct a message formatter with default localization settings
func NewDefaultMessageFormatter(msgFactory MessageFactory) MessageFormatter {
	return NewMessageFormatter(nil, nil, nil, msgFactory)
}

func NewMessageFormatter(acceptLocale *language.Tag,
	formatLocale *language.Tag,
	timezone *time.Location,
	msgFactory MessageFactory) MessageFormatter {
	initDateTimeFormats()
	catalog := msgFactory.Catalog()
	if acceptLocale == nil || !contains(*acceptLocale, catalog.Languages()) {
		acceptLocale = &language.English
	}
	if formatLocale == nil {
		// English or accept locale?
		//formatLocale = &language.English
		formatLocale = acceptLocale
	}
	if timezone == nil {
		timezone = time.UTC
	}
	precision := msgFactory.Precision
	if precision == nil {
		p := int64(2)
		precision = &p
	}
	dateTimeFormat := msgFactory.DateTimeFormat
	if dateTimeFormat == nil {
		s := DateTimeFormat_SHORT_DATE_TIME
		dateTimeFormat = &s
	}
	floatFormat := fmt.Sprintf("%%.%df", *precision)
	messagePrinter := message.NewPrinter(*acceptLocale, message.Catalog(catalog))
	log.Infof("Message formatter created for %v, %v, %v, %v, %v",
		*acceptLocale, *formatLocale, timezone.String(), *precision, *dateTimeFormat)
	return MessageFormatter{*acceptLocale, messagePrinter,
		*formatLocale, *timezone, *precision,
		floatFormat, *dateTimeFormat}
}

// Loads and formats a message with positional arguments. It is expected
// that the caller has formatted any date time and number objects to
// strings using the other methods of this class
func (f MessageFormatter) GetPositionalArgMessage(msgID string, args ...interface{}) string {
	template := f.messagePrinter.Sprintf(f.messageKey(msgID), args...)
	return template
}

// Loads and formats a message with named arguments. It is expected that the
// caller has formatted any date time and number objects to strings using
// the other methods of this class
func (f MessageFormatter) GetLocalizedMessage(msgId string, args map[string]string) string {
	message := f.messagePrinter.Sprintf(f.messageKey(msgId))
	for key, val := range args {
		if !strings.Contains(message, key) {
			log.Errorf("Error constructing message %s: Invalid parameter %s", msgId, key)
		}
		message = strings.Replace(message, "{"+key+"}", val, -1)
	}
	if strings.Contains(message, "{") {
		log.Errorf("Error constructing message %s: Insufficient parameters provided for message %s", msgId, message)
	}
	return message
}

func (f MessageFormatter) FormatInterfaceParam(param interface{}) string {
	str := fmt.Sprintf("%v", param)
	if p, ok := param.(float32); ok {
		str = f.FormatFloat(float64(p), nil)
	} else if p, ok := param.(float64); ok {
		str = f.FormatFloat(p, nil)
	} else if p, ok := param.(time.Time); ok {
		str = f.FormatDateTime(p, nil)
	}
	return str
}

func (f MessageFormatter) FormatFloat(param float64, precision *int64) string {
	var floatFormat string
	if precision == nil {
		floatFormat = f.floatFormat
	} else {
		floatFormat = fmt.Sprintf("%%.%df", *precision)
	}
	return fmt.Sprintf(floatFormat, param)
}

func (mf MessageFormatter) FormatDateTime(param time.Time, format *string) string {
	// TODO [l18n-support]: No built in l18n support in go. Use external library
	// to get date format based on format locale
	// https://stackoverflow.com/questions/21400121/localization-when-using-time-format
	// https://github.com/nicksnyder/go-i18n/issues/45
	// Apply format locale using https://github.com/go-playground/universal-translator
	var dtFormat string
	if format == nil {
		dtFormat = DateTimeFormats[mf.dateTimeFormat]
	} else {
		dtFormat = DateTimeFormats[*format]
	}
	return param.In(&mf.timezone).Format(dtFormat)
}

func (mf MessageFormatter) DateTimeFormat() string {
	return mf.dateTimeFormat
}

func (mf MessageFormatter) FloatPrecision() int64 {
	return mf.floatPrecision
}

func (f MessageFormatter) messageKey(msgID string) message.Reference {
	invalidId := f.messagePrinter.Sprintf(INVALID_MESSAGE_ID, msgID)
	return message.Key(msgID, invalidId)
}

func contains(a language.Tag, list []language.Tag) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
