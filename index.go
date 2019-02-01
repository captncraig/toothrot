package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

type GameIndex struct {
	RoomNames  RNAM
	MaxValues  MAXS
	RoomDir    OffsetMap // DROO
	ScriptDir  OffsetMap // DSCR
	SoundDir   OffsetMap // DSOU
	CostumeDir OffsetMap // DCOS
	CharsetDir OffsetMap // DCHR
	ObjectDir  OffsetMap // DOBJ
}

type RNAM map[byte]string

type MAXS struct {
	BlockName uint32
	BlockSize uint32
	Variables uint16
	// unknown uint16
	BitVars          uint16
	LocalObjects     uint16
	Charsets         uint16
	Verbs            uint16
	Array            uint16
	InventoryObjects uint16
}

type OffsetMap map[byte]uint32

func ParseIndex(dat []byte) *GameIndex {
	gi := &GameIndex{}
	off := 0
	for off < len(dat) {
		blk, newOff := readBlock(dat[off:], off)
		switch blk.Type {
		case "RNAM":
			gi.RoomNames = ParseRNAM(blk.Content)
		case "DROO":
			fmt.Println("DROO")
			gi.RoomDir = ParseOffsets(blk.Content)
		case "DSCR":
			fmt.Println("DSCR")
			gi.RoomDir = ParseOffsets(blk.Content)
		default:
			log.Printf("WARNING: Unimplemented block type found in game index file: '%s'", blk.Type)
		}
		off = newOff
	}
	return gi
}

func ParseRNAM(dat []byte) RNAM {
	names := map[byte]string{}
	for len(dat) > 1 {
		rNum := dat[0]
		for i := 1; i <= 10; i++ {
			dat[i] ^= 0xff
		}
		name := string(dat[1:10])
		dat = dat[10:]
		names[rNum] = name
	}
	return names
}

func ParseOffsets(dat []byte) OffsetMap {
	fmt.Println(dat)
	offs := OffsetMap{}
	count := binary.LittleEndian.Uint16(dat)
	dat = dat[2:]
	nums := make([]byte, count)
	for i := range nums {
		nums[i] = dat[i]
	}
	dat = dat[count:]
	for _, n := range nums {
		offs[n] = binary.LittleEndian.Uint32(dat[:4])
		dat = dat[4:]
		fmt.Println(n, offs[n])
	}
	return offs
}

func ParseInlineOffsets(dat []byte) OffsetMap {
	offs := OffsetMap{}
	count := int(dat[0])
	dat = dat[1:]
	for i := 0; i < count; i++ {
		offs[dat[0]] = binary.LittleEndian.Uint32(dat[1:])
		dat = dat[5:]
	}
	return offs
}
