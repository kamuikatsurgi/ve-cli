package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/kamuikatsurgi/ve-cli/internal"
)

var (
	endpoint   string
	httpClient = &http.Client{Timeout: 5 * time.Second}
)

// fetchAndDecode handles the core logic of fetching and decoding blocks.
func fetchAndDecode(start, end int64) error {
	if start == end {
		fmt.Printf("Fetching and decoding VE for block height %d...\n", start)
		resp, err := internal.FetchAndDecodeVE(httpClient, endpoint, start)
		if err != nil {
			return fmt.Errorf("failed to fetch and decode VE at height %d: %w", start, err)
		}
		internal.PrintExtendedCommitInfo(start, resp)
		fmt.Printf("Successfully fetched and decoded VE at height %d.\n", start)
	} else {
		fmt.Printf("Fetching and decoding VEs for blocks from height %d to %d...\n", start, end)
		resp, err := internal.FetchAndDecodeVEs(httpClient, endpoint, start, end)
		if err != nil {
			return fmt.Errorf("failed to fetch and decode VEs from height %d to %d: %w", start, end, err)
		}
		for i, ve := range resp {
			internal.PrintExtendedCommitInfo(start+int64(i), ve)
		}
		fmt.Printf("Successfully fetched and decoded VEs from height %d to %d.\n", start, end)
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "ve-cli",
	Short: "Extract and decode VEs from Heimdall-v2 blocks",
	Long:  "Use the 'block' subcommand for a single block, or the 'blocks' subcommand to process a range of blocks.",
}

var blockCmd = &cobra.Command{
	Use:   "block <height>",
	Short: "Process a single block height",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		height, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid block height: %v", err)
		}
		return fetchAndDecode(height, height)
	},
}

var blocksCmd = &cobra.Command{
	Use:   "blocks <start-height> <end-height>",
	Short: "Process a range of block heights",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid start height: %v", err)
		}
		end, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid end height: %v", err)
		}
		if start > end {
			return fmt.Errorf("start height (%d) cannot be greater than end height (%d)", start, end)
		}
		return fetchAndDecode(start, end)
	},
}

func init() {
	rootCmd.AddCommand(blockCmd, blocksCmd)
	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "http://localhost:26657", "Heimdall-v2 RPC URL")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
