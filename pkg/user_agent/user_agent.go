package useragent

import "strings"

var browsers = map[string]string{
	"Opera":   "Opera",
	"OPR":     "Opera",
	"Edg":     "Edge",
	"Chrome":  "Chrome",
	"Safari":  "Safari",
	"Firefox": "Firefox",
	"MSIE":    "IE",
}

// Analysis user-agent to indentify browser
func indentifyBrowser(userAgent string) (browserName string) {
	for _, browserIdentifier := range browsers {
		if strings.Contains(userAgent, browserIdentifier) {
			return browserIdentifier
		}
	}
	return browserName
}

// Analysis user-agent to indentify version of the browser
func indentifyBrowserVersion(userAgent string, browserName string) (browserVersion string) {
	segments := strings.Split(userAgent, " ")
	for _, segment := range segments {
		if strings.HasPrefix(segment, browserName) {
			version := strings.Split(segment, "/")
			if len(version) > 1 {
				browserVersion = version[1]
			}
			break
		}
	}
	return browserVersion
}

// Analysis user-agent to get browser info
func GetDeviceInfo(userAgent string) (browserName, browserVersion string) {
	if len(userAgent) == 0 {
		return "", ""
	}

	browserName = indentifyBrowser(userAgent)

	if len(browserName) == 0 {
		return "", ""
	}

	browserVersion = indentifyBrowserVersion(userAgent, browserName)
	return browserName, browserVersion
}
