package structs

import (
	"ForensicsTool/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ghostiam/binstruct"
	"io"
)

type UserV struct {
	_                    []byte `bin:"len:12"`
	NameOffset           uint32
	NameLength           uint32
	_                    uint32
	FullNameOffset       uint32
	FullNameLength       uint32
	_                    uint32
	CommentOffset        uint32
	CommentLength        uint32
	_                    uint32
	UserCommentOffset    uint32
	UserCommentLength    uint32
	_                    uint32
	_                    []byte `bin:"len:12"`
	HomeDirOffset        uint32
	HomeDirLength        uint32
	_                    uint32
	HomeDirConnectOffset uint32
	HomeDirConnectLength uint32
	_                    uint32
	ScriptPathOffset     uint32
	ScriptPathLength     uint32
	_                    uint32
	ProfilePathOffset    uint32
	ProfilePathLength    uint32
	_                    uint32
	WorkstationsOffset   uint32
	WorkstationsLength   uint32
	_                    uint32
	HoursAllowedOffset   uint32
	HoursAllowedLength   uint32
	_                    uint32
	_                    []byte `bin:"len:12"`
	LMHashOffset         uint32
	LMHashLength         uint32
	_                    uint32
	NTHashOffset         uint32
	NTHashLength         uint32
	_                    uint32
	_                    []byte      `bin:"len:24"`
	Name                 string      `bin:"ParseName"`
	FullName             string      `bin:"ParseFullName"`
	Comment              string      `bin:"ParseComment"`
	UserComment          string      `bin:"ParseUserComment"`
	HomeDir              string      `bin:"ParseHomeDir"`
	HomeDirConnect       string      `bin:"ParseHomeDirConnect"`
	ScriptPath           string      `bin:"ParseScriptPath"`
	ProfilePath          string      `bin:"ParseProfilePath"`
	Workstations         string      `bin:"ParseWorkstations"`
	HoursAllowed         string      `bin:"ParseHoursAllowed"`
	LMHash               interface{} `bin:"ParseLMHash"`
	NTHash               interface{} `bin:"ParseNTHash"`
}

type SAMHash struct {
	PekID    uint16
	Revision uint16
	Hash     []byte `bin:"len:16"`
}

func getSAMHash(binData []byte) *SAMHash {
	var samHash SAMHash
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&samHash)
	if err != nil {
		return nil
	}
	return &samHash
}

type SAMHashAes struct {
	PekID      uint16
	Revision   uint16
	DataOffset uint32
	Salt       []byte `bin:"len:16"`
	Data       []byte `bin:"ParseData"`
}

func (s *SAMHashAes) ParseData(r binstruct.Reader) ([]byte, error) {
	binData, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return binData, nil
}

func getSAMHashAes(binData []byte) *SAMHashAes {
	var samHashAes SAMHashAes
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&samHashAes)
	if err != nil {
		return nil
	}
	return &samHashAes
}

const VOffset = 204 // UserV固定读取前204个字节

func (u *UserV) ParseName(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.NameOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.NameLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseFullName(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.FullNameOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.FullNameLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseComment(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.CommentOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.CommentLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseUserComment(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.UserCommentOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.UserCommentLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseHomeDir(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.HomeDirOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.HomeDirLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseHomeDirConnect(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.HomeDirConnectOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.HomeDirConnectLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseScriptPath(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.ScriptPathOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.ScriptPathLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseProfilePath(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.ProfilePathOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.ProfilePathLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseWorkstations(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.WorkstationsOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.WorkstationsLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseHoursAllowed(r binstruct.Reader) (string, error) {
	_, err := r.Seek(int64(u.HoursAllowedOffset+VOffset), io.SeekStart)
	if err != nil {
		return "", err
	}
	binData, err := r.Peek(int(u.HoursAllowedLength))
	if err != nil {
		return "", err
	}
	return utils.UTF16leBytesToString(binData), nil
}

func (u *UserV) ParseLMHash(r binstruct.Reader) (interface{}, error) {
	_, err := r.Seek(int64(u.NTHashOffset+VOffset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	head, err := r.Peek(1)
	if err != nil {
		return nil, err
	}
	if head[0] == byte(1) {
		if u.LMHashLength == 20 {
			_, err = r.Seek(int64(u.LMHashOffset+VOffset), io.SeekStart)
			if err != nil {
				return nil, err
			}
			binData, err := r.Peek(int(u.LMHashLength))
			if err != nil {
				return nil, err
			}
			return getSAMHash(binData), nil
		}
	} else {
		if u.LMHashLength == 24 {
			_, err = r.Seek(int64(u.LMHashOffset+VOffset), io.SeekStart)
			if err != nil {
				return nil, err
			}
			binData, err := r.Peek(int(u.LMHashLength))
			if err != nil {
				return nil, err
			}
			return getSAMHashAes(binData), nil
		}
	}
	return nil, fmt.Errorf("err lmhash")
}

func (u *UserV) ParseNTHash(r binstruct.Reader) (interface{}, error) {
	_, err := r.Seek(int64(u.NTHashOffset+VOffset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	head, err := r.Peek(1)
	if err != nil {
		return nil, err
	}
	if head[0] == byte(1) {
		if u.NTHashLength == 20 {
			_, err = r.Seek(int64(u.NTHashOffset+VOffset), io.SeekStart)
			if err != nil {
				return nil, err
			}
			binData, err := r.Peek(int(u.NTHashLength))
			if err != nil {
				return nil, err
			}
			return getSAMHash(binData), nil
		}
	} else {
		_, err = r.Seek(int64(u.NTHashOffset+VOffset), io.SeekStart)
		if err != nil {
			return nil, err
		}
		binData, err := r.Peek(int(u.NTHashLength))
		if err != nil {
			return nil, err
		}
		return getSAMHashAes(binData), nil
	}
	return nil, fmt.Errorf("err nthash")
}

func GetUserV(binData []byte) *UserV {
	var userV UserV
	decoder := binstruct.NewDecoder(bytes.NewReader(binData), binary.LittleEndian)
	err := decoder.Decode(&userV)
	if err != nil {
		return nil
	}
	return &userV
}
