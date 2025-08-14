package structs

import (
	"ForensicsTool/utils"
	"bytes"
	"encoding/binary"
	"github.com/ghostiam/binstruct"
)

// UserF _是暂时不知道含义的字段
type UserF struct {
	_                   []byte `bin:"len:8"`
	LastLoginTime       string `bin:"ParseWindowsFileTimestamp"`
	_                   []byte `bin:"len:8"`
	LastPwdChangeTime   string `bin:"ParseWindowsFileTimestamp"`
	_                   []byte `bin:"len:8"`
	LastFailedLoginTime string `bin:"ParseWindowsFileTimestamp"`
	RID                 uint32
	_                   []byte `bin:"len:4"`
	UserAttribute       uint32
	_                   []byte `bin:"len:4"`
	LogonCount          int16
	InValidLoginCount   int16
	_                   []byte `bin:"len:12"`
}

// ParseWindowsFileTimestamp 会在解码的时候被调用
func (*UserF) ParseWindowsFileTimestamp(r binstruct.Reader) (string, error) {
	ts, err := r.ReadUint64()
	if err != nil {
		return "", err
	}
	return utils.WindowsFileTimeToDatetime(ts, "", ""), nil
}

func GetUserF(binData []byte) *UserF {
	var userF UserF
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&userF)
	if err != nil {
		return nil
	}
	return &userF
}
