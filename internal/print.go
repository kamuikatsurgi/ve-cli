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

	PrintHeader(height, info.Round)

	for i, v := range info.Votes {
		err := PrintVote(height, i+1, v, voteExts[i])
		if err != nil {
			return err
		}
	}

	err := PrintSummary(info.Votes, voteExts)
	if err != nil {
		return err
	}

	return nil
}

func PrintHeader(height int64, round int32) {
	fmt.Println()
	fmt.Println("================ Extended Commit Info ================")
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Round: %d\n", round)
	fmt.Println("======================================================")
	fmt.Println()
}

func PrintVote(height int64, index int, voteInfo abci.ExtendedVoteInfo, voteExt sidetxs.VoteExtension) error {
	fmt.Printf("Vote %d:\n", index)
	fmt.Println("------------------------------------------------------")

	PrintValidatorInfo(voteInfo)
	PrintVoteExtensionInfo(voteExt)
	err := PrintNonRpVoteExtAndSignatures(height, voteInfo)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------")
	fmt.Println()

	return nil
}

func PrintValidatorInfo(v abci.ExtendedVoteInfo) {
	fmt.Printf("Validator: %s\n", hex.EncodeToString(v.Validator.Address))
	fmt.Printf("Power: %d\n", v.Validator.Power)
	fmt.Printf("BlockIdFlag: %s\n", v.BlockIdFlag.String())
}

func PrintVoteExtensionInfo(voteExt sidetxs.VoteExtension) {
	fmt.Println("VoteExtension:")
	fmt.Printf("  BlockHash: %s\n", hex.EncodeToString(voteExt.BlockHash))
	fmt.Printf("  Height: %d\n", voteExt.Height)

	if len(voteExt.SideTxResponses) == 0 {
		fmt.Println("  SideTxResponses: []")
	} else {
		fmt.Println("  SideTxResponses:")
		for j, resp := range voteExt.SideTxResponses {
			fmt.Printf("    Response %d:\n", j+1)
			fmt.Printf("      TxHash: %s\n", hex.EncodeToString(resp.TxHash))
			fmt.Printf("      Result: %s\n", resp.Result.String())
		}
	}

	if voteExt.MilestoneProposition != nil {
		mp := voteExt.MilestoneProposition
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

func PrintNonRpVoteExtAndSignatures(height int64, v abci.ExtendedVoteInfo) error {
	fmt.Printf("ExtensionSignature: %s\n", hex.EncodeToString(v.ExtensionSignature))
	err := PrintNonRpVoteExtension(height, v.NonRpVoteExtension)
	if err != nil {
		return err
	}
	fmt.Printf("NonRpExtensionSignature: %s\n", hex.EncodeToString(v.NonRpExtensionSignature))

	return nil
}
