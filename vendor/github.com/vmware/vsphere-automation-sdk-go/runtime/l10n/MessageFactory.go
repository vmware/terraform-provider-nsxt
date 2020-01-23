/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package l10n

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// Interface to create formatter from a message catalog
type LocalizableMessageFactory interface {
	GetDefaultFormatter() *MessageFormatter
	GetFormatterForLocalizationParams(acceptLanguage *string,
		formatLanguage *string,
		timezone *string) (MessageFormatter, error)
}

// Utility struct for maintaining message catalog and creating message formatter
type MessageFactory struct {
	messageCatalog   *catalog.Builder
	defaultFormatter *MessageFormatter
	DateTimeFormat   *string
	Precision        *int64
}

const INVALID_MESSAGE_ID = "messageformatter.invalid.message.id"

func NewMessageFactory() MessageFactory {
	return NewMessageFactoryWithSettings(nil, nil)
}
func NewMessageFactoryWithSettings(dateTimeFormat *string, precision *int64) MessageFactory {
	// Fallback is the language used when performing MatchString
	cat := catalog.NewBuilder(catalog.Fallback(language.English))
	mf := MessageFactory{cat, nil, dateTimeFormat, precision}
	df := NewDefaultMessageFormatter(mf)
	mf.defaultFormatter = &df
	return mf
}

func NewMessageFactoryWithBundle(localesAndFiles map[string]string,
	dateTimeFormat *string, precision *int64) (MessageFactory, error) {
	mf := NewMessageFactoryWithSettings(dateTimeFormat, precision)
	mf.AddBundles(localesAndFiles)
	if !mf.supportsEnglish() {
		var emptyMf MessageFactory
		errorMsg := "Invalid catalog state. No english messages available. " +
			"'en' catalog should be added because it is the default language"
		log.Error(errorMsg)
		return emptyMf, errors.New(errorMsg)
	}
	return mf, nil
}

func (mf MessageFactory) supportsEnglish() bool {
	log.Info("Supported languages are: ", mf.messageCatalog.Languages())
	for _, tag := range mf.messageCatalog.Languages() {
		if tag == language.English {
			return true
		}
	}
	return false
}

// Load key value pairs of message IDs and message strings from properties file
// for a provided locale
func (mf MessageFactory) AddBundle(localeStr string, messages io.Reader) error {

	locale, err := language.Parse(localeStr)
	if err != nil {
		log.Error(err)
		return err
	}
	scanner := bufio.NewScanner(messages)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		var value string
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}
		if len(value) == 0 {
			return fmt.Errorf("No value for key: %s in catalog file", key)
		}
		if len(key) == 0 {
			return fmt.Errorf("No key for value: %s in catalog file", value)
		}
		setStringErr := mf.messageCatalog.SetString(locale, key, value)
		if setStringErr != nil {
			return setStringErr
		}
	}
	setStringErr := mf.messageCatalog.SetString(locale, INVALID_MESSAGE_ID, "Unknown message ID %s requested")
	if setStringErr != nil {
		return setStringErr
	}
	if err := scanner.Err(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Load key value pairs of message IDs and message strings from map of locale
// and its corresponding properties file
func (mf MessageFactory) AddBundles(localesAndFiles map[string]string) error {
	for localeStr, messageFileName := range localesAndFiles {
		file, err := os.Open(messageFileName)
		if err != nil {
			log.Error(err)
			return err
		}
		defer file.Close()
		if err := mf.AddBundle(localeStr, file); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// Locale string to map of message id to message template
type LocaleToMessages func(string) map[string]string

// Available locales is the list of locales supplied by the function in second parameter
// Second parameter is a function which returns a map of message id to message template for the locale provided
func (mf MessageFactory) AddMessages(availableLocales []string, localeToMessages LocaleToMessages) error {
	for _, localeStr := range availableLocales {
		locale, err := language.Parse(localeStr)
		if err != nil {
			log.Error(err)
			return err
		}
		for key, value := range localeToMessages(localeStr) {
			mf.messageCatalog.SetString(locale, key, value)
		}
	}
	return nil
}

func (m MessageFactory) Catalog() catalog.Catalog {
	return m.messageCatalog
}

// Construct a message formatter from the localization headers from application
// context in the execution context of a request
func (m MessageFactory) GetFormatterForLocalizationParams(acceptLanguage *string,
	formatLanguage *string,
	timezone *string) (MessageFormatter, error) {
	if !m.supportsEnglish() {
		var emptyMf MessageFormatter
		errorMsg := "Invalid catalog state. No english messages available. " +
			"'en' catalog should be added because it is the default language"
		log.Error(errorMsg)
		return emptyMf, errors.New(errorMsg)
	}
	var bundleLocale *language.Tag
	var formatLocale *language.Tag
	var timeZoneLocation *time.Location

	var emptyMf MessageFormatter

	if acceptLanguage != nil && *acceptLanguage != "" {
		loc, _ := language.MatchStrings(m.messageCatalog.Matcher(), *acceptLanguage)
		bundleLocale = &loc
	}

	if formatLanguage != nil && *formatLanguage != "" {
		// TODO [l18n-support]: Create matcher from available l18n locales from external library
		matcher := language.NewMatcher(m.messageCatalog.Languages())
		loc, _ := language.MatchStrings(matcher, *formatLanguage)
		formatLocale = &loc
	}

	if timezone != nil && *timezone != "" {
		tz, err := time.LoadLocation(*timezone)
		if err != nil {
			log.Errorf("Error loading time zone from %s", *timezone)
			return emptyMf, err
		}
		timeZoneLocation = tz
	}

	return NewMessageFormatter(bundleLocale, formatLocale, timeZoneLocation, m), nil
}

func (m MessageFactory) GetDefaultFormatter() *MessageFormatter {
	return m.defaultFormatter
}
