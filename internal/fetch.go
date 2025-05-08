package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	abci "github.com/cometbft/cometbft/abci/types"
)

// BlockResponse mirrors the JSON returned by the heimdall-v2 /block?height=<height> endpoint.
type BlockResponse struct {
	Result struct {
		Block struct {
			Header struct {
				ChainID string `json:"chain_id"`
			}
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
		return "", fmt.Errorf("unexpected status fetching VE: %s", resp.Status)
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

// FetchChainID fetches the chain_id from the block at the specified height.
func FetchChainID(httpClient *http.Client, endpoint string, height int64) (string, error) {
	url := fmt.Sprintf("%s/block?height=%d", endpoint, height)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status fetching chainID: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var br BlockResponse
	if err := json.Unmarshal(body, &br); err != nil {
		return "", err
	}

	return br.Result.Block.Header.ChainID, nil
}

// FetchTotalVotingPower queries the endpoint at /stake/total-power and returns the total voting power as an int64.
func FetchTotalVotingPower(httpClient *http.Client, endpoint string) (int64, error) {
	url := fmt.Sprintf("%s/stake/total-power", endpoint)
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status fetching total-power: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var totalPowerResponse struct {
		TotalPower string `json:"total_power"`
	}
	if err := json.Unmarshal(body, &totalPowerResponse); err != nil {
		return 0, err
	}

	power, err := strconv.ParseInt(totalPowerResponse.TotalPower, 10, 64)
	if err != nil {
		return 0, err
	}

	return power, nil
}
