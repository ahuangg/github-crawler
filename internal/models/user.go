package models

type User struct {
	Username       string
	Location      string
	LanguageStats map[string]float64 
}