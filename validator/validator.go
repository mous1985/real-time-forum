package validator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(data interface{}) error {
	t := reflect.TypeOf(data)
	structName := t.String()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kind := field.Type.Kind()
		data := reflect.ValueOf(data).Field(i).Interface()

		flagsStr, ok := field.Tag.Lookup("validator")
		if !ok && kind != reflect.Struct {
			continue
		}
		flags := parseFlags(flagsStr)

		if kind == reflect.Struct {
			err := Validate(data)
			if err != nil {
				return err
			}
		} else if kind == reflect.Array || kind == reflect.Slice {
			s := reflect.ValueOf(data)
			if hasRequiredFlag(flags) && s.Len() == 0 {
				return newError(field.Name, "is required")
			}
			for i := 0; i < s.Len(); i++ {
				d := s.Index(i).Interface()
				err := checkAll(d, flags, field.Name, structName)
				if err != nil {
					return err
				}
			}
		} else {
			err := checkAll(data, flags, field.Name, structName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkAll(data interface{}, flags []string, fieldName string, structName string) error {
	for _, flag := range flags {
		keyValue := strings.Split(flag, "=")

		switch len(keyValue) {
		case 1:
			key := keyValue[0]
			switch key {
			case "required":
				if err := checkRequired(data); err != nil {
					return newError(fieldName, err.Error())
				}
			case "password":
				if err := checkPassword(fmt.Sprintf("%s", data)); err != nil {
					return newError(fieldName, err.Error())
				}
			case "email":
				if err := checkEmail(fmt.Sprintf("%s", data)); err != nil {
					return err
				}
			case "username":
				if err := checkUsername(fmt.Sprintf("%s", data)); err != nil {
					return err
				}
			default:
				log.Fatalf("'%v' invalid key at %v", key, structName)
			}

		case 2:
			key, value := keyValue[0], keyValue[1]
			switch key {
			case "min":
				if err := checkMin(data, value); err != nil {
					return newError(fieldName, err.Error())
				}
			case "max":
				if err := checkMax(data, value); err != nil {
					return newError(fieldName, err.Error())
				}
			default:
				return fmt.Errorf("%v invalid key", key)
			}
		default:
			log.Fatalf("'%v' invalid format at %v", flag, structName)
		}
	}
	return nil
}

func checkRequired(data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Int:
		if v.Int() == 0 {
			return newError("is required")
		}
	case reflect.String:
		if len(strings.TrimSpace(v.String())) == 0 {
			return newError("is required")
		}
	}
	return nil
}

func checkPassword(pass string) error {
	numRegex := regexp.MustCompile(`[0-9]{1}`)
	lowercaseRegex := regexp.MustCompile(`[a-z]{1}`)
	uppercaseRegex := regexp.MustCompile(`[A-Z]{1}`)
	symbolRegex := regexp.MustCompile(`[!@#~$%^&*()+|_]{1}`)

	if !numRegex.MatchString(pass) {
		return newError("must contain at least one number")
	}
	if !lowercaseRegex.MatchString(pass) {
		return newError("must contain at least one lowercase letter")
	}
	if !uppercaseRegex.MatchString(pass) {
		return newError("must contain at least one uppercase letter")
	}
	if !symbolRegex.MatchString(pass) {
		return newError("must contain at least one symbol\n(!, @, #, ~, $, %, ^, &, *, (, ), +, |, _, )")
	}
	return nil
}

func checkEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,24}$`)

	if !emailRegex.MatchString(email) {
		return newError("e-mail is invalid")
	}
	return nil
}

func checkUsername(username string) error {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	if !usernameRegex.MatchString(username) {
		return newError("username is invalid")
	}
	return nil
}
func checkMin(data interface{}, minStr string) error {
	min, err := parseInt(minStr)
	if err != nil {
		return err
	}

	switch v := data.(type) {
	case int:
		if v < min {
			return fmt.Errorf("value (%v) is lower than minimum value (%v)", v, min)
		}
	case string:
		l := len(v)
		if l < min {
			return fmt.Errorf("length (%v) is lower than minimum length (%v)", l, min)
		}
	}
	return nil
}

func checkMax(data interface{}, maxStr string) error {
	max, err := parseInt(maxStr)
	if err != nil {
		return err
	}

	switch v := data.(type) {
	case int:
		if v > max {
			return fmt.Errorf("value (%v) is higher than maximum value (%v)", v, max)
		}
	case string:
		l := len(v)
		if l > max {
			return fmt.Errorf("length (%v) length is higher than maximim length (%v)", l, max)
		}
	}
	return nil
}

func parseInt(intStr string) (int, error) {
	n, err := strconv.Atoi(intStr)
	if err != nil {
		return n, fmt.Errorf("%v is not integer", intStr)
	}
	return n, nil
}

func parseFlags(flags string) []string {
	flags = trimWhiteSpaces(flags)
	flagsArray := strings.Split(flags, ",")
	return flagsArray
}

func trimWhiteSpaces(s string) string {
	whiteSpaces := []string{" ", "\t", "\v", "\n"}
	for _, w := range whiteSpaces {
		s = strings.ReplaceAll(s, w, "")
	}
	return s
}

func newError(errs ...string) error {
	return errors.New(strings.Join(errs, " "))
}

func hasRequiredFlag(flags []string) bool {
	for _, f := range flags {
		if f == "required" {
			return true
		}
	}
	return false
}
