package i18n

import (
	"testing"
)

func TestCurrencyFormat(t *testing.T) {

	price := 1234.00

	priceString := ApplyLocaletoCurrency(price, "en-us")
	if priceString != "$1234.00" {
		t.Errorf("Should be $1234.00 but got " + priceString)
	}

	priceString = ApplyLocaletoCurrency(price, "fr-ca")
	if priceString != "CAD 1234,00" {
		t.Errorf("Should be CAD 1234,00 but got " + priceString)
	}

	priceString = ApplyLocaletoCurrency(price, "en-ca")
	if priceString != "CAD 1234.00" {
		t.Errorf("Should be CAD 1234.00 but got " + priceString)
	}

}

func TestLanguageLocale(t *testing.T) {

	locale := "en-us"

	langLocale := GetCultureLocale(locale)
	if langLocale != "en-us" {
		t.Errorf("Should be en-us but got " + langLocale)
	}

	locale = "cb-en-us"

	langLocale = GetCultureLocale(locale)
	if langLocale != "en-us" {
		t.Errorf("Should be en-us but got " + langLocale)
	}

	locale = "cn-fr-ca"

	langLocale = GetCultureLocale(locale)
	if langLocale != "fr-ca" {
		t.Errorf("Should be fr-ca but got " + langLocale)
	}

}

func TestTranslation(t *testing.T) {

	locale := "cb-en-us"

	enLangLocale := GetCultureLocale(locale)

	locale = "cn-fr-ca"

	frLangLocale := GetCultureLocale(locale)

	if translation, err := TranslateString("rouge", frLangLocale, enLangLocale); err != nil || translation != "red" {
		t.Errorf("Should be \"red\" but got " + translation)
	}

}

func TestInvalidLanguageTranslation(t *testing.T) {

	locale := "cb-en-us"

	enLangLocale := GetCultureLocale(locale)

	locale = "cb-sg-sg"

	invalidLangLocale := GetCultureLocale(locale)

	if translation, err := TranslateString("red", invalidLangLocale, enLangLocale); err == nil || translation != "red" {
		t.Errorf("Should have gotten and error and \"red\" as a result.  But got no error and/or got " + translation)
	}

}
