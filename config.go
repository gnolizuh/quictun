package main

// Config for server
type Config struct {
	Listen  string `json:"listen"`
	Target  string `json:"target"`
	Timeout int    `json:"mtu"`
}