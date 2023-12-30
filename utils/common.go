package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func ParseStringToBool(s string) (bool, error) {
	// Parse the string into a boolean
	result, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}

	return result, nil
}

func GenerateCode(prefix string, length int) string {
	// Use the current timestamp and a random number to generate a unique identifier
	rand.Seed(time.Now().UnixNano())
	uniqueID := rand.Intn(9999) // You can adjust the range as needed

	// Create the product code by combining the prefix, unique identifier, and possibly a suffix
	code := fmt.Sprintf("%s%04d", prefix, uniqueID)

	// Optionally, you can add a suffix
	// productCode := fmt.Sprintf("%s%04d%s", prefix, uniqueID, suffix)

	return code
}

func ParsePoitnerToString(ptr *string) string {
	if ptr != nil {
		return *ptr
	} else {
		return ""
	}
}

func ParseStringToIn64(s string) (int64, error) {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func ParseIn64ToPointer(n int64) *int64 {
	if n == 0 {
		return nil
	} else {
		return &n
	}
}

func InterfaceToString(value interface{}) string {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		return v.String()
	}
	return ""
}

func parsePointerToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}


