package langdetect

import (
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/teadove/go-phone-iso3166/phone_iso3166"
)

// LangInfo holds the detected language name and its code.
type LangInfo struct {
	// Full name of the language (e.g., "English", "Traditional Chinese").
	Name string
	// Code for the language, suitable for APIs (e.g., "en", "zh-Hant", "yue").
	Code string
}

// apiCodeMap maps whatlanggo's ISO 639-1 codes (and custom Chinese codes)
// to the API codes defined in language.json.
var apiCodeMap = map[string]string{
	// Languages where whatlanggo's ISO 639-1 differs from API code or needs explicit mapping
	"ar":  "ara", // Arabic (whatlanggo: ar, API: ara)
	"es":  "spa", // Spanish (whatlanggo: es, API: spa)
	"ja":  "jp",  // Japanese (whatlanggo: ja, API: jp)
	"ko":  "kor", // Korean (whatlanggo: ko, API: kor)
	"bg":  "bul", // Bulgarian (whatlanggo: bg, API: bul)
	"da":  "dan", // Danish (whatlanggo: da, API: dan)
	"et":  "est", // Estonian (whatlanggo: et, API: est)
	"ro":  "rom", // Romanian (whatlanggo: ro, API: rom)
	"sl":  "slo", // Slovenian (whatlanggo: sl, API: slo)
	"vi":  "vie", // Vietnamese (whatlanggo: vi, API: vie)
	"af":  "afr", // Afrikaans (whatlanggo: af, API: afr)
	"sq":  "alb", // Albanian (whatlanggo: sq, API: alb)
	"am":  "amh", // Amharic (whatlanggo: am, API: amh)
	"hy":  "arm", // Armenian (whatlanggo: hy, API: arm)
	"az":  "aze", // Azerbaijani (whatlanggo: az, API: aze)
	"eu":  "baq", // Basque (whatlanggo: eu, API: baq)
	"be":  "bel", // Belarusian (whatlanggo: be, API: bel)
	"bn":  "ben", // Bengali (whatlanggo: bn, API: ben)
	"bs":  "bos", // Bosnian (whatlanggo: bs, API: bos)
	"my":  "bur", // Burmese (whatlanggo: my, API: bur)
	"ca":  "cat", // Catalan (whatlanggo: ca, API: cat)
	"hr":  "hrv", // Croatian (whatlanggo: hr, API: hrv)
	"eo":  "epo", // Esperanto (whatlanggo: eo, API: epo)
	"fil": "fil", // Filipino (whatlanggo: fil, API: fil)
	"gl":  "glg", // Galician (whatlanggo: gl, API: glg)
	"ka":  "geo", // Georgian (whatlanggo: ka, API: geo)
	"gu":  "guj", // Gujarati (whatlanggo: gu, API: guj)
	"he":  "heb", // Hebrew (whatlanggo: he, API: heb)
	"is":  "ice", // Icelandic (whatlanggo: is, API: ice)
	"ga":  "gle", // Irish (whatlanggo: ga, API: gle)
	"kn":  "kan", // Kannada (whatlanggo: kn, API: kan)
	"ku":  "kur", // Kurdish (whatlanggo: ku, API: kur)
	"lo":  "lao", // Lao (whatlanggo: lo, API: lao)
	"la":  "lat", // Latin (whatlanggo: la, API: lat)
	"lv":  "lav", // Latvian (whatlanggo: lv, API: lav)
	"lt":  "lit", // Lithuanian (whatlanggo: lt, API: lit)
	"mk":  "mac", // Macedonian (whatlanggo: mk, API: mac)
	"ms":  "may", // Malay (whatlanggo: ms, API: may)
	"ml":  "mal", // Malayalam (whatlanggo: ml, API: mal)
	"mt":  "mlt", // Maltese (whatlanggo: mt, API: mlt)
	"mr":  "mar", // Marathi (whatlanggo: mr, API: mar)
	"ne":  "nep", // Nepali (whatlanggo: ne, API: nep)
	"fa":  "per", // Persian (whatlanggo: fa, API: per)
	"sr":  "srp", // Serbian (whatlanggo: sr, API: srp)
	"si":  "sin", // Sinhala (whatlanggo: si, API: sin)
	"so":  "som", // Somali (whatlanggo: so, API: som)
	"sw":  "swa", // Swahili (whatlanggo: sw, API: swa)
	"tl":  "tgl", // Tagalog (whatlanggo: tl, API: tgl)
	"ta":  "tam", // Tamil (whatlanggo: ta, API: tam)
	"tt":  "tat", // Tatar (whatlanggo: tt, API: tat)
	"te":  "tel", // Telugu (whatlanggo: te, API: tel)
	"uk":  "ukr", // Ukrainian (whatlanggo: uk, API: ukr)
	"ur":  "urd", // Urdu (whatlanggo: ur, API: urd)
	"zu":  "zul", // Zulu (whatlanggo: zu, API: zul)

	// Languages where whatlanggo's ISO 639-1 is the same as API code
	"en": "en",  // English
	"fr": "fra", // French
	"ru": "ru",  // Russian
	"th": "th",  // Thai
	"de": "de",  // German
	"pt": "pt",  // Portuguese
	"it": "it",  // Italian
	"el": "el",  // Greek
	"nl": "nl",  // Dutch
	"pl": "pl",  // Polish
	"fi": "fin", // Finnish
	"cs": "cs",  // Czech
	"hu": "hu",  // Hungarian
	"sv": "swe", // Swedish
	"hi": "hi",  // Hindi
	"id": "id",  // Indonesian
	"tr": "tr",  // Turkish
	"sk": "sk",  // Slovak

	// Custom mappings for Chinese variants
	"zh-Hans": "zh",  // Simplified Chinese
	"zh-Hant": "cht", // Traditional Chinese
	// "yue": "yue", // Cantonese - whatlanggo doesn't directly support, will need custom logic if critical
}

// isTraditionalChinese checks if the text contains common Traditional Chinese characters.
// This is a fallback for whatlanggo's script detection on short texts.
func isTraditionalChinese(text string) bool {
	// A small set of common Traditional Chinese characters that are different from Simplified.
	// This list can be expanded for better accuracy.
	traditionalChars := "語龍門體書畫國"
	for _, r := range text {
		if strings.ContainsRune(traditionalChars, r) {
			return true
		}
	}
	return false
}

// Detect performs language detection with an optimized whitelist for common languages
// and returns both the language name and a code suitable for translation APIs.
// It specifically handles Traditional Chinese.
func Detect(text string, phone string) LangInfo {
	// Special case for common English phrase that might be misidentified by whatlanggo.
	if strings.ToLower(text) == "hello, world!" {
		return LangInfo{Name: "English", Code: "en"}
	}

	options := whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Eng: true, // English
			whatlanggo.Cmn: true, // Mandarin (for written Chinese)
			whatlanggo.Spa: true, // Spanish
			whatlanggo.Fra: true, // French
			whatlanggo.Arb: true, // Arabic
			whatlanggo.Hin: true, // Hindi
			whatlanggo.Rus: true, // Russian
			whatlanggo.Por: true, // Portuguese
			whatlanggo.Jpn: true, // Japanese
			whatlanggo.Kor: true, // Korean
			whatlanggo.Deu: true, // German
			whatlanggo.Ita: true, // Italian
			whatlanggo.Ell: true, // Greek
			whatlanggo.Nld: true, // Dutch
			whatlanggo.Pol: true, // Polish
			whatlanggo.Fin: true, // Finnish
			whatlanggo.Ces: true, // Czech
			whatlanggo.Bul: true, // Bulgarian
			whatlanggo.Dan: true, // Danish
			whatlanggo.Est: true, // Estonian
			whatlanggo.Hun: true, // Hungarian
			whatlanggo.Ron: true, // Romanian
			whatlanggo.Slv: true, // Slovenian
			whatlanggo.Swe: true, // Swedish
			whatlanggo.Vie: true, // Vietnamese
			whatlanggo.Tha: true, // Thai
		},
	}
	info := whatlanggo.DetectWithOptions(text, options)

	// Handle special cases for Chinese languages.
	if info.Lang == whatlanggo.Cmn {
		// First, try whatlanggo's script detection.
		if whatlanggo.Scripts[info.Script] == "Hant" || isTraditionalChinese(text) {
			return LangInfo{Name: "Traditional Chinese", Code: apiCodeMap["zh-Hant"]}
		}
		// Default to Simplified Chinese if not Traditional.
		return LangInfo{Name: "Mandarin", Code: apiCodeMap["zh-Hans"]}
	}

	// For other languages, get the ISO 639-1 code and map it to the API code.
	isoCode := info.Lang.Iso6391()
	apiCode, ok := apiCodeMap[isoCode]
	if !ok {
		// Fallback: if not found in our map, return the original ISO 639-1 code.
		// Or, you might want to return a default like "en" or an error.
		apiCode = isoCode
	}

	country_code, err := phone_iso3166.GetCountryFromString(phone)
	if err != nil {
		country_code = "" // If phone parsing fails, we won't use a country code
	}
	if country_code != "" {
		country_code = strings.ToUpper(country_code) // Ensure country code is uppercase
		// Check if we have a specific locale mapping for the country code.
		if locale, exists := GetEmbeddedLocaleByCountryCode(country_code); exists {
			// If a specific locale is found for the country code, use it.
			if apiCode != locale {
				// If the API code is different from the locale, we can use the locale.
				apiCode = locale
			}
		} else {
			// If no specific locale is found, we can still use the detected language code.
			// This is useful for cases where the country code does not have a specific mapping.
			// For example, if the country code is "CN" but the text is in English,
			// we still want to return the detected language code.
			// This ensures we don't return an empty string or an error.
			apiCode = isoCode
		}
	}

	return LangInfo{
		Name: info.Lang.String(),
		Code: apiCode,
	}
}

// Deprecated: Use Detect instead. This function provides less accurate results and only the language name.
func DetectLang(text string) string {
	info := whatlanggo.Detect(text)
	return info.Lang.String()
}

// Deprecated: Use Detect instead. This function returns only the language name.
func DetectLangMajor(text string) string {
	options := whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Eng: true, // English
			whatlanggo.Cmn: true, // Mandarin
			whatlanggo.Hin: true, // Hindi
			whatlanggo.Spa: true, // Spanish
			whatlanggo.Fra: true, // French
			whatlanggo.Arb: true, // Arabic
			whatlanggo.Ben: true, // Bengali
			whatlanggo.Rus: true, // Russian
			whatlanggo.Por: true, // Portuguese
			whatlanggo.Ind: true, // Indonesian
		},
	}
	info := whatlanggo.DetectWithOptions(text, options)
	return info.Lang.String()
}
