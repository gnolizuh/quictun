package main

// Config for echo server
type Config struct {
	Listen  string `json:"listen"`
	Quiet   bool   `json:"quiet"`
}