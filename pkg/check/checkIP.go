package check

import (
	"regexp"
	"strings"
)

//ipv4RegexPattern matches ipv4 string like "127.0.0.1"
var ipv4RegexPattern = `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`

//ipv6onlyRegexPattern matching string like "2001:da8::"
var ipv6onlyRegexPattern = "^[0-9a-fA-F]*:([0-9a-fA-F]|:)*$"

// ipv6WithSquareBracketPattern matching string like "[2001:da8::]"
var ipv6WithSquareBracketPattern = "^\\[[0-9a-fA-F]*:([0-9a-fA-F]|:)*\\]$"

//portRegexPattern matches port string like "8080"
var portRegexPattern = `^([0-9]|[1-9][0-9]|[1-9][0-9]{2}|[1-9][0-9]{3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`

func IPType(ip string) (bool, string) {
	if ip == "" {
		return false, "IP is empty"
	}
	if strings.Contains(ip, ":") {
		if strings.Contains(ip, "[") && strings.Contains(ip, "]") {
			return checkIPv6WithSquareBracket(ip), "6"
		}
		return checkIPv6WithOutSquareBracket(ip), "6"
	}
	return checkIPv4(ip), "4"
}

//checkIPv4 checks ipv4
func checkIPv4(ip string) bool {
	regexp.MustCompile(ipv4RegexPattern)
	return regexp.MustCompile(ipv4RegexPattern).MatchString(ip)
}

func checkIPv6(ip string) bool {
	if checkIPv6WithSquareBracket(ip) || checkIPv6WithOutSquareBracket(ip) {
		return false
	}
	return true
}
func checkIPv6WithOutSquareBracket(ip string) bool {
	regexp.MustCompile(ipv6onlyRegexPattern)
	return regexp.MustCompile(ipv6onlyRegexPattern).MatchString(ip)
}

func checkIPv6WithSquareBracket(ip string) bool {
	regexp.MustCompile(ipv6WithSquareBracketPattern)
	return regexp.MustCompile(ipv6WithSquareBracketPattern).MatchString(ip)
}

func IP(str string) bool {
	if str == "" {
		return false
	}
	if strings.Contains(str, ":") {
		if strings.Contains(str, "[") && strings.Contains(str, "]") {
			return checkIPv6WithSquareBracket(str)
		}
		return checkIPv6WithOutSquareBracket(str)
	}
	return checkIPv4(str)
}

func Port(str string) bool {
	regexp.MustCompile(portRegexPattern)
	return regexp.MustCompile(portRegexPattern).MatchString(str)
}

func Protocol(str string) bool {
	if str == "" {
		return false
	}
	if str == "tcp" || str == "udp" || str == "both" || str == "raw" {
		return true
	}
	return false
}
