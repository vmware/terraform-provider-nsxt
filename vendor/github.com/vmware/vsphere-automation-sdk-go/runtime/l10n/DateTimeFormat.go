/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package l10n

// Constants that map to the enum values in std.LocalizationParam.DateTimeFormat
const (
	//The date and time value will be formatted as short date for example
	//<em>2019-01-28</em>
	DateTimeFormat_SHORT_DATE = "SHORT_DATE"
	//The date and time value will be formatted as medium date for example
	//<em>2019 Jan 28</em>
	DateTimeFormat_MED_DATE = "MED_DATE"
	//The date and time value will be formatted as long date for example
	//<em>2019 Jan 28</em>
	DateTimeFormat_LONG_DATE = "LONG_DATE"
	//The date and time value will be formatted as full date for example
	//<em>2019 Jan 28, Mon</em>
	DateTimeFormat_FULL_DATE = "FULL_DATE"
	//The date and time value will be formatted as short time for example
	//<em>12:59</em>
	DateTimeFormat_SHORT_TIME = "SHORT_TIME"
	//The date and time value will be formatted as medium time for example
	//<em>12:59:01</em>
	DateTimeFormat_MED_TIME = "MED_TIME"
	//The date and time value will be formatted as long time for example
	//<em>12:59:01 Z</em>
	DateTimeFormat_LONG_TIME = "LONG_TIME"
	//The date and time value will be formatted as full time for example
	//<em>12:59:01 Z</em>
	DateTimeFormat_FULL_TIME = "FULL_TIME"
	//The date and time value will be formatted as short date and time
	//for example <em>2019-01-28 12:59</em>
	DateTimeFormat_SHORT_DATE_TIME = "SHORT_DATE_TIME"
	//The date and time value will be formatted as medium date and time
	//for example <em>2019 Jan 28 12:59:01</em>
	DateTimeFormat_MED_DATE_TIME = "MED_DATE_TIME"
	//The date and time value will be formatted as long date and time for
	//example <em>2019 Jan 28 12:59:01 Z</em>
	DateTimeFormat_LONG_DATE_TIME = "LONG_DATE_TIME"
	//The date and time value will be formatted as full date and time
	//<em>2019 Jan 28, Mon 12:59:01 Z</em>
	DateTimeFormat_FULL_DATE_TIME = "FULL_DATE_TIME"
)
// TODO [l18n-support]: No built in l18n support in go. Use external library
// to get date format based on format locale
// https://stackoverflow.com/questions/21400121/localization-when-using-time-format
// https://github.com/nicksnyder/go-i18n/issues/45
// Apply format locale using https://github.com/go-playground/universal-translator

// TODO[l18n-support]: REMOVE this. It is map of temporary date and time formatting
// options. It combines the Unicode CLDR format types -
// full, long, medium and short with 3 different presentations - date only,
// time only and combined date and time presentation.
// Each layout string (used for formatting) is a representation of the time stamp,
//Mon Jan 2 15:04:05 -0700 MST 2006
var DateTimeFormats map[string]string
func initDateTimeFormats() {
	if len(DateTimeFormats) == 0 {
	DateTimeFormats = map[string]string {
		// The date and time value will be formatted as short date for example
		// <em>2019-01-28</em>
		DateTimeFormat_SHORT_DATE: "2006-01-02",

		// The date and time value will be formatted as medium date for example
		// <em>2019 Jan 28</em>
		DateTimeFormat_MED_DATE: "2006 Jan 02",

		// The date and time value will be formatted as long date for example
		// <em>2019 Jan 28</em>
		DateTimeFormat_LONG_DATE: "2006 Jan 02",

		// The date and time value will be formatted as full date for example
		// <em>2019 Jan 28, Mon</em>
		DateTimeFormat_FULL_DATE: "2006 Jan 02, Mon",

		// The date and time value will be formatted as short time for example
		// <em>12:59</em>
		DateTimeFormat_SHORT_TIME: "15:04",

		// The date and time value will be formatted as medium time for example
		// <em>12:59:01</em>
		DateTimeFormat_MED_TIME: "15:04:05",

		// The date and time value will be formatted as long time for example
		// <em>12:59:01 Z</em>
		DateTimeFormat_LONG_TIME: "15:04:05 MST",

		// The date and time value will be formatted as full time for example
		// <em>12:59:01 Z</em>
		DateTimeFormat_FULL_TIME: "15:04:05 MST",

		// The date and time value will be formatted as short date and time
		// for example <em>2019-01-28 12:59</em>
		DateTimeFormat_SHORT_DATE_TIME: "2006-01-02 15:04",

		// The date and time value will be formatted as medium date and time
		// for example <em>2019 Jan 28 12:59:01</em>
		DateTimeFormat_MED_DATE_TIME: "2006 Jan 02 15:04:05",

		// The date and time value will be formatted as long date and time for
		// example <em>2019 Jan 28 12:59:01 Z</em>
		DateTimeFormat_LONG_DATE_TIME: "2006 Jan 02 15:04:05 MST",

		// The date and time value will be formatted as full date and time
		// <em>2019 Jan 28, Mon 12:59:01 Z</em>
		DateTimeFormat_FULL_DATE_TIME: "2006 Jan 02, Mon 15:04:05 MST",
	}
}
}