package parser

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	EQUIHASH_SIZE = 1344 // size of an Equihash solution in bytes
)

// A block header as defined in version 2018.0-beta-29 of the Zcash Protocol Spec.
type RawBlockHeader struct {
	// The block version number indicates which set of block validation rules
	// to follow. The current and only defined block version number for Zcash
	// is 4.
	Version int32

	// A SHA-256d hash in internal byte order of the previous block's header. This
	// ensures no previous block can be changed without also changing this block's
	// header.
	HashPrevBlock [32]byte

	// A SHA-256d hash in internal byte order. The merkle root is derived from
	// the hashes of all transactions included in this block, ensuring that
	// none of those transactions can be modified without modifying the header.
	HashMerkleRoot [32]byte

	// [Pre-Sapling] A reserved field which should be ignored.
	// [Sapling onward] The root LEBS2OSP_256(rt) of the Sapling note
	// commitment tree corresponding to the final Sapling treestate of this
	// block.
	HashFinalSaplingRoot [32]byte

	// The block time is a Unix epoch time (UTC) when the miner started hashing
	// the header (according to the miner).
	Time uint32

	// An encoded version of the target threshold this block's header hash must
	// be less than or equal to, in the same nBits format used by Bitcoin.
	NBits [4]byte

	// An arbitrary field that miners can change to modify the header hash in
	// order to produce a hash less than or equal to the target threshold.
	Nonce [32]byte

	// The size of an Equihash solution in bytes (always 1344).
	SolutionSize EquihashSize

	// The Equihash solution.
	Solution [EQUIHASH_SIZE]byte
}

// EquihashSize is a concrete instance of Bitcoin's CompactSize encoding.
type EquihashSize struct {
	_    byte   // always the byte value 253
	Size uint16 // always 1344
}

func ReadBlockHeader(r io.Reader) (*RawBlockHeader, error) {
	var blockHeader RawBlockHeader
	err := binary.Read(r, binary.LittleEndian, &blockHeader)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading block header")
	}
	return &blockHeader, nil
}