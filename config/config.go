package config

import (
	"net/http"
)

var (
	ChainID          string
	CometEndpoint    string
	HeimdallEndpoint string

	HttpClient *http.Client

	StartHeight int64
	EndHeight   int64
)
