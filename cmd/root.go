package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	blockFlag  bool
	blocksFlag bool

	StartHeight int64
	EndHeight   int64
)

var rootCmd = &cobra.Command{
	Use:   "ve-cli --block <height> | --blocks <start-height> <end-height>",
	Short: "Extract and decode VEs from Heimdall-v2 blocks",
	Long:  `Use --block <height> for a single block, or --blocks <start-height> <end-height> to process a range of blocks.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if blockFlag && blocksFlag {
			return fmt.Errorf("cannot use both --block and --blocks together")
		}
		if !blockFlag && !blocksFlag {
			return fmt.Errorf("must provide either --block or --blocks")
		}
		if blockFlag && len(args) != 1 {
			return fmt.Errorf("you must pass exactly one argument with --block: <height>")
		}
		if blocksFlag && len(args) != 2 {
			return fmt.Errorf("you must pass exactly two arguments with --blocks: <start> <end>")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if blockFlag {
			StartHeight, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println("❌ Invalid block height:", err)
				os.Exit(1)
			}
			EndHeight = StartHeight

			fmt.Printf("Fetching and decoding block height %d...\n", StartHeight)
		} else {
			StartHeight, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println("❌ Invalid start height:", err)
				os.Exit(1)
			}
			EndHeight, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				fmt.Println("❌ Invalid end height:", err)
				os.Exit(1)
			}
			if StartHeight > EndHeight {
				fmt.Println("❌ Start height cannot be greater than end height.")
				os.Exit(1)
			}

			fmt.Printf("Fetching and decoding from block height %d to %d...\n", StartHeight, EndHeight)
		}
	},
}

func Execute() {
	rootCmd.Flags().BoolVar(&blockFlag, "block", false, "Specify a single block height")
	rootCmd.Flags().BoolVar(&blocksFlag, "blocks", false, "Enable range mode with 2 args: <start-height> <end-height>")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ CLI execution error:", err)
		os.Exit(1)
	}
}
