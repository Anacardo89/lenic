package parse

import "regexp"

func ParseAtString(input string) []string {
	re := regexp.MustCompile(`@[\w]+`)
	return re.FindAllString(input, -1)
}
