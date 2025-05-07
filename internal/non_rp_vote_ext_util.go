package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	checkpointTypes "github.com/0xPolygon/heimdall-v2/x/checkpoint/types"

	"github.com/kamuikatsurgi/ve-cli/config"
)

var dummyNonRpVoteExtension = []byte("\t\r\n#HEIMDALL-VOTE-EXTENSION#\r\n\t")

// PrintNonRpVoteExtension prints the non-RP vote extension.
func PrintNonRpVoteExtension(height int64, nonRpVoteExt []byte) error {
	dummy, err := IsDummyNonRpVoteExtension(height, nonRpVoteExt)
	if err != nil {
		return err
	}
	if dummy {
		fmt.Printf("NonRpVoteExtension [DUMMY #HEIMDALL-VOTE-EXTENSION#]: %s\n", hex.EncodeToString(nonRpVoteExt))
	} else {
		msg, err := GetCheckpointMsg(nonRpVoteExt)
		if err != nil {
			return err
		}
		fmt.Println("NonRpVoteExtension [CHECKPOINT MSG]:")
		fmt.Printf("  Proposer: %s\n", msg.Proposer)
		fmt.Printf("  StartBlock: %d\n", msg.StartBlock)
		fmt.Printf("  EndBlock: %d\n", msg.EndBlock)
		fmt.Printf("  RootHash: %s\n", hex.EncodeToString(msg.RootHash))
		fmt.Printf("  AccountRootHash: %s\n", hex.EncodeToString(msg.AccountRootHash))
		fmt.Printf("  BorChainId: %s\n", msg.BorChainId)
	}

	return nil
}

// IsDummyNonRpVoteExtension returns true if the given byte slice matches the dummy extension.
func IsDummyNonRpVoteExtension(height int64, nonRpVoteExt []byte) (bool, error) {
	dummyVoteExt, err := GetDummyNonRpVoteExtension(height-1, config.ChainID)
	if err != nil {
		return false, err
	}
	return bytes.Equal(nonRpVoteExt, dummyVoteExt), nil
}

// GetDummyNonRpVoteExtension returns a dummy non-rp vote extension for given height and chain id.
func GetDummyNonRpVoteExtension(height int64, chainID string) ([]byte, error) {
	var buf bytes.Buffer

	writtenBytes, err := buf.Write(dummyNonRpVoteExtension)
	if err != nil {
		return nil, err
	}
	if writtenBytes != len(dummyNonRpVoteExtension) {
		return nil, errors.New("failed to write dummy vote extension")
	}
	if err := buf.WriteByte('|'); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, height); err != nil {
		return nil, err
	}
	if err := buf.WriteByte('|'); err != nil {
		return nil, err
	}
	writtenBytes, err = buf.WriteString(chainID)
	if err != nil {
		return nil, err
	}
	if writtenBytes != len(chainID) {
		return nil, errors.New("failed to write chainID")
	}

	return buf.Bytes(), nil
}

// GetCheckpointMsg returns the checkpoint message from the non-rp vote extension.
func GetCheckpointMsg(nonRpVoteExt []byte) (*checkpointTypes.MsgCheckpoint, error) {
	// Skip leading marker byte
	checkpointMsg, err := checkpointTypes.UnpackCheckpointSideSignBytes(nonRpVoteExt[1:])
	if err != nil {
		return nil, err
	}

	return checkpointMsg, nil
}
