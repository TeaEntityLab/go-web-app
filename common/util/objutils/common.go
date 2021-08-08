package objutils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func ReplaceFormat(format string, p map[string]interface{}) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "${" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

func DeduplicateStringSlice(list []string) []string {
	var result []string
	hashmap := map[string]bool{}

	for _, item := range list {
		if hashmap[item] {
			continue
		}
		hashmap[item] = true
		result = append(result, item)
	}

	return result
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func StrOrDefault(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}
