package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Give filename to open")
	}
	dat, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for i := range dat {
		dat[i] = dat[i] ^ 0x69
	}
	gi := ParseGameData(dat)
	fmt.Println(gi)
}

type Block struct {
	Type    string
	Content []byte
	Offset  int
}

func readBlock(dat []byte, offset int) (Block, int) {
	blk := Block{
		Type:   string(dat[:4]),
		Offset: offset,
	}
	dat = dat[4:]
	size := int(binary.BigEndian.Uint32(dat[:4])) - 8
	dat = dat[4:]
	blk.Content = dat[:size]
	return blk, offset + 8 + size
}
