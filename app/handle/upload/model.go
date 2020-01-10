package upload

import (
	"gfs/app/common"
	"strconv"
)

type SmallFUResult struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Url  string `json:"url"`
	Size int64  `json:"size"`
}

func newSmallFUResult(name, hash, url string, size int64) *SmallFUResult {
	return &SmallFUResult{Name: name, Hash: hash, Url: url, Size: size}
}

type BigFUIResult struct {
	IsExist        uint8    `json:"isExist"` //1：存在 2：不存在
	Name           string   `json:"name"`
	Hash           string   `json:"hash"`
	Url            string   `json:"url"`
	ChunkInfoArray []string `json:"chunkInfoArray"`
}

func newBigFUIResultIsExist(isExist uint8, name, hash, url string) *BigFUIResult {
	return &BigFUIResult{IsExist: isExist, Name: name, Hash: hash, Url: url}
}

func newBigFUIResultUnExist(isExist uint8, chunkInfoArray *[]string) *BigFUIResult {
	return &BigFUIResult{IsExist: isExist, ChunkInfoArray: *chunkInfoArray}
}

type BigFUCResult struct {
	FileHash  string `json:"fileHash"`
	ChunkHash string `json:"chunkHash"`
	Size      int64  `json:"size"`
}

func newBigFUCResult(fileHash, chunkHash string, size int64) *BigFUCResult {
	return &BigFUCResult{FileHash: fileHash, ChunkHash: chunkHash, Size: size}
}

type bigFileInfo struct {
	name         string
	hash         string
	chunkCount   uint16
	chunkInfoMap *map[string]*ChunkInfo
}

func newBigFileInfo(name, hash, chunkCount string) (*bigFileInfo, error) {
	c, err := strconv.ParseInt(chunkCount, 10, 16)
	if err != nil {
		e := common.NewGfsError("param chunkSize conversion to int16 failed")
		return nil, &e
	}
	chunkInfoMap := make(map[string]*ChunkInfo, uint16(c))
	return &bigFileInfo{name: name, hash: hash, chunkCount: uint16(c), chunkInfoMap: &chunkInfoMap}, nil
}

type ChunkInfo struct {
	Hash  string `json:"hash"`
	Index uint16 `json:"index"`
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

func newChunkInfo(hash, index, start, end string) (*ChunkInfo, error) {
	i, err := strconv.ParseInt(index, 10, 16)
	if err != nil {
		e := common.NewGfsError("param index conversion to int16 failed")
		return nil, &e
	}
	s, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		e := common.NewGfsError("param start conversion to int64 failed")
		return nil, &e
	}
	ed, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		e := common.NewGfsError("param index conversion to int64 failed")
		return nil, &e
	}
	return &ChunkInfo{Hash: hash, Index: uint16(i), Start: uint64(s), End: uint64(ed)}, nil
}
