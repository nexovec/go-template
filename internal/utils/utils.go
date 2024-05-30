package utils

import (
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v3"
)

func Initialize() error {
	// TODO:
	return nil
}

func CountTrues(arr ...bool) int {
	count := 0
	for _, b := range arr {
		if b {
			count++
		}
	}
	return count
}

func Must[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}
	return result
}

func GetCaller() string {
	pc, file, line, ok := runtime.Caller(2)
	if pc == 0 || !ok {
		return ""
	}
	return file + ":" + fmt.Sprint(line)
}

func IsAjaxRequest(c fiber.Ctx) bool {
	return c.Get("HX-Request") == "true"
}

func GetPackageName(filename string) string {
	// Parse the package name from the filename
	// This is just a basic implementation
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '/' {
			return filename[i+1:]
		}
	}
	return filename
}
