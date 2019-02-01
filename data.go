package main

import (
	"fmt"
	"log"
)

type LECF struct {
	LOFF  OffsetMap
	Disks []*LFLF
}

type LFLF struct{}

func ParseGameData(dat []byte) *LECF {

	blk, newOff := readBlock(dat, 0)
	if blk.Type != "LECF" || newOff != len(dat) {
		fmt.Println(blk.Type, newOff, len(dat))
		log.Fatal("FILE SHOULD HAVE 1 BLOCK LECF")
	}
	return ParseLECF(blk.Content)
}

func ParseLECF(dat []byte) *LECF {
	game := &LECF{}
	off := 0
	for off < len(dat) {
		blk, newOff := readBlock(dat[off:], off)
		switch blk.Type {
		case "LOFF":
			game.LOFF = ParseInlineOffsets(blk.Content)
		case "LFLF":
			game.Disks = append(game.Disks, ParseLFLF(blk.Content))
		default:
			log.Printf("WARNING: Unimplemented block type found in LECF: '%s'", blk.Type)
		}
		off = newOff
	}
	return game
}

func ParseLFLF(dat []byte) *LFLF {
	return &LFLF{}
}
