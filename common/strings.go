package common

func Bool2String(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func StripQuotes(s string) string {
	return s[1 : len(s)-1]
}
