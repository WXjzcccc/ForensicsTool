package reader

import (
	"ForensicsTool/analyzers/reader/readers"
	"context"

	"github.com/iancoleman/orderedmap"
)

type FileReader struct {
	ctx context.Context
}

func NewFileReader() *FileReader {
	return &FileReader{}
}

type ReadMapResult struct {
	Data *orderedmap.OrderedMap `json:"data"`
	Err  string                 `json:"err"`
}

func (p *FileReader) InitCtx(ctx context.Context) {
	p.ctx = ctx
}

func (p *FileReader) ReadLevelDB(path string) *ReadMapResult {
	result, err := readers.LevelDBReader(path)
	if err != nil {
		return &ReadMapResult{nil, err.Error()}
	}
	return &ReadMapResult{result, ""}
}

func (p *FileReader) ReadMMKV(path, password string) *ReadMapResult {
	result, err := readers.MMKVReader(path, password)
	if err != nil {
		return &ReadMapResult{nil, err.Error()}
	}
	return &ReadMapResult{result, ""}
}
