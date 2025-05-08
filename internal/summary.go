package internal

import (
	"encoding/hex"
	"fmt"

	sidetxs "github.com/0xPolygon/heimdall-v2/sidetxs"
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/kamuikatsurgi/ve-cli/config"
)

// PrintSummary prints a compact summary of:
//   - milestone block hash voting power and percentage
//   - side-tx voting power per result and percentage
//   - non-RP vote-extension voting power and percentage
//
// It fetches the total voting power internally.
func PrintSummary(votes []abci.ExtendedVoteInfo, extensions []sidetxs.VoteExtension) error {
	totalPower, err := FetchTotalVotingPower(config.HttpClient, config.HeimdallEndpoint)
	if err != nil {
		return err
	}

	// 1) Tally milestone voting power
	milestoneVP := make(map[string]int64)
	// 2) Tally side-tx voting power
	sideTxVP := make(map[string]map[sidetxs.Vote]int64)
	// 3) Tally non-RP vote-extension voting power
	nonRpVP := make(map[string]int64)

	for i, voteInfo := range votes {
		power := voteInfo.Validator.Power

		// Milestone propositions
		if prop := extensions[i].MilestoneProposition; prop != nil {
			for _, blockHash := range prop.BlockHashes {
				hashKey := hex.EncodeToString(blockHash)
				milestoneVP[hashKey] += power
			}
		}

		// Side-tx responses
		for _, response := range extensions[i].SideTxResponses {
			txKey := hex.EncodeToString(response.TxHash)
			if _, exists := sideTxVP[txKey]; !exists {
				sideTxVP[txKey] = make(map[sidetxs.Vote]int64)
			}
			sideTxVP[txKey][response.Result] += power
		}

		// Non-RP vote extension
		nrpKey := hex.EncodeToString(voteInfo.NonRpVoteExtension)
		nonRpVP[nrpKey] += power
	}

	formatToPercent := func(vp int64) string {
		percent := float64(vp) / float64(totalPower) * 100
		return fmt.Sprintf("%d (%.2f%%)", vp, percent)
	}

	fmt.Println("=============== Summary ===============")

	fmt.Println("Milestone Block Hash Voting Power:")
	if len(milestoneVP) == 0 {
		fmt.Println("  None")
	} else {
		for hash, vp := range milestoneVP {
			fmt.Printf("  %s: %s\n", hash, formatToPercent(vp))
		}
	}
	fmt.Println()

	fmt.Println("Side-Tx Voting Power by Result:")
	if len(sideTxVP) == 0 {
		fmt.Println("  None")
	} else {
		for tx, results := range sideTxVP {
			fmt.Printf("  Tx %s:\n", tx)
			fmt.Printf("    YES:         %s\n", formatToPercent(results[sidetxs.Vote_VOTE_YES]))
			fmt.Printf("    NO:          %s\n", formatToPercent(results[sidetxs.Vote_VOTE_NO]))
			fmt.Printf("    UNSPECIFIED: %s\n", formatToPercent(results[sidetxs.Vote_UNSPECIFIED]))
		}
	}
	fmt.Println()

	fmt.Println("Non-RP Vote-Extension Voting Power:")
	if len(nonRpVP) == 0 {
		fmt.Println("  None")
	} else {
		for ext, vp := range nonRpVP {
			fmt.Printf("  %s: %s\n", ext, formatToPercent(vp))
		}
	}
	fmt.Println("========================================")

	return nil
}
