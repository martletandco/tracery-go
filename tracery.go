package tracery

import "regexp"

// Flatten resolves a grammer tree
func Flatten(rule string) string {
	actionRe := regexp.MustCompile(`^\[(.*?):(.*?)\]`)
	plainRe := regexp.MustCompile(`^([^\[#]+)`)
	tagRe := regexp.MustCompile(`^#([^#]+)#`)

	ctx := make(map[string]string)
	out := ""
	var index []int
	for {
		if len(rule) == 0 {
			break
		}
		index = actionRe.FindStringIndex(rule)
		if index != nil {
			match := actionRe.FindStringSubmatch(rule[index[0]:index[1]])
			ctx[match[1]] = match[2]
			rule = rule[index[1]:]
			continue
		}
		index = tagRe.FindStringIndex(rule)
		if index != nil {
			match := tagRe.FindStringSubmatch(rule[index[0]:index[1]])
			value, ok := ctx[match[1]]
			if ok {
				out = out + value
			}
			rule = rule[index[1]:]
			continue
		}
		index = plainRe.FindStringIndex(rule)
		if index != nil {
			match := plainRe.FindStringSubmatch(rule[index[0]:index[1]])
			if match != nil {
				out = out + match[1]
			}
			rule = rule[index[1]:]
			continue
		}
	}
	return out
}
