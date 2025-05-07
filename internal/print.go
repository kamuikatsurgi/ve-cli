package internal

import (
	"encoding/hex"
	"fmt"

	sidetxs "github.com/0xPolygon/heimdall-v2/sidetxs"
	abci "github.com/cometbft/cometbft/abci/types"
	goproto "github.com/cosmos/gogoproto/proto"
)

// DecodeAndPrintExtendedCommitInfo decodes and prints all fields of ExtendedCommitInfo.
func DecodeAndPrintExtendedCommitInfo(height int64, info *abci.ExtendedCommitInfo) error {
	voteExts := make([]sidetxs.VoteExtension, len(info.Votes))
	for i, v := range info.Votes {
		if err := goproto.Unmarshal(v.VoteExtension, &voteExts[i]); err != nil {
			return err
		}
	}

	printHeader(height, info.Round)

	for i, v := range info.Votes {
		printVote(i+1, v, voteExts[i])
	}

	return nil
}

func printHeader(height int64, round int32) {
	fmt.Println()
	fmt.Println("================ Extended Commit Info ================")
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Round: %d\n", round)
	fmt.Println("======================================================")
	fmt.Println()
}

func printVote(index int, voteInfo abci.ExtendedVoteInfo, ext sidetxs.VoteExtension) {
	fmt.Printf("Vote %d:\n", index)
	fmt.Println("------------------------------------------------------")

	printValidatorInfo(voteInfo)
	printVoteExtensionInfo(ext)
	printRawSignatures(voteInfo)

	fmt.Println("------------------------------------------------------")
	fmt.Println()
}

func printValidatorInfo(v abci.ExtendedVoteInfo) {
	fmt.Printf("Validator: %s\n", hex.EncodeToString(v.Validator.Address))
	fmt.Printf("Power: %d\n", v.Validator.Power)
	fmt.Printf("BlockIdFlag: %s\n", v.BlockIdFlag.String())
}

func printVoteExtensionInfo(ext sidetxs.VoteExtension) {
	fmt.Println("VoteExtension:")
	fmt.Printf("  BlockHash: %s\n", hex.EncodeToString(ext.BlockHash))
	fmt.Printf("  Height: %d\n", ext.Height)

	if len(ext.SideTxResponses) == 0 {
		fmt.Println("  SideTxResponses: []")
	} else {
		fmt.Println("  SideTxResponses:")
		for j, resp := range ext.SideTxResponses {
			fmt.Printf("    Response %d:\n", j+1)
			fmt.Printf("      TxHash: %s\n", hex.EncodeToString(resp.TxHash))
			fmt.Printf("      Result: %s\n", resp.Result.String())
		}
	}

	if ext.MilestoneProposition != nil {
		mp := ext.MilestoneProposition
		fmt.Println("  MilestoneProposition:")
		for k, bh := range mp.BlockHashes {
			fmt.Printf("    BlockHash[%d]: %s\n", k, hex.EncodeToString(bh))
		}
		fmt.Printf("    StartBlockNumber: %d\n", mp.StartBlockNumber)
		fmt.Printf("    ParentHash: %s\n", hex.EncodeToString(mp.ParentHash))
	} else {
		fmt.Println("  MilestoneProposition: nil")
	}
}

func printRawSignatures(v abci.ExtendedVoteInfo) {
	fmt.Printf("ExtensionSignature: %s\n", hex.EncodeToString(v.ExtensionSignature))
	fmt.Printf("NonRpVoteExtension: %s\n", hex.EncodeToString(v.NonRpVoteExtension))
	fmt.Printf("NonRpExtensionSignature: %s\n", hex.EncodeToString(v.NonRpExtensionSignature))
}
