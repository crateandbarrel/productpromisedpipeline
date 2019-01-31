package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//TranslationResult Use to capture google translate result
type TranslationResult struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

//TranslateString Use to translate from one locale to another
func TranslateString(str string, fromLocale string, toLocale string) (string, error) {
	var returnError error
	var returnString string

	source := GetCultureLocale(fromLocale)
	target := GetCultureLocale(toLocale)

	source = (strings.Split(source, "-"))[0]
	target = (strings.Split(target, "-"))[0]

	host := "www.googleapis.com"
	path := "/language/translate/v2"
	apiKey := "AIzaSyBzcKt7siOprMiMUeoIcrsZ_6sadpQpN54"

	path = path +
		"?key=" + apiKey +
		"&q=" + url.QueryEscape(str) +
		"&source=" + source +
		"&target=" + target

	url := "https://" + host + path
	resp, err := http.Get(url)
	if err != nil {
		returnError = err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var result = TranslationResult{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		returnError = errors.New("Couldn't marshall result: " + err.Error())
	}

	if resp.StatusCode < 300 {
		returnString = result.Data.Translations[0].TranslatedText
	} else {
		returnError = errors.New("the source (" + source + ") or target (" + target + ") language code is invalid")
		returnString = str
	}

	return returnString, returnError
}

//ApplyLocaletoCurrency Use to convert numbers into locale-specific currencies
func ApplyLocaletoCurrency(value float64, locale string) string {
	var returnValue string

	returnValue = padDecimal(value)

	switch {
	case locale == "en-us":
		returnValue = "$" + returnValue
	case locale == "fr-ca":
		returnValue = strings.Replace(returnValue, ".", ",", -1)
		returnValue = "CAD " + returnValue
	case locale == "en-ca":
		returnValue = "CAD " + returnValue
	}
	return returnValue
}

func padDecimal(number float64) string {
	numberString := fmt.Sprintf("%.2f", number)
	if strings.Index(numberString, ".") == -1 {
		numberString = numberString + "."
	}
	for len(numberString) < strings.Index(numberString, ".")+3 {
		numberString = numberString + "0"
	}
	return numberString
}

//GetCultureLocale Used to get the IEEE locale from our company-based locales
func GetCultureLocale(locale string) string {
	//full cratebrowser locale looks like company-language-country - e.g. cb-en-us
	if len(locale) >= 8 {
		parts := strings.Split(locale, "-")
		locale = parts[1] + "-" + parts[2]
	}

	return locale
}
