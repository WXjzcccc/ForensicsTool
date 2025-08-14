package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ghostiam/binstruct"
)

type DomainAccountV struct {
	_        []byte `bin:"len:408"`
	DomainID uint32
	SID1     uint32
	SID2     uint32
	SID3     uint32
}

func GetMachineSid(binData []byte) string {
	var domainAccountV DomainAccountV
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&domainAccountV)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("S-1-5-%v-%v-%v-%v", domainAccountV.DomainID, domainAccountV.SID1, domainAccountV.SID2, domainAccountV.SID3)
}
