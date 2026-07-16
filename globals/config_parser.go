package globals

// * NOTE: THIS IS ALL LIBRARY CODE, INTENDED TO BE REMOVED FROM THIS REPO IN THE FUTURE.
// *       THIS IS ONLY HERE FOR NOW SO I CAN PLAY AROUND WITH THE IDEA.

// *       MESSING AROUND WITH THIS BECAUSE I DIDN'T REALLY LIKE THE WAY EXISTING CONFIG
// *       PARSERS WORKED, THEY ALL HAD WEIRD QUIRKS LIKE ODD SEMANTICS FOR STRUCT TAGS,
// *       COULDN'T HANDLE COMPLEX SLICES AND MAPS CLEANLY, ETC.

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type fieldTagOptions struct {
	optional        bool
	defaultValue    string
	hasDefault      bool
	envNameOverride string
}

func parseFieldTag(tag string) fieldTagOptions {
	options := fieldTagOptions{}
	if tag == "" {
		return options
	}

	if i := strings.Index(tag, "default:"); i != -1 {
		options.defaultValue = tag[i+8:]
		options.hasDefault = true
		tag = strings.TrimSuffix(tag[:i], ",")
	}

	for _, part := range strings.Split(tag, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if part == "optional" {
			options.optional = true
		} else if strings.HasPrefix(part, "env:") {
			options.envNameOverride = strings.TrimPrefix(part, "env:")
		}
	}

	return options
}

type ConfigParser[T any] struct {
	config                   T
	prefix                   string
	initialisms              map[string]bool
	allowedPluralInitialisms map[string]bool
}

func (cp *ConfigParser[T]) toPascalCase(str string) string {
	lowercase := strings.ToLower(str)
	words := strings.Split(lowercase, "_")

	var pascalCase strings.Builder

	for _, word := range words {
		pascalCase.WriteString(cp.capitalizeWord(word))
	}

	return pascalCase.String()
}

func (cp *ConfigParser[T]) capitalizeWord(word string) string {
	if word == "" {
		return word
	}

	upper := strings.ToUpper(word)

	if cp.initialisms[upper] {
		return upper
	}

	endsWithS := strings.HasSuffix(upper, "S")
	withoutS := strings.TrimSuffix(upper, "S")

	if cp.initialisms[withoutS] && endsWithS && cp.allowedPluralInitialisms[withoutS] {
		return withoutS + "s"
	} else if cp.initialisms[withoutS] {
		return upper
	}

	r, n := utf8.DecodeRuneInString(word)
	if r == utf8.RuneError && n == 0 {
		return word
	}

	return string(unicode.ToUpper(r)) + word[n:]
}

func (cp *ConfigParser[T]) toEnvVarName(fieldName string) string {
	var result strings.Builder
	if cp.prefix != "" {
		result.WriteString(cp.prefix)
	}

	runes := []rune(fieldName)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == 's' && i > 0 && unicode.IsUpper(runes[i-1]) {
			result.WriteRune('S')
			continue
		}

		if i > 0 && unicode.IsUpper(r) {
			prevIsLower := unicode.IsLower(runes[i-1])
			nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			nextIsPluralS := i+1 < len(runes) && runes[i+1] == 's' && (i+2 >= len(runes) || unicode.IsUpper(runes[i+2]))

			if prevIsLower || (nextIsLower && !nextIsPluralS) {
				result.WriteRune('_')
			}
		}

		result.WriteRune(unicode.ToUpper(r))
	}

	return result.String()
}

func (cp *ConfigParser[T]) SetPrefix(prefix string) *ConfigParser[T] {
	cp.prefix = prefix + "_"

	return cp
}

func (cp *ConfigParser[T]) AddInitialisms(initialisms map[string]bool) *ConfigParser[T] {
	for key, value := range initialisms {
		cp.initialisms[key] = value
	}

	return cp
}

func (cp *ConfigParser[T]) ParseFromEnv() T {
	envAsPascal := make(map[string]string)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		value := strings.TrimSpace(pair[1])

		if !strings.HasPrefix(key, cp.prefix) {
			continue
		}

		fieldName := cp.toPascalCase(strings.TrimPrefix(key, cp.prefix))

		envAsPascal[fieldName] = value
	}

	v := reflect.ValueOf(cp.config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldValue := v.Field(i)
		fieldOptions := parseFieldTag(field.Tag.Get("envconf"))
		errors := make([]string, 0)
		warnings := make([]string, 0)

		envVarName := cp.toEnvVarName(fieldName)
		envValue, exists := envAsPascal[fieldName]

		if fieldOptions.envNameOverride != "" {
			envVarName = fieldOptions.envNameOverride
			envValue, exists = os.LookupEnv(envVarName)
		}

		if !exists && fieldOptions.hasDefault {
			envValue = fieldOptions.defaultValue
			exists = true

			warnings = append(warnings, fmt.Sprintf("Optional field %s does not have a corresponding %s environment variable. Using default value \"%s\"", fieldName, envVarName, envValue))
		}

		if exists && fieldValue.CanSet() {
			fieldType := fieldValue.Type()
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(envValue)
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(envValue); err == nil {
					fieldValue.SetBool(boolVal)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intVal, err := strconv.ParseInt(envValue, 10, fieldValue.Type().Bits()); err == nil {
					fieldValue.SetInt(intVal)
				} else {
					errors = append(errors, fmt.Sprintf("Error parsing %s: %v", envVarName, err))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if uintVal, err := strconv.ParseUint(envValue, 10, fieldValue.Type().Bits()); err == nil {
					fieldValue.SetUint(uintVal)
				} else {
					errors = append(errors, fmt.Sprintf("Error parsing %s: %v", envVarName, err))
				}
			case reflect.Float32, reflect.Float64:
				if floatVal, err := strconv.ParseFloat(envValue, fieldValue.Type().Bits()); err == nil {
					fieldValue.SetFloat(floatVal)
				} else {
					errors = append(errors, fmt.Sprintf("Error parsing %s: %v", envVarName, err))
				}
			case reflect.Slice, reflect.Map:
				sliceOrMap := reflect.New(fieldType)
				if err := json.Unmarshal([]byte(envValue), sliceOrMap.Interface()); err != nil {
					errors = append(errors, fmt.Sprintf("Error parsing %s: %v", envVarName, err))
				}
				fieldValue.Set(sliceOrMap.Elem())
			}
		} else if !exists {
			if !fieldOptions.optional {
				errors = append(errors, fmt.Sprintf("Required field %s does not have a corresponding %s environment variable", fieldName, envVarName))
			} else if !fieldOptions.hasDefault {
				warnings = append(warnings, fmt.Sprintf("Optional field %s does not have a corresponding %s environment variable and no default value. Skipping", fieldName, envVarName))
			}
		}

		if len(warnings) != 0 {
			for _, warning := range warnings {
				fmt.Println(warning)
			}
		}

		if len(errors) != 0 {
			for _, err := range errors {
				fmt.Println(err)
			}

			os.Exit(0)
		}
	}

	return cp.config
}

func NewConfigParser[T any](config T) *ConfigParser[T] {
	return &ConfigParser[T]{
		config: config,
		initialisms: map[string]bool{ // * https://go.googlesource.com/lint/+/818c5a804067/lint.go#767
			"ACL":   true,
			"API":   true,
			"ASCII": true,
			"CPU":   true,
			"CSS":   true,
			"DNS":   true,
			"EOF":   true,
			"GUID":  true,
			"HTML":  true,
			"HTTP":  true,
			"HTTPS": true,
			"ID":    true,
			"IP":    true,
			"JSON":  true,
			"LHS":   true,
			"QPS":   true,
			"RAM":   true,
			"RHS":   true,
			"RPC":   true,
			"SLA":   true,
			"SMTP":  true,
			"SQL":   true,
			"SSH":   true,
			"TCP":   true,
			"TLS":   true,
			"TTL":   true,
			"UDP":   true,
			"UI":    true,
			"UID":   true,
			"UUID":  true,
			"URI":   true,
			"URL":   true,
			"UTF8":  true,
			"VM":    true,
			"XML":   true,
			"XMPP":  true,
			"XSRF":  true,
			"XSS":   true,
			"NEX":   true, // * Start of our custom ones
			"GRPC":  true,
			"AES":   true,
			"PID":   true,
		},
		allowedPluralInitialisms: map[string]bool{
			"API":  true,
			"GUID": true,
			"ID":   true,
			"IP":   true,
			"UUID": true,
			"URI":  true,
			"URL":  true,
		},
	}
}
