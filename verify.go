package main

import (
	"fmt"
	"regexp"
)

func VerifyRequest(request ExecuteRequest) error {
	err := VerifyDirectoryName(request.directory)
	if err != nil {
		return err
	}
	err = VertifyToken(request.token)
	if err != nil {
		return err
	}
	return nil
}

func VerifyDirectoryName(directory string) error {
	if len(directory) > 20 {
		return fmt.Errorf("Directory name is too long")
	}
	re := "^[a-zA-Z\\-0-9]+$"
	regex := regexp.MustCompile(re)
	if !regex.MatchString(directory) {
		return fmt.Errorf("Directory contains invalid characters")
	}
	return nil
}

func VertifyToken(token string) error {
	if token == "" {
		return fmt.Errorf("No token provided")
	}
	if len(token) > 500 {
		return fmt.Errorf("Token is too long")
	}
	return nil
}
