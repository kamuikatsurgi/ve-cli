package config

import (
	"net/http"
)

var (
	ChainID  string
	Endpoint string

	HTTPClient *http.Client

	StartHeight int64
	EndHeight   int64
)
