package internal

import (
	"encoding/hex"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
)

// PrintExtendedCommitInfo prints ExtendedCommitInfo in a vertical table layout.
func PrintExtendedCommitInfo(height int64, info *abci.ExtendedCommitInfo) {
	fmt.Println()
	fmt.Println("================ Extended Commit Info ================")
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Round: %d\n", info.Round)
	fmt.Println("======================================================")
	fmt.Println()

	for i, v := range info.Votes {
		fmt.Printf("Vote %d\n", i+1)
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Validator: %s\n", hex.EncodeToString(v.Validator.Address))
		fmt.Printf("Power: %d\n", v.Validator.Power)
		fmt.Printf("BlockIdFlag: %s\n", v.BlockIdFlag.String())
		fmt.Printf("VoteExtension: %s\n", hex.EncodeToString(v.VoteExtension))
		fmt.Printf("ExtensionSignature: %s\n", hex.EncodeToString(v.ExtensionSignature))
		fmt.Printf("NonRpVoteExtension: %s\n", hex.EncodeToString(v.NonRpVoteExtension))
		fmt.Printf("NonRpExtensionSignature: %s\n", hex.EncodeToString(v.NonRpExtensionSignature))
		fmt.Println("------------------------------------------------------")
		fmt.Println()
	}
}
