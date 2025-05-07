package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	abci "github.com/cometbft/cometbft/abci/types"
)

// BlockResponse mirrors the JSON returned by the heimdall-v2 /block?height=<height> endpoint.
type BlockResponse struct {
	Result struct {
		Block struct {
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	} `json:"result"`
}

// FetchAndDecodeVEs retrieves and decodes the VEs for the given block range.
func FetchAndDecodeVEs(httpClient *http.Client, endpoint string, start, end int64) ([]*abci.ExtendedCommitInfo, error) {
	var results []*abci.ExtendedCommitInfo
	for h := start; h <= end; h++ {
		veStr, err := FetchVE(httpClient, endpoint, h)
		if err != nil {
			return nil, err
		}
		if veStr == "" {
			continue
		}

		veBytes, err := base64.StdEncoding.DecodeString(veStr)
		if err != nil {
			return nil, err
		}

		decoded, err := DecodeVE(veBytes)
		if err != nil {
			return nil, err
		}
		results = append(results, decoded)
	}
	return results, nil
}

// FetchAndDecodeVE fetches and decodes the VE from the given block height,
func FetchAndDecodeVE(httpClient *http.Client, endpoint string, height int64) (*abci.ExtendedCommitInfo, error) {
	veStr, err := FetchVE(httpClient, endpoint, height)
	if err != nil {
		return nil, err
	}
	if veStr == "" {
		return nil, err
	}

	veBytes, err := base64.StdEncoding.DecodeString(veStr)
	if err != nil {
		return nil, err
	}

	decoded, err := DecodeVE(veBytes)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// FetchVE retrieves the VE (as a Base64 string) at the specified block height.
func FetchVE(httpClient *http.Client, endpoint string, height int64) (string, error) {
	resp, err := httpClient.Get(fmt.Sprintf("%s/block?height=%d", endpoint, height))
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("endpoint returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var br BlockResponse
	if err := json.Unmarshal(body, &br); err != nil {
		return "", err
	}

	if len(br.Result.Block.Data.Txs) == 0 {
		return "", nil
	}
	return br.Result.Block.Data.Txs[0], nil
}

// DecodeVE unmarshals the VE bytes into an ExtendedCommitInfo protobuf.
func DecodeVE(bz []byte) (*abci.ExtendedCommitInfo, error) {
	var info abci.ExtendedCommitInfo
	if err := info.Unmarshal(bz); err != nil {
		return nil, err
	}
	return &info, nil
}
