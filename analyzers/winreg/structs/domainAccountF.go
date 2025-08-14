package structs

import (
	"ForensicsTool/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ghostiam/binstruct"
)

type DomainAccountF struct {
	Revision                 uint32
	_                        uint32
	CreationTime             string `bin:"ParseWindowsFileTimestamp"`
	DomainModifiedCount      uint64
	MaxPwdAge                string `bin:"ParseWindowsFileTimestamp"`
	MinPwdAge                string `bin:"ParseWindowsFileTimestamp"`
	ForceLogOff              string `bin:"ParseWindowsFileTimestamp"`
	LockOutDuration          uint64
	LockOutObservationWindow uint64
	_                        uint64
	NextRid                  uint32
	PwdProperties            uint32
	MinPwdLength             uint16
	PwdHistoryLength         uint16
	LockoutTreshold          uint16
	_                        uint16
	ServerState              uint32
	ServerRole               uint16
	UasCompatibilityReq      uint16
	_                        uint64
	Key0                     interface{} `bin:"ParseSamKeyData"`
}

type SamKeyData struct {
	Revision uint32
	Length   uint32
	Salt     []byte `bin:"len:16"`
	Key      []byte `bin:"len:16"`
	CheckSum []byte `bin:"len:16"`
	Reserved uint64
}

func getSamKeyData(binData []byte) *SamKeyData {
	var samKeyData SamKeyData
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&samKeyData)
	if err != nil {
		return nil
	}
	return &samKeyData
}

type SamKeyDataAes struct {
	Revision       uint32
	Length         uint32
	CheckSumLength uint32
	DataLength     uint32
	Salt           []byte `bin:"len:16"`
	Data           []byte `bin:"len:DataLength"`
}

func getSamKeyDataAes(binData []byte) *SamKeyDataAes {
	var samKeyDataAes SamKeyDataAes
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&samKeyDataAes)
	if err != nil {
		return nil
	}
	return &samKeyDataAes
}

func (*DomainAccountF) ParseWindowsFileTimestamp(r binstruct.Reader) (string, error) {
	ts, err := r.ReadUint64()
	if err != nil {
		return "", err
	}
	return utils.WindowsFileTimeToDatetime(ts, "", ""), nil
}

func (d *DomainAccountF) ParseSamKeyData(r binstruct.Reader) (interface{}, error) {
	marker, err := r.Peek(1)
	if err != nil {
		return nil, err
	}
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	switch marker[0] {
	case byte(1):
		return getSamKeyData(data), nil
	case byte(2):
		return getSamKeyDataAes(data), nil
	}
	return nil, fmt.Errorf("invalid marker")
}

func GetDomainAccountF(binData []byte) *DomainAccountF {
	var domainAccountF DomainAccountF
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&domainAccountF)
	if err != nil {
		return nil
	}
	return &domainAccountF
}
