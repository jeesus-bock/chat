package utils

func ContainsStr(ss []string, s string) bool {
	for _, str := range ss {
		if s == str {
			return true
		}
	}
	return false
}
