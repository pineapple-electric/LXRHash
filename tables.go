package lxr

import (
	"fmt"
	"os"
	"io/ioutil"
)


// constants for building different sized lookup tables (ByteMap).  Right now, the lookup table is hard coded as
// a 1K table, but it can be far larger.
const(
	firstrand = int64(0x13ef13156da2756b)
	Mapsiz    = 0x1000
	MapMask   = Mapsiz - 1
)

// GenerateTable
// Build a table with a rather simplistic but with many passes, adequately randomly ordered bytes.
// We do some straight forward bitwise math to initialize and scramble our ByteMap.
func (w *LXRHash) GenerateTable(rounds int) {

	// Our own "random" generator that really is just used to shuffle values
	offset := firstrand
	var b int64
	rand := func(i int64) int64 {
		b = int64(w.ByteMap[(offset&i^b)&MapMask]) ^ b<<9 ^ b>>1
		offset = offset<<9 ^ offset>>1 ^ offset>>7 ^ i ^ int64(w.ByteMap[(b+i)&MapMask])
		return (b^offset) & MapMask
	}

	// Fill the ByteMap with bytes ranging from 0 to 255.  As long as Mapsize%256 == 0, this
	// looping and masking works just fine.
	for i := range w.ByteMap {
		w.ByteMap[i] = byte(i)
	}

	// Now what we want to do is just mix it all up.  Take every byte in the ByteMap list, and exchange it
	// for some other byte in the ByteMap list. Note that we do this over and over, mixing and more mixing
	// the ByteMap, but maintaining the ratio of each byte value in the ByteMap list.
	for loops := 0; loops < rounds; loops++ {
		fmt.Println("Pass ", loops)
		for i := range w.ByteMap {
			j := rand(int64(i))
			w.ByteMap[i], w.ByteMap[j] = w.ByteMap[j], w.ByteMap[i]
		}
	}
}

// WriteTable
// When playing around with the algorithm, it is nice to generate files and use them off the disk.  This
// allows the user to do that, and save the cost of regeneration between test runs.
func (w *LXRHash) WriteTable(filename string) {
	// Ah, the data file isn't good for us.  Delete it (if it exists)
	os.Remove(filename)

	// open output file
	fo, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// write a chunk
	if _, err := fo.Write(w.ByteMap[:]); err != nil {
		panic(err)
	}

}

// ReadTable
// When a lookup table is on the disk, this will allow one to read it.
func (w *LXRHash) ReadTable(filename string) {

	// Try and load our byte map.
	dat, err := ioutil.ReadFile(filename)

	// If loading fails, or it is the wrong size, generate it.  Otherwise just use it.
	if err != nil || len(dat) != Mapsiz {
		w.GenerateTable(200)
		w.WriteTable(filename)
	} else {
		copy(w.ByteMap[:Mapsiz], dat)
	}

}

// Init()
// We use our own algorithm for initializing the map struct.  This is an fairly large table of
// byte values we use to map bytes to other byte values to enhance the avalanche nature of the hash
// as well as increase the memory footprint of the hash.
func (w *LXRHash) Init() {
	byteMap := []byte {
		0x30, 0x52, 0xcb, 0xe0, 0xdd, 0x03, 0x89, 0xea, 0xe9, 0xf0, 0xff, 0x73, 0x40, 0xf3, 0x7e, 0x1f,
		0x12, 0x23, 0x3f, 0xb7, 0x13, 0xd4, 0x3d, 0xf1, 0xea, 0xda, 0xc7, 0x0e, 0xfd, 0xaf, 0x90, 0x0a,
		0xb2, 0xcb, 0xab, 0xb0, 0x98, 0x59, 0x53, 0x51, 0xf3, 0x39, 0x6c, 0xa7, 0x43, 0xb1, 0x28, 0x16,
		0xfa, 0x0b, 0xe7, 0x57, 0x64, 0xb4, 0x4f, 0xbe, 0xb4, 0x9f, 0xad, 0x0c, 0x8a, 0xdc, 0x91, 0x18,
		0x56, 0xa5, 0x5f, 0x70, 0x36, 0xe3, 0x12, 0x31, 0xc3, 0x97, 0xfc, 0xdd, 0x2d, 0x64, 0x23, 0x12,
		0x61, 0xe0, 0xd2, 0x8e, 0x1c, 0x50, 0x4e, 0x12, 0xe7, 0x69, 0xb6, 0x4e, 0xa6, 0x38, 0xe4, 0xd7,
		0x11, 0x42, 0x37, 0x40, 0x5b, 0xc7, 0x6f, 0xee, 0x6c, 0xc3, 0x6d, 0xc3, 0xc3, 0x12, 0x8f, 0xe8,
		0x0b, 0x4b, 0x2e, 0x1e, 0x39, 0x7f, 0x7c, 0x35, 0x65, 0x3f, 0xa8, 0xcd, 0x37, 0xc5, 0x94, 0x10,
		0x64, 0x1f, 0x04, 0x83, 0x27, 0xa4, 0x0a, 0x25, 0x9d, 0xbf, 0xc8, 0xad, 0x5a, 0x51, 0xeb, 0xb7,
		0xc9, 0xe5, 0x3c, 0x78, 0x34, 0x85, 0x7b, 0x48, 0x23, 0xa3, 0x0a, 0x30, 0x7a, 0x9d, 0x70, 0x4b,
		0x7c, 0xae, 0x50, 0x6f, 0xae, 0x47, 0x2e, 0xe9, 0xd1, 0xe8, 0x17, 0xee, 0xb3, 0x08, 0x0d, 0xa2,
		0x04, 0x39, 0x68, 0x89, 0x8b, 0x06, 0xbb, 0x86, 0x3d, 0xbd, 0x74, 0xa4, 0xc4, 0x53, 0x2e, 0xd6,
		0xb5, 0x28, 0x56, 0x5f, 0x50, 0xba, 0xbb, 0xd7, 0xc4, 0x51, 0xd4, 0x3d, 0x27, 0x5c, 0x4d, 0xff,
		0x96, 0xf2, 0x0d, 0xaa, 0x75, 0xa4, 0xf8, 0x56, 0xd7, 0xe8, 0x82, 0xd3, 0xaa, 0x0e, 0x65, 0x84,
		0x38, 0xd1, 0xc7, 0xdc, 0xd5, 0x44, 0x11, 0x15, 0x3b, 0xf9, 0xf7, 0x0a, 0xec, 0x5d, 0x2d, 0xea,
		0xa9, 0x02, 0xa6, 0x8c, 0x32, 0x22, 0xc1, 0xfc, 0x86, 0x36, 0xfc, 0x06, 0x1b, 0xee, 0x06, 0xfd,
		0x97, 0xc6, 0x8c, 0x61, 0x80, 0x45, 0x9f, 0xb2, 0x8d, 0xd2, 0xed, 0xd7, 0xb5, 0x25, 0x07, 0xfe,
		0x4d, 0xe3, 0x15, 0x7a, 0x4e, 0xb1, 0xf0, 0xd0, 0xcc, 0x66, 0xae, 0x7d, 0x00, 0x0f, 0x77, 0x64,
		0x58, 0xce, 0x66, 0xb8, 0xdf, 0xa6, 0x80, 0xd3, 0x22, 0x96, 0x8e, 0x58, 0x79, 0x19, 0x44, 0x58,
		0xd2, 0xb6, 0x0f, 0x77, 0xd2, 0x75, 0xd3, 0x73, 0x2a, 0x79, 0x82, 0xdc, 0x5c, 0xf1, 0x84, 0x18,
		0x63, 0x92, 0x7f, 0xd6, 0x63, 0x98, 0x5f, 0x51, 0xb7, 0x39, 0x97, 0xaa, 0x30, 0x19, 0x9b, 0x88,
		0x04, 0xa6, 0x3c, 0xe0, 0xb8, 0x23, 0x6e, 0xc0, 0xa1, 0x13, 0xcd, 0xab, 0xda, 0xa0, 0xdd, 0x06,
		0x17, 0xba, 0x26, 0xc5, 0xe8, 0x54, 0xa9, 0xe1, 0x7e, 0x23, 0xb2, 0xbc, 0xa9, 0x03, 0x88, 0x8f,
		0x63, 0x42, 0x92, 0x13, 0xbc, 0x25, 0xe0, 0x77, 0xcb, 0x30, 0xa5, 0x76, 0x5a, 0xa3, 0xe9, 0x8f,
		0xc7, 0x6c, 0xfb, 0x9d, 0x1a, 0x52, 0x44, 0x9c, 0x5b, 0xe3, 0xb1, 0x6e, 0x5a, 0x20, 0x54, 0x49,
		0xed, 0x3a, 0x09, 0xa1, 0x15, 0xe2, 0x15, 0x6c, 0x21, 0x21, 0x8a, 0xb3, 0x48, 0xd6, 0x58, 0xd9,
		0x7e, 0xde, 0x30, 0x46, 0x70, 0x72, 0xd9, 0xde, 0xa8, 0x52, 0x56, 0xee, 0xe2, 0xfb, 0x59, 0x6d,
		0x00, 0xff, 0x42, 0x53, 0xde, 0xcf, 0x98, 0xec, 0x87, 0x53, 0x66, 0x47, 0x10, 0x5b, 0xcf, 0xb3,
		0xa5, 0x10, 0x0f, 0x9c, 0x05, 0x18, 0x7a, 0xdb, 0x3e, 0xed, 0xae, 0xde, 0x43, 0xa7, 0x33, 0x1f,
		0x52, 0xbd, 0x5d, 0xf1, 0xa9, 0x1d, 0xb6, 0xd2, 0x43, 0xd8, 0x4e, 0xac, 0xd9, 0xc2, 0xe3, 0x68,
		0xec, 0x49, 0x54, 0x99, 0x56, 0x19, 0x3a, 0x49, 0x36, 0x3c, 0x00, 0xb3, 0x48, 0xeb, 0xf2, 0x33,
		0xc7, 0x9b, 0x19, 0xf5, 0xbb, 0xee, 0x52, 0x39, 0xb4, 0x70, 0x07, 0x0a, 0x7c, 0x35, 0x93, 0x85,
		0xc7, 0x6a, 0xaa, 0xad, 0x49, 0xa6, 0x32, 0x08, 0x88, 0xf1, 0x7d, 0x4c, 0x88, 0x43, 0x0b, 0x3f,
		0x7f, 0x63, 0x60, 0x09, 0x0c, 0x00, 0x9d, 0xfa, 0x74, 0x99, 0xae, 0x58, 0x94, 0xa6, 0xb9, 0x07,
		0x86, 0xb6, 0x9e, 0x1c, 0x0b, 0xba, 0xd0, 0x58, 0xa0, 0xf0, 0xcf, 0xff, 0xd5, 0x74, 0xb6, 0x9f,
		0x3d, 0x98, 0x23, 0xd2, 0x2f, 0x26, 0x41, 0x65, 0x00, 0x9c, 0x74, 0x33, 0x7c, 0x9f, 0x4a, 0xab,
		0x63, 0x79, 0x13, 0x04, 0xeb, 0x88, 0xc9, 0xd5, 0x72, 0x3c, 0x0b, 0x77, 0x2d, 0x62, 0x4c, 0xb7,
		0x26, 0x9c, 0x4c, 0xf0, 0x5b, 0x7c, 0x24, 0xf6, 0xf0, 0xe7, 0x51, 0xad, 0xfe, 0xe2, 0x2e, 0xe3,
		0x81, 0xb3, 0x5d, 0xec, 0xa1, 0x7b, 0xf3, 0x98, 0x2e, 0x5f, 0xb5, 0xbc, 0x0d, 0x73, 0x71, 0xeb,
		0x4b, 0x08, 0x1e, 0x4c, 0xd5, 0x00, 0xdb, 0x69, 0x4a, 0x46, 0xd1, 0x26, 0xeb, 0x5d, 0xa8, 0x0e,
		0xf6, 0xbb, 0xbf, 0x69, 0xa8, 0xe9, 0x6d, 0xaa, 0x82, 0x2e, 0x1e, 0x76, 0x0f, 0xf3, 0x34, 0xfb,
		0x9b, 0xcb, 0xde, 0x91, 0x80, 0x7a, 0xf9, 0x29, 0x07, 0x83, 0xc6, 0x73, 0x00, 0x88, 0x29, 0x5f,
		0xc0, 0x38, 0x24, 0xc8, 0x50, 0x0a, 0xfd, 0x73, 0x8b, 0x37, 0xa6, 0xdb, 0xd5, 0x6b, 0x57, 0x05,
		0x14, 0x3d, 0x7a, 0xaf, 0xc2, 0x62, 0x50, 0xe0, 0x82, 0x2c, 0x65, 0xdc, 0xf9, 0x1e, 0xc5, 0x67,
		0xfb, 0x22, 0xc7, 0xd4, 0x45, 0x8c, 0x88, 0xa0, 0x51, 0x2f, 0xe5, 0xfc, 0x74, 0x1f, 0x34, 0x2e,
		0x73, 0x55, 0xde, 0xa3, 0x09, 0xd1, 0xc5, 0xc0, 0xa5, 0xff, 0xff, 0x9b, 0x18, 0x55, 0x8d, 0xd1,
		0xe9, 0xfc, 0xd6, 0x69, 0x46, 0x99, 0xdf, 0x21, 0x0c, 0x9c, 0xfc, 0xb4, 0x28, 0x16, 0xd3, 0x9a,
		0x6f, 0x84, 0x41, 0x4a, 0xef, 0x39, 0x6b, 0x97, 0x2a, 0x08, 0x39, 0x72, 0x38, 0x96, 0x88, 0xae,
		0xc2, 0x68, 0xcd, 0x6e, 0x51, 0x33, 0x63, 0x66, 0x45, 0xb5, 0x4b, 0x2a, 0x9a, 0xac, 0xb4, 0x92,
		0x10, 0x39, 0x8c, 0xc9, 0x70, 0x2f, 0xfc, 0xfd, 0x09, 0x05, 0x9c, 0xbf, 0x69, 0x48, 0xa9, 0x86,
		0xae, 0x21, 0x17, 0x68, 0x83, 0x88, 0x83, 0xbd, 0x7e, 0xcf, 0xa7, 0xf7, 0x88, 0xec, 0x99, 0x63,
		0xe7, 0xcf, 0xd6, 0x33, 0x90, 0xd3, 0x77, 0x40, 0xf1, 0x85, 0xaf, 0x27, 0xa4, 0x1d, 0x9f, 0x78,
		0xe9, 0xc0, 0xb3, 0x30, 0x98, 0x68, 0x26, 0x0d, 0x45, 0xc1, 0x7b, 0x9b, 0x43, 0xa1, 0x44, 0xb3,
		0xb3, 0x7f, 0xdb, 0x76, 0x23, 0xbf, 0xa5, 0xf7, 0x21, 0xa9, 0xd7, 0x6f, 0x42, 0xee, 0xd5, 0x20,
		0xb0, 0x7e, 0xe9, 0x56, 0xf7, 0x37, 0x5c, 0x3c, 0x67, 0xe6, 0x3c, 0x7f, 0x3b, 0xbf, 0x15, 0xef,
		0xd8, 0x40, 0x6f, 0x15, 0x72, 0xf2, 0xe1, 0xb0, 0x98, 0xfe, 0x9e, 0x8e, 0x8d, 0xe0, 0xf2, 0x23,
		0x87, 0xa2, 0xb8, 0x9e, 0x71, 0xdd, 0x14, 0xa9, 0x7a, 0x67, 0x7f, 0xfd, 0x9a, 0x47, 0x39, 0x6a,
		0x59, 0xbc, 0xb8, 0xc4, 0xc4, 0x4f, 0x02, 0xd1, 0x28, 0x31, 0x03, 0xf5, 0xf8, 0xcc, 0x4f, 0x43,
		0x49, 0x60, 0x08, 0xb2, 0xdf, 0x29, 0x11, 0x31, 0xb4, 0x9a, 0x94, 0xec, 0x7f, 0x90, 0x14, 0xd6,
		0x13, 0xdb, 0xdf, 0x22, 0xd7, 0x9b, 0x6d, 0x9e, 0xc5, 0xb1, 0x10, 0xb5, 0x17, 0xfb, 0x6e, 0x78,
		0x13, 0x39, 0x66, 0x6f, 0xb9, 0xa2, 0xf4, 0xc4, 0xda, 0xbe, 0x62, 0x30, 0x93, 0x87, 0xb6, 0xda,
		0xe4, 0x3a, 0x97, 0x3b, 0xa3, 0x1e, 0x8a, 0xc5, 0xe3, 0x50, 0x50, 0x53, 0x62, 0x26, 0x97, 0x62,
		0x7a, 0xa8, 0xee, 0x65, 0x8c, 0xce, 0x0d, 0xe9, 0x09, 0x2c, 0x02, 0x26, 0x5e, 0xd1, 0x84, 0xc9,
		0x8d, 0xbf, 0x8f, 0x85, 0xde, 0x19, 0x1e, 0x01, 0x4f, 0xa3, 0xea, 0x3b, 0x37, 0xf8, 0x01, 0xfb,
		0x0a, 0x9c, 0xc2, 0x47, 0xac, 0xe4, 0x83, 0xca, 0x53, 0x68, 0xc7, 0x9c, 0x6b, 0xab, 0xe1, 0x59,
		0xea, 0x89, 0x37, 0x7e, 0x08, 0x49, 0x5a, 0xca, 0xbc, 0xa6, 0x16, 0x66, 0x33, 0xd2, 0x9e, 0x8b,
		0x33, 0x77, 0x47, 0x6f, 0x75, 0xa7, 0xe0, 0xaf, 0xbd, 0x1a, 0x2d, 0x3a, 0x81, 0x60, 0x41, 0xf3,
		0x19, 0x16, 0x52, 0xa4, 0x6f, 0xdf, 0x82, 0xd8, 0x59, 0x1a, 0x5c, 0xca, 0xf8, 0x5f, 0xeb, 0xee,
		0x04, 0xd9, 0xb6, 0xfb, 0xe3, 0x7b, 0xcd, 0x4f, 0x61, 0xc2, 0x69, 0xba, 0x4a, 0xed, 0x92, 0x8c,
		0x0e, 0xc6, 0x44, 0x0d, 0x78, 0xa8, 0x52, 0x7d, 0x2b, 0xdd, 0x5e, 0x67, 0x34, 0xb3, 0x9f, 0x90,
		0x1d, 0x88, 0x85, 0xb4, 0x35, 0x31, 0xc0, 0xd8, 0x3c, 0x98, 0xe1, 0x35, 0x16, 0x5c, 0xc6, 0x32,
		0xb5, 0x4a, 0x25, 0x25, 0x41, 0xf4, 0x6c, 0x18, 0x83, 0x4b, 0x33, 0x38, 0xf1, 0x9f, 0x19, 0x41,
		0x12, 0x06, 0x86, 0x08, 0xeb, 0x61, 0xcd, 0xb8, 0x5d, 0xdb, 0x9a, 0x0f, 0x07, 0xb2, 0x91, 0x3b,
		0xf9, 0xac, 0x74, 0x45, 0x3f, 0x3d, 0x14, 0x52, 0x7e, 0x6a, 0x4a, 0xac, 0x95, 0x57, 0x48, 0x34,
		0x1f, 0xd9, 0x7e, 0x6c, 0xe8, 0x83, 0xcd, 0x35, 0x89, 0x28, 0xc1, 0x40, 0x33, 0x5f, 0x9b, 0x2c,
		0x6b, 0x10, 0xd4, 0xc6, 0xec, 0x0b, 0x6e, 0x04, 0x87, 0x5c, 0xb2, 0x19, 0x05, 0xfc, 0x2b, 0xe8,
		0x92, 0x5c, 0x52, 0x76, 0x9d, 0x18, 0xd3, 0xdf, 0xfe, 0xf9, 0x1a, 0x78, 0x6d, 0x61, 0x11, 0x95,
		0xa6, 0xbb, 0xfc, 0x33, 0x20, 0xf3, 0xa2, 0x77, 0xe6, 0x9d, 0xa3, 0x0d, 0xae, 0x67, 0x02, 0xc6,
		0x6e, 0x36, 0x65, 0x7b, 0xdb, 0x32, 0xa1, 0x23, 0x39, 0x49, 0xa5, 0xf9, 0xc0, 0x9b, 0x42, 0xb6,
		0xd2, 0x46, 0x7f, 0xed, 0xd4, 0x7d, 0x65, 0xf7, 0xc4, 0x96, 0x4f, 0x25, 0x0b, 0x49, 0x1f, 0x3c,
		0xb1, 0xe5, 0xbe, 0x18, 0xdd, 0xee, 0x45, 0x42, 0x30, 0xb9, 0x06, 0x6c, 0x10, 0xd3, 0x5b, 0x37,
		0x1e, 0xd1, 0xb7, 0x0c, 0xda, 0x8d, 0x5b, 0xc1, 0x17, 0xb0, 0x88, 0xbb, 0x23, 0x09, 0xdb, 0x95,
		0x8b, 0xc8, 0x28, 0xb8, 0x0d, 0xc6, 0xef, 0xaf, 0xf8, 0x34, 0x46, 0x4e, 0xee, 0x83, 0xaf, 0x18,
		0x52, 0x01, 0xe4, 0xc6, 0xd7, 0xa9, 0xc2, 0x45, 0x86, 0xf6, 0x8d, 0xb5, 0x39, 0x44, 0xec, 0x21,
		0xdf, 0x5a, 0xc5, 0x8b, 0x1c, 0xb9, 0xba, 0x98, 0xad, 0xa8, 0x89, 0x62, 0xca, 0x8a, 0xfe, 0x85,
		0xc0, 0x29, 0x10, 0x85, 0xc9, 0x05, 0x8b, 0xab, 0x4c, 0x4a, 0x93, 0x81, 0xee, 0x91, 0x06, 0xe8,
		0x3a, 0x42, 0xf8, 0x7b, 0x2a, 0x4b, 0x36, 0x3f, 0x40, 0x1e, 0xec, 0x91, 0x00, 0x22, 0xd6, 0xef,
		0x7d, 0xb8, 0x9e, 0xd4, 0x96, 0xf5, 0xde, 0x47, 0x6c, 0x4b, 0x1a, 0xcf, 0xea, 0x1e, 0x32, 0xcb,
		0x80, 0xd4, 0xbb, 0x9a, 0x84, 0x15, 0x48, 0x95, 0xc5, 0xb5, 0x79, 0x90, 0xab, 0x86, 0x05, 0x50,
		0x28, 0xa0, 0x22, 0x05, 0x33, 0x6b, 0xb1, 0xc3, 0x58, 0xe2, 0x9d, 0x6b, 0x8f, 0x67, 0x36, 0x5b,
		0x14, 0x69, 0xad, 0x25, 0x89, 0x22, 0x01, 0x31, 0xaa, 0x8e, 0xcb, 0x79, 0x5b, 0xff, 0xa6, 0x51,
		0xa1, 0xd8, 0xf1, 0xd0, 0x18, 0xb7, 0xa7, 0x59, 0x40, 0x72, 0x4d, 0x55, 0x7e, 0x91, 0x28, 0x5f,
		0x72, 0xbd, 0x38, 0x1a, 0xc2, 0x2f, 0xfa, 0x07, 0x3d, 0xd9, 0xfa, 0xc6, 0xa6, 0xd6, 0x46, 0xb8,
		0x0e, 0xaf, 0x3b, 0xe6, 0x28, 0xa0, 0x02, 0xe2, 0x30, 0xce, 0xf6, 0xe7, 0xde, 0x59, 0xb1, 0x56,
		0xbc, 0xc1, 0x55, 0x3f, 0x3d, 0x72, 0x1c, 0x6a, 0xb0, 0xf9, 0x70, 0x47, 0xde, 0xdf, 0xe7, 0x2e,
		0xfc, 0x24, 0xac, 0x7c, 0x99, 0xeb, 0x8b, 0x8a, 0x7e, 0x47, 0x35, 0x0c, 0x15, 0x14, 0x29, 0x5e,
		0x13, 0x4a, 0x0e, 0xe6, 0xef, 0xaf, 0x18, 0x28, 0xb2, 0xcf, 0xd4, 0x71, 0x23, 0xcb, 0x84, 0xe4,
		0xfc, 0x80, 0xa1, 0x20, 0x93, 0xda, 0x87, 0x02, 0x99, 0xb7, 0xf2, 0x2d, 0x82, 0xce, 0xa9, 0x9e,
		0x78, 0x70, 0x82, 0xf3, 0x72, 0xe6, 0x90, 0x0a, 0x0b, 0xbb, 0x54, 0x97, 0x6c, 0xb4, 0xd3, 0xa0,
		0x0d, 0xd1, 0xd6, 0x2a, 0x6d, 0x3a, 0x0d, 0x92, 0xc8, 0x43, 0x6b, 0xe5, 0xa5, 0x54, 0x94, 0x1d,
		0x5e, 0x20, 0xac, 0x7b, 0xd3, 0xf6, 0x06, 0xd0, 0xad, 0x8c, 0x8f, 0xfc, 0x91, 0x1a, 0xa2, 0x00,
		0x62, 0xd7, 0x0a, 0x37, 0x8d, 0x3b, 0x2a, 0xa1, 0xed, 0xe8, 0x8d, 0x9f, 0x09, 0xf7, 0x41, 0x49,
		0xe3, 0x5f, 0x43, 0xfa, 0x36, 0xf8, 0xf6, 0x96, 0x11, 0x18, 0xef, 0x2a, 0xbd, 0x42, 0x82, 0xa4,
		0x5f, 0x7c, 0x40, 0x82, 0x03, 0xf8, 0x69, 0x87, 0x28, 0x30, 0x49, 0xcb, 0xe8, 0xfa, 0x83, 0x2b,
		0x53, 0x07, 0xa8, 0x6f, 0x04, 0xbd, 0x2b, 0x8c, 0xaf, 0x4e, 0x54, 0x71, 0xab, 0x8b, 0x6d, 0x31,
		0x49, 0xea, 0xff, 0x19, 0x00, 0x78, 0xac, 0x09, 0x2a, 0x46, 0xce, 0x22, 0xef, 0x9f, 0x75, 0xbe,
		0x8d, 0xaa, 0x87, 0x5a, 0x1b, 0x4c, 0x37, 0xbf, 0xb7, 0xa3, 0x11, 0x4a, 0x13, 0xfa, 0xe5, 0xe5,
		0x25, 0x31, 0x69, 0x2d, 0x76, 0x7a, 0xfa, 0xe2, 0xbe, 0x71, 0x1a, 0x78, 0x2f, 0x94, 0x3d, 0x37,
		0x92, 0xeb, 0x29, 0x03, 0x79, 0x56, 0xe7, 0x96, 0xb9, 0x2b, 0xb9, 0xf3, 0xfd, 0xba, 0x94, 0xe6,
		0x0e, 0x83, 0x50, 0xee, 0xbd, 0x4a, 0x37, 0xba, 0x15, 0x58, 0x17, 0x02, 0xa3, 0x4e, 0xc1, 0x8f,
		0x94, 0x14, 0x46, 0x1c, 0x41, 0x49, 0x4b, 0x24, 0x0d, 0xf2, 0xbe, 0x96, 0xcc, 0x95, 0xc2, 0x97,
		0x65, 0x8e, 0xf7, 0x3e, 0x19, 0x91, 0x35, 0x8a, 0x24, 0x37, 0x1c, 0x02, 0x60, 0xfc, 0x81, 0xaf,
		0x64, 0x32, 0x1c, 0x27, 0xab, 0x6e, 0xa2, 0xe4, 0x80, 0xe6, 0xcf, 0xb6, 0x4d, 0x06, 0xb4, 0xe3,
		0xba, 0x7b, 0x3a, 0x77, 0xd3, 0xdf, 0xd2, 0x78, 0x98, 0xe1, 0x42, 0x5a, 0xe4, 0xd9, 0x34, 0x61,
		0x93, 0xeb, 0x95, 0xfd, 0x7d, 0x62, 0x08, 0x0f, 0x57, 0x6d, 0x13, 0xd0, 0x85, 0x09, 0x1f, 0xb0,
		0xce, 0x23, 0xf2, 0xc7, 0x92, 0x5e, 0x1b, 0x6a, 0x90, 0x8f, 0xf9, 0xf7, 0x76, 0x46, 0x7d, 0x04,
		0x6e, 0x92, 0x3e, 0x8b, 0xf5, 0x12, 0x34, 0x5d, 0xa2, 0x09, 0xae, 0x95, 0x7c, 0xf6, 0xe6, 0x45,
		0x00, 0x53, 0x5a, 0x25, 0x3d, 0x2b, 0xc6, 0x8c, 0x9d, 0xc9, 0xe9, 0xf6, 0xb9, 0x30, 0x87, 0x20,
		0x79, 0x2b, 0x38, 0xcb, 0xe7, 0xf6, 0x66, 0xc8, 0x69, 0x32, 0xed, 0x72, 0xf6, 0x3b, 0x3b, 0x7b,
		0x0f, 0x75, 0x9f, 0x93, 0x5f, 0xc2, 0x0e, 0xb8, 0xc6, 0x1b, 0x84, 0xaa, 0x26, 0x00, 0x0e, 0xa2,
		0xf5, 0xc9, 0x06, 0x8b, 0x20, 0x7c, 0xf5, 0xe8, 0x9f, 0x6d, 0x2c, 0x3b, 0x1b, 0x39, 0xa7, 0x88,
		0xd5, 0xa7, 0x65, 0x6b, 0x11, 0x1f, 0x2a, 0x03, 0xd4, 0x0b, 0xf7, 0x0f, 0xf8, 0x8d, 0x15, 0x8f,
		0x14, 0x60, 0x57, 0x97, 0x07, 0xa9, 0xda, 0xb0, 0x34, 0x0c, 0x24, 0xd1, 0x98, 0x9c, 0x0c, 0x3f,
		0x6b, 0x16, 0x80, 0x13, 0xbe, 0x9e, 0x46, 0x87, 0xb0, 0x99, 0x86, 0xcb, 0xc5, 0x4e, 0x71, 0x0b,
		0x30, 0x74, 0xb5, 0xf8, 0x20, 0xa1, 0xdc, 0xdf, 0xa5, 0x00, 0xce, 0xd5, 0x33, 0xaa, 0xdf, 0xd6,
		0x4a, 0x29, 0xdc, 0xe1, 0xf9, 0x60, 0x40, 0x61, 0x50, 0x70, 0x26, 0xf6, 0x41, 0xd0, 0x55, 0xd8,
		0x36, 0x22, 0x10, 0x64, 0x18, 0xba, 0x25, 0xfb, 0xb1, 0x41, 0x2d, 0x81, 0xc1, 0xf3, 0x21, 0x52,
		0xe5, 0x41, 0xa3, 0x0a, 0x5b, 0xaa, 0xf0, 0x3e, 0xee, 0x9a, 0x9b, 0xc4, 0x24, 0x4a, 0x1e, 0xf9,
		0x19, 0xe8, 0x2f, 0xa4, 0xd1, 0xb9, 0x05, 0x78, 0x25, 0xc2, 0x7c, 0x80, 0x79, 0x4b, 0xec, 0x15,
		0x7f, 0x1b, 0xc2, 0xe0, 0xbf, 0x36, 0xb3, 0x35, 0xbe, 0x31, 0x5c, 0x97, 0x4c, 0xad, 0x64, 0x22,
		0x64, 0x4f, 0xb4, 0xf0, 0x21, 0x63, 0xd5, 0x2f, 0xe6, 0x87, 0xe5, 0xbd, 0xef, 0xfb, 0xec, 0xf5,
		0x32, 0x94, 0xf8, 0xae, 0x5e, 0xd4, 0x9e, 0x8f, 0xc3, 0x7c, 0x63, 0x31, 0xe4, 0xa5, 0x82, 0x5d,
		0xa0, 0xd2, 0xe6, 0xc8, 0xc4, 0x8a, 0x32, 0x85, 0xd0, 0xc1, 0xe8, 0x98, 0xf3, 0x76, 0xd9, 0xbc,
		0x5a, 0x93, 0xe4, 0x94, 0x9d, 0x79, 0xbc, 0xaa, 0x67, 0x60, 0xc1, 0x97, 0x3a, 0x81, 0x4e, 0xab,
		0xf4, 0x54, 0x72, 0xaa, 0xf2, 0x2c, 0x5f, 0x40, 0x20, 0xab, 0x59, 0x01, 0x40, 0x6e, 0x79, 0x21,
		0x11, 0x5e, 0xe7, 0x96, 0x8f, 0xa9, 0x13, 0xa5, 0x6e, 0x0a, 0x07, 0x62, 0x35, 0x57, 0xb8, 0x18,
		0x51, 0x34, 0x56, 0x9f, 0xe5, 0x1d, 0x8e, 0xd5, 0x26, 0xeb, 0xeb, 0x53, 0xaa, 0xc8, 0x31, 0x63,
		0x9a, 0xcc, 0xb1, 0x1b, 0x3a, 0xb1, 0xd6, 0xc3, 0x3d, 0x6d, 0x32, 0x74, 0xbb, 0x5c, 0x8b, 0xd6,
		0xb2, 0x8c, 0x57, 0x67, 0x36, 0xa8, 0x62, 0x43, 0xbd, 0x7a, 0x80, 0x29, 0xfd, 0x3c, 0x09, 0xf4,
		0x52, 0x17, 0x54, 0x9c, 0xf9, 0x1c, 0x57, 0x16, 0x5a, 0x32, 0x3e, 0xd5, 0x8b, 0xa8, 0x92, 0xf2,
		0x36, 0x8d, 0x1f, 0x8e, 0x2f, 0x19, 0x71, 0x4a, 0xe9, 0x1e, 0xe2, 0x75, 0x86, 0x3f, 0x22, 0x9a,
		0xca, 0xb6, 0x51, 0x9d, 0x97, 0x8f, 0x41, 0x44, 0x99, 0xfe, 0xd7, 0x94, 0x75, 0x73, 0x21, 0x09,
		0x02, 0xd1, 0x4f, 0x39, 0x98, 0xfb, 0xc9, 0x4d, 0xda, 0x3f, 0x4e, 0x5f, 0xe1, 0x43, 0x1d, 0xc3,
		0xf7, 0x54, 0x0a, 0x61, 0x2c, 0xda, 0x5c, 0x85, 0x75, 0x87, 0x9c, 0x5f, 0x65, 0x7a, 0x99, 0x10,
		0x20, 0xb5, 0x46, 0xc5, 0x8a, 0x5a, 0x81, 0x2a, 0x04, 0x68, 0xd4, 0x3b, 0x58, 0x44, 0xf8, 0x2c,
		0xa1, 0xb3, 0xcf, 0x7e, 0x7b, 0x5e, 0xb9, 0x93, 0xed, 0xdb, 0xed, 0x1d, 0x9c, 0xed, 0x14, 0x1f,
		0x8c, 0xf0, 0x0e, 0x4d, 0x84, 0x96, 0xf4, 0x7e, 0xd3, 0x03, 0x0c, 0x84, 0x77, 0xbf, 0xc7, 0xfb,
		0x9a, 0xcc, 0x34, 0x12, 0x50, 0x8c, 0x6b, 0x66, 0xd0, 0xd9, 0xbc, 0xcc, 0x08, 0xca, 0x51, 0x68,
		0x0b, 0x8b, 0x1c, 0xfb, 0xd3, 0x84, 0xa0, 0xcd, 0xfc, 0x76, 0xb4, 0xba, 0xdd, 0x46, 0x6f, 0x2d,
		0x28, 0x1d, 0xef, 0x83, 0xa2, 0x16, 0xcd, 0xfe, 0xe1, 0x63, 0x25, 0x2f, 0xb9, 0x4b, 0x82, 0x48,
		0x76, 0xc2, 0xba, 0xfd, 0x80, 0x4e, 0xde, 0x7c, 0x5e, 0x8e, 0xf0, 0xa7, 0xc1, 0xf6, 0x96, 0xf2,
		0x03, 0x09, 0x24, 0x46, 0xa3, 0x63, 0x7f, 0x1b, 0x56, 0x2b, 0xca, 0xb3, 0x3e, 0x59, 0x0a, 0x80,
		0xaa, 0x8b, 0x62, 0xc1, 0x31, 0xd7, 0xa4, 0x2e, 0x8a, 0x01, 0xd4, 0x3e, 0x2d, 0x99, 0x45, 0x81,
		0xdf, 0x66, 0x67, 0xe7, 0x40, 0xb1, 0xf3, 0xa3, 0xa2, 0x71, 0x4e, 0x6a, 0xdf, 0xf4, 0x87, 0xf2,
		0xf4, 0x57, 0xfa, 0x77, 0x49, 0xc0, 0x97, 0x85, 0x43, 0xd7, 0x45, 0x26, 0x00, 0xbe, 0x85, 0x23,
		0x7d, 0xc8, 0x6f, 0x20, 0xb5, 0xaf, 0x96, 0xe3, 0x7e, 0xbc, 0x6d, 0x34, 0xda, 0x16, 0x9d, 0x89,
		0xb1, 0x82, 0x5c, 0xad, 0x89, 0x29, 0x50, 0xcb, 0xc2, 0x03, 0x1a, 0x7e, 0xb4, 0x61, 0x1c, 0x34,
		0x1f, 0xbd, 0xb7, 0x08, 0x31, 0x9b, 0xd5, 0xa1, 0xa8, 0x8a, 0x05, 0x11, 0xc6, 0xc8, 0xdb, 0x1b,
		0xa5, 0xa4, 0xef, 0x5d, 0xfd, 0x81, 0x9f, 0xe1, 0x42, 0xac, 0xe6, 0x3c, 0x50, 0x08, 0xb7, 0x15,
		0xe9, 0x91, 0x49, 0x74, 0x48, 0x05, 0xc8, 0x70, 0x8d, 0x6c, 0x3e, 0x47, 0xac, 0xa0, 0x6c, 0x13,
		0x93, 0x53, 0x9b, 0x74, 0x27, 0xe0, 0xcb, 0xbb, 0x1f, 0xa6, 0xe2, 0xc8, 0x4b, 0xba, 0xb3, 0xcf,
		0x16, 0xb0, 0x07, 0x42, 0x83, 0xa0, 0xff, 0x58, 0x3f, 0x93, 0x44, 0x37, 0xaf, 0xff, 0x38, 0x95,
		0xc3, 0x9d, 0x20, 0x0f, 0x3a, 0xe6, 0xfd, 0xa5, 0x35, 0x62, 0xc6, 0xb0, 0x8a, 0xa0, 0x58, 0xa9,
		0xbb, 0xf1, 0x4a, 0xe1, 0x5c, 0x66, 0x2c, 0xb6, 0x69, 0xfa, 0xe4, 0x4b, 0xad, 0xec, 0x6a, 0x75,
		0x0e, 0x29, 0x3c, 0xa3, 0x22, 0x70, 0x50, 0x11, 0x01, 0xf4, 0x81, 0x9b, 0xce, 0x8f, 0xa7, 0xc3,
		0x64, 0x37, 0x1e, 0x0a, 0xc5, 0x67, 0x78, 0x4d, 0x7c, 0x61, 0xdd, 0x06, 0xec, 0x76, 0xcd, 0x1e,
		0x0d, 0xae, 0x32, 0x14, 0xef, 0x92, 0xe7, 0x35, 0x1d, 0x8f, 0xdc, 0x25, 0xbb, 0x6e, 0xf9, 0xce,
		0x80, 0xa4, 0x11, 0xc1, 0x7f, 0x2d, 0x3c, 0xe8, 0xdf, 0xd8, 0x26, 0xbd, 0x55, 0x81, 0x61, 0x2c,
		0x3c, 0x6d, 0x68, 0x45, 0x75, 0x24, 0x9d, 0xf2, 0x7f, 0x2b, 0xbe, 0xbd, 0x64, 0x5e, 0xd1, 0x6a,
		0xf8, 0xdc, 0xb6, 0xa7, 0x64, 0xf4, 0x07, 0x89, 0x1e, 0xce, 0xa3, 0xac, 0x95, 0x2c, 0xa4, 0x12,
		0x0d, 0x02, 0xf1, 0x08, 0xed, 0xf4, 0xa4, 0xf4, 0xb2, 0xdd, 0x26, 0xe0, 0x1f, 0x23, 0x56, 0xb6,
		0xed, 0xd2, 0x3b, 0xe9, 0x01, 0x45, 0x93, 0x9a, 0xc0, 0xfd, 0xe6, 0x3f, 0xb1, 0xac, 0x6e, 0x82,
		0x37, 0x13, 0x80, 0x01, 0xd6, 0x57, 0x43, 0x24, 0xd9, 0x0f, 0x8e, 0x91, 0x95, 0xa8, 0xa7, 0xed,
		0x3b, 0x34, 0x1a, 0x1a, 0x01, 0x3e, 0x87, 0x7d, 0x48, 0x54, 0x7a, 0x89, 0x78, 0x04, 0x54, 0xd2,
		0x38, 0xd7, 0x5a, 0x11, 0x31, 0x72, 0xc6, 0x24, 0xbc, 0x03, 0x57, 0x61, 0xe7, 0x6e, 0x8e, 0x54,
		0xf4, 0x99, 0x75, 0xe2, 0xea, 0xa2, 0xb2, 0x79, 0x17, 0xc9, 0x41, 0xdb, 0x50, 0xb7, 0x9b, 0x22,
		0xb1, 0x77, 0xb5, 0xec, 0xc7, 0x5d, 0x2e, 0xba, 0xf3, 0x0b, 0x29, 0x52, 0x47, 0x11, 0xad, 0xfe,
		0x94, 0x41, 0x7a, 0x2b, 0x3e, 0x16, 0xd8, 0xcd, 0x7f, 0x74, 0x62, 0x78, 0x0d, 0x59, 0xb2, 0xcd,
		0xfe, 0x93, 0xa0, 0x90, 0x7d, 0x20, 0x27, 0xe2, 0x05, 0xfc, 0xdd, 0x1c, 0xcf, 0x98, 0x7f, 0xe3,
		0xb5, 0x81, 0x46, 0x68, 0x4f, 0x87, 0xf5, 0x28, 0x3a, 0x5c, 0x20, 0xf5, 0x30, 0x02, 0x41, 0x7d,
		0xea, 0xb8, 0x56, 0x3a, 0x12, 0xc4, 0xf2, 0x40, 0xe5, 0x86, 0x90, 0x83, 0x60, 0xb3, 0x2b, 0x36,
		0x8b, 0xb7, 0x53, 0xe5, 0xe8, 0xf7, 0xa9, 0xf2, 0xcc, 0x51, 0x87, 0x0b, 0xf9, 0xa1, 0xae, 0x13,
		0xff, 0x1e, 0xb0, 0x64, 0xdd, 0x1c, 0x6f, 0x51, 0x86, 0x35, 0x64, 0xda, 0xff, 0x5d, 0x5e, 0x14,
		0x59, 0x42, 0x5c, 0x26, 0x06, 0x4c, 0xb5, 0x73, 0xcc, 0x8a, 0x27, 0x15, 0x3e, 0x15, 0x94, 0x48,
		0x7e, 0x25, 0x9e, 0xfd, 0xa5, 0xd3, 0x1b, 0x36, 0x1a, 0xc0, 0x5e, 0x59, 0x53, 0xcc, 0x0c, 0x83,
		0x0f, 0x56, 0xda, 0xc3, 0x0c, 0xc4, 0x55, 0xcb, 0xf5, 0x04, 0x97, 0x12, 0xc1, 0x60, 0xf0, 0x0c,
		0xdc, 0xd2, 0x47, 0xff, 0xce, 0x1b, 0xb3, 0x18, 0xdc, 0x2e, 0xf0, 0x31, 0x71, 0xb9, 0x90, 0xee,
		0x55, 0x06, 0x77, 0x27, 0x1d, 0xab, 0x71, 0xa4, 0x86, 0xe5, 0xe0, 0x01, 0x88, 0x47, 0xa0, 0x45,
		0xb9, 0x55, 0x61, 0x3e, 0x0c, 0x9e, 0xb8, 0xdc, 0x16, 0xc6, 0xca, 0x30, 0xcf, 0x2f, 0xcc, 0x34,
		0xe8, 0x47, 0x32, 0xca, 0x8f, 0x10, 0xf6, 0x59, 0xbe, 0x1a, 0x31, 0x15, 0x9d, 0xbc, 0xca, 0x68,
		0x12, 0xbf, 0x95, 0x5e, 0xb8, 0xbb, 0xfd, 0x34, 0xd8, 0x99, 0x23, 0x73, 0x66, 0x90, 0xf1, 0x08,
		0xc0, 0xb2, 0x0f, 0x7d, 0x85, 0x8c, 0xe9, 0x82, 0x49, 0x95, 0x73, 0x30, 0x56, 0x3c, 0xe5, 0xaa,
		0x74, 0x6b, 0x19, 0x33, 0x2f, 0x38, 0x47, 0x76, 0x19, 0x5c, 0x60, 0x6d, 0xb4, 0xd8, 0x4b, 0xe2,
		0x2f, 0x3d, 0xd2, 0x89, 0xd9, 0x44, 0x10, 0xfb, 0x84, 0xa7, 0x27, 0x0f, 0x89, 0xe0, 0x92, 0xea,
		0xd7, 0xe5, 0xc7, 0xb0, 0x4c, 0xa4, 0x6f, 0xcf, 0x41, 0x39, 0xbc, 0xd4, 0x17, 0x51, 0xe0, 0x38,
		0x3d, 0x17, 0xa0, 0x63, 0x70, 0x99, 0x90, 0xa2, 0x10, 0x28, 0xbf, 0x60, 0x32, 0x83, 0x21, 0x68,
		0xfe, 0xd7, 0x21, 0xad, 0xb5, 0x1b, 0xee, 0xa9, 0x4e, 0xbe, 0x6b, 0xf9, 0xd2, 0x5b, 0x35, 0x67,
		0x24, 0xc7, 0x82, 0x8e, 0x44, 0x69, 0x01, 0x96, 0x79, 0x94, 0x3a, 0x48, 0x1b, 0x58, 0xdd, 0xc7,
		0x4d, 0x55, 0xbb, 0x60, 0xbe, 0x02, 0x79, 0x86, 0xc8, 0x1c, 0xa8, 0xc0, 0xd4, 0x6c, 0xd9, 0x08,
		0x8c, 0x70, 0x9c, 0x32, 0x91, 0x67, 0xf1, 0x99, 0xc1, 0x02, 0x5a, 0x9e, 0x3e, 0x05, 0xfa, 0x73,
		0xe2, 0xde, 0xcc, 0x69, 0x82, 0x4f, 0x65, 0x5d, 0x91, 0x79, 0x61, 0x1f, 0x8e, 0x60, 0xeb, 0x3a,
		0x7c, 0x78, 0x2e, 0x2f, 0x86, 0x89, 0x26, 0x89, 0xa7, 0x7b, 0x5c, 0x29, 0x6d, 0x3b, 0x77, 0xca,
		0xd8, 0x4c, 0x55, 0xd0, 0x4d, 0x91, 0xdd, 0x93, 0xae, 0x5f, 0x1d, 0xe9, 0x01, 0xb8, 0xa9, 0x9b,
		0xe7, 0x74, 0x68, 0x09, 0x20, 0xf9, 0x04, 0x7d, 0xd0, 0x5d, 0x5b, 0x0b, 0x16, 0x5b, 0x1a, 0xd0,
		0xa1, 0x69, 0x3e, 0xf7, 0x43, 0x5b, 0xfa, 0x16, 0x21, 0x02, 0x8d, 0xa4, 0x11, 0xd1, 0xbc, 0x38,
		0x63, 0x58, 0x66, 0xc5, 0xf4, 0xac, 0xd8, 0xa8, 0x18, 0x97, 0x16, 0x52, 0xe1, 0xd4, 0x75, 0xcc,
		0xd1, 0xcc, 0xc4, 0xee, 0x9f, 0xcb, 0x9b, 0xf3, 0x58, 0x7d, 0xa6, 0x60, 0x44, 0xde, 0xf6, 0x9a,
		0x99, 0xe4, 0x30, 0x36, 0x8c, 0xd8, 0x9e, 0xd8, 0x3e, 0xb2, 0xb9, 0x57, 0x9a, 0xdf, 0x52, 0xc8,
		0xed, 0x0e, 0x68, 0x0f, 0xcc, 0x6a, 0x24, 0xfe, 0x0f, 0x67, 0x6a, 0x90, 0x4c, 0x6c, 0xa0, 0x0b,
		0xa7, 0x96, 0xf1, 0x7a, 0x38, 0xd3, 0x83, 0xf4, 0xc9, 0xaf, 0xff, 0x6f, 0xed, 0xe9, 0xca, 0xd9,
		0x7c, 0xf7, 0x40, 0x5b, 0xea, 0x59, 0xef, 0xa6, 0x6b, 0x43, 0x4f, 0xb1, 0x76, 0x72, 0x69, 0x65,
		0x4f, 0x5a, 0x2c, 0x7b, 0xbf, 0x6e, 0x14, 0x71, 0x3b, 0xbf, 0xcf, 0x2d, 0x24, 0x10, 0x2c, 0x06,
		0x6a, 0xfb, 0x48, 0x67, 0x36, 0x56, 0xea, 0xe1, 0xb0, 0x28, 0x19, 0x5b, 0xeb, 0x41, 0x6a, 0x5e,
		0x94, 0x92, 0xb9, 0x45, 0xa5, 0x29, 0x2b, 0xe5, 0xa3, 0xbd, 0x3a, 0x8e, 0xd4, 0xe3, 0x2a, 0x2d,
		0x48, 0xa5, 0xb8, 0xf8, 0xa7, 0x53, 0xad, 0xe4, 0xc8, 0x76, 0x19, 0xc9, 0x73, 0xde, 0x4d, 0xe4,
		0xe2, 0xaf, 0xc0, 0x95, 0xc2, 0xb4, 0x66, 0xab, 0x65, 0x2d, 0xae, 0x08, 0x88, 0xcf, 0xb0, 0x7a,
		0x28, 0xb5, 0x0c, 0xc0, 0xac, 0xad, 0xd2, 0xc9, 0xd0, 0xe1, 0x1f, 0x38, 0x55, 0x38, 0x2b, 0x4d,
		0x9c, 0x2a, 0x72, 0x04, 0x9f, 0x4c, 0xfb, 0x05, 0x8e, 0x9a, 0x13, 0x36, 0x35, 0x6a, 0xe1, 0x04,
		0xf1, 0xe3, 0xae, 0xf5, 0xf3, 0x81, 0x74, 0x17, 0x7b, 0xcd, 0x6a, 0xab, 0x6c, 0x27, 0xea, 0x00,
		0xc4, 0x06, 0x46, 0x3c, 0x68, 0xc0, 0xf3, 0xd0, 0xe2, 0xdc, 0xfe, 0xf5, 0x81, 0x4c, 0xc5, 0xe8,
		0x17, 0x4b, 0x77, 0xb8, 0xef, 0x67, 0x5d, 0x05, 0xe1, 0x03, 0xc3, 0xad, 0x33, 0x48, 0x93, 0xdb,
		0x92, 0x78, 0xcb, 0xca, 0x71, 0x64, 0x01, 0xbb, 0xe4, 0x09, 0xf8, 0x85, 0xb7, 0x4c, 0xcc, 0x80,
		0xc3, 0xeb, 0x4f, 0xec, 0x2f, 0x4c, 0x4b, 0x75, 0xbe, 0x6d, 0x95, 0x47, 0xa2, 0xea, 0x55, 0x31,
		0x1c, 0x71, 0x44, 0x91, 0x1f, 0x4d, 0x71, 0x8f, 0xcd, 0x6a, 0x7a, 0x57, 0x74, 0x4a, 0xc7, 0x07,
		0x78, 0x81, 0x22, 0x4a, 0x1b, 0xe7, 0x22, 0x4d, 0xc0, 0xa5, 0x14, 0xc4, 0xb6, 0xd3, 0x79, 0x85,
		0x12, 0x53, 0x2f, 0x64, 0x6c, 0xc9, 0xca, 0x23, 0x07, 0x1d, 0x94, 0xc5, 0xfe, 0x21, 0x42, 0x94,
		0x4c, 0x32, 0x2d, 0xb0, 0xf2, 0xcd, 0x9a, 0xbe, 0xc4, 0x17, 0x6e, 0xaf, 0x0e, 0xd5, 0x8a, 0x71,
		0x1d, 0x62, 0xfd, 0x27, 0x88, 0x1c, 0xea, 0x2a, 0x49, 0x75, 0xe1, 0x2a, 0x7b, 0xd5, 0xe3, 0xc3,
		0xfa, 0xd0, 0x2c, 0xb7, 0xac, 0x51, 0x5a, 0xd1, 0x07, 0x75, 0xec, 0xd5, 0x5e, 0xcd, 0x5d, 0xfe,
		0xba, 0xce, 0xb9, 0xde, 0x60, 0x7b, 0xd9, 0x48, 0x9e, 0xab, 0x70, 0xb2, 0xfb, 0xb7, 0x9f, 0xf0,
		0x95, 0x17, 0x8b, 0x50, 0xf7, 0x2e, 0xf2, 0xce, 0xab, 0x96, 0xe6, 0xf3, 0xf0, 0xf0, 0xc6, 0xde,
		0xda, 0x06, 0xf5, 0x98, 0xfc, 0x57, 0x44, 0x84, 0xef, 0x01, 0xc8, 0x65, 0xbf, 0x80, 0x21, 0xe7,
		0x63, 0x2a, 0xea, 0xef, 0x6b, 0x2c, 0xc5, 0x72, 0x5d, 0x81, 0xf8, 0xf7, 0xa8, 0x86, 0x2b, 0xf1,
		0x2a, 0x66, 0xe7, 0x63, 0x14, 0xf1, 0xce, 0x3d, 0x54, 0x2e, 0x7e, 0x56, 0xed, 0xc1, 0x61, 0x85,
		0x3d, 0x35, 0x89, 0xa1, 0x60, 0x8d, 0x7b, 0xdd, 0x9a, 0x72, 0xd6, 0xbe, 0xa8, 0xb4, 0x8a, 0x3f,
		0x9e, 0x02, 0xc4, 0xfd, 0xe9, 0x03, 0xdc, 0x27, 0xd9, 0xfe, 0x3d, 0xb6, 0xeb, 0x03, 0x48, 0x99,
		0x73, 0x0b, 0xd5, 0x2e, 0x6c, 0x64, 0xa7, 0x2f, 0x77, 0x33, 0x1e, 0xa3, 0x09, 0x0c, 0x24, 0xda,
		0x84, 0xc2, 0xc8, 0x8a, 0x89, 0x42, 0x1a, 0xdc, 0xf1, 0x4d, 0xf4, 0xa2, 0xbf, 0x4e, 0x59, 0x39,
		0x6a, 0xe0, 0x8d, 0x5f, 0xff, 0x74, 0x37, 0x8e, 0x42, 0x3a, 0x43, 0x04, 0xb2, 0x7a, 0x16, 0x9e,
		0xd3, 0x95, 0x1d, 0x71, 0xd7, 0x29, 0x08, 0xa3, 0x2c, 0xcc, 0x45, 0x4e, 0xf0, 0xce, 0x84, 0x3f,
		0x4b, 0x51, 0xe6, 0x72, 0x18, 0xd6, 0x65, 0x17, 0xa1, 0x25, 0x1d, 0xad, 0x24, 0x84, 0x05, 0x5a,
		0x0d, 0xac, 0x11, 0x1d, 0x9d, 0xaa, 0x70, 0xf7, 0x54, 0xb1, 0x69, 0x27, 0x40, 0x57, 0x0e, 0xd0,
		0x62, 0xda, 0xa0, 0x14, 0x38, 0x7d, 0x42, 0x96, 0x40, 0x1b, 0x70, 0x03, 0xab, 0x0a, 0x33, 0xc5,
		0xdd, 0xa6, 0x5b, 0xbf, 0xc3, 0xf5, 0x3c, 0xcf, 0x73, 0xdb, 0xae, 0xc4, 0xe2, 0xa6, 0x27, 0x52,
		0xa1, 0xdc, 0x2b, 0x58, 0xa4, 0x91, 0x9c, 0x2c, 0x4f, 0xd6, 0x0c, 0x27, 0x57, 0x55, 0x8d, 0xfa,
		0x10, 0xdb, 0x12, 0x61, 0x7f, 0xbb, 0xba, 0x90, 0xc9, 0xc2, 0x14, 0x90, 0xc3, 0x27, 0x76, 0x6f,
		0xe0, 0x00, 0xb6, 0x0c, 0x04, 0x8c, 0xbc, 0xfa, 0x9c, 0xb7, 0x8a, 0x59, 0x92, 0xda, 0x55, 0xaf,
		0xd7, 0x2e, 0x02, 0x87, 0x29, 0x76, 0x9d, 0x3f, 0x44, 0xfe, 0x7f, 0xdd, 0x1a, 0x35, 0x80, 0xf5,
		0xf0, 0x46, 0x66, 0xf6, 0x86, 0x55, 0xc9, 0x4f, 0x2d, 0x5d, 0x01, 0x26, 0xa2, 0xff, 0xa2, 0x84,
		0xb2, 0xf9, 0x17, 0x3b, 0x15, 0xb3, 0x05, 0x8e, 0x54, 0x66, 0x5e, 0xdf, 0xa9, 0x77, 0x92, 0x0e,
		0xf4, 0x73, 0x1b, 0x07, 0x90, 0xb0, 0xd8, 0xe3, 0xc1, 0xca, 0x12, 0xd0, 0xe5, 0x6f, 0x03, 0x03,
		0x76, 0xca, 0xd8, 0x2b, 0x93, 0xe0, 0x7d, 0x3f, 0xce, 0x13, 0xb9, 0x47, 0xbd, 0xef, 0x79, 0xf6,
		0xb4, 0x91, 0x98, 0x4d, 0x6b, 0xe4, 0x6d, 0x4f, 0xdc, 0x65, 0x10, 0x07, 0xbc, 0x43, 0xe3, 0x4d,
		0x22, 0x75, 0xe6, 0x97, 0x2d, 0x6e, 0x62, 0x53, 0x20, 0x54, 0xba, 0x53, 0x0f, 0x0e, 0xdb, 0x0d,
		0xc7, 0x7c, 0x1c, 0x70, 0x29, 0x9b, 0xbd, 0x67, 0x93, 0xcd, 0x3e, 0xd9, 0xcb, 0xc9, 0x58, 0x45,
		0x73, 0x0a, 0x3f, 0xfa, 0x21, 0xf5, 0x25, 0xe2, 0x44, 0x6b, 0x4e, 0x68, 0xe4, 0x12, 0x11, 0xdb,
	}
	if len(w.ByteMap) != Mapsiz {
		w.ReadTable("lrx256.dat")
	}


	copy(w.ByteMap[:],byteMap)
}

