package main

// Config for client
type Config struct {
	LocalAddr  string `json:"localaddr"`
	RemoteAddr string `json:"remoteaddr"`
	Timeout    int    `json:"timeout"`
	Retry      int    `json:"retry"`
	Quiet      bool   `json:"quiet"`
}
