package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/kamuikatsurgi/ve-cli/config"
	"github.com/kamuikatsurgi/ve-cli/internal"
)

// FetchAndDecode handles the core logic of fetching and decoding blocks.
func FetchAndDecode() error {
	if config.StartHeight == config.EndHeight {
		fmt.Printf("Fetching and decoding VE for block height %d...\n", config.StartHeight)
		resp, err := internal.FetchAndDecodeVE(config.HTTPClient, config.Endpoint, config.StartHeight)
		if err != nil {
			return err
		}
		err = internal.DecodeAndPrintExtendedCommitInfo(config.StartHeight, resp)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully fetched and decoded VE at height %d.\n", config.StartHeight)
	} else {
		fmt.Printf("Fetching and decoding VEs for blocks from height %d to %d...\n", config.StartHeight, config.EndHeight)
		resp, err := internal.FetchAndDecodeVEs(config.HTTPClient, config.Endpoint, config.StartHeight, config.EndHeight)
		if err != nil {
			return err
		}
		for i, ve := range resp {
			err = internal.DecodeAndPrintExtendedCommitInfo(config.StartHeight+int64(i), ve)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Successfully fetched and decoded VEs from height %d to %d.\n", config.StartHeight, config.EndHeight)
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:          "ve-cli",
	Short:        "Extract and decode VEs from Heimdall-v2 blocks",
	Long:         "Use the 'block' subcommand for a single block, or the 'blocks' subcommand to process a range of blocks.",
	SilenceUsage: true,
}

var blockCmd = &cobra.Command{
	Use:   "block <height>",
	Short: "Process a single block height",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		height, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid block height: %v", err)
		}
		if height < 0 {
			return fmt.Errorf("block height cannot be negative")
		}
		if height == 1 {
			return fmt.Errorf("vote extensions are not enabled at block height 1")
		}
		config.StartHeight = height
		config.EndHeight = height

		chainID, err := internal.FetchChainID(config.HTTPClient, config.Endpoint, config.StartHeight)
		if err != nil {
			return fmt.Errorf("failed to fetch chain ID: %v", err)
		}
		config.ChainID = chainID

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return FetchAndDecode()
	},
}

var blocksCmd = &cobra.Command{
	Use:   "blocks <start-height> <end-height>",
	Short: "Process a range of block heights",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		start, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid start height: %v", err)
		}
		end, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid end height: %v", err)
		}
		if start < 0 {
			return fmt.Errorf("start height cannot be negative")
		}
		if end < 0 {
			return fmt.Errorf("end height cannot be negative")
		}
		if start > end {
			return fmt.Errorf("start height (%d) cannot be greater than end height (%d)", start, end)
		}
		if start == 1 || end == 1 {
			return fmt.Errorf("vote extensions are not enabled at block height 1")
		}
		config.StartHeight = start
		config.EndHeight = end

		chainID, err := internal.FetchChainID(config.HTTPClient, config.Endpoint, config.StartHeight)
		if err != nil {
			return fmt.Errorf("failed to fetch chain ID: %v", err)
		}
		config.ChainID = chainID

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return FetchAndDecode()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(blockCmd, blocksCmd)
	rootCmd.PersistentFlags().StringVarP(&config.Endpoint, "endpoint", "e", "http://localhost:26657", "Heimdall-v2 RPC URL")
}

func initConfig() {
	config.HTTPClient = &http.Client{Timeout: 25 * time.Second}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
