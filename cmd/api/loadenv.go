package main

import "github.com/joho/godotenv"

// loads environment files
func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
