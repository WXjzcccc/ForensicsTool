package readers

import (
	"fmt"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	"github.com/iancoleman/orderedmap"
	"github.com/syndtr/goleveldb/leveldb"
)

func formatBytes(data []byte) string {
	var builder strings.Builder
	builder.WriteString("b'")
	for _, b := range data {
		if b >= 32 && b <= 126 { // 可打印的 ASCII 字符
			builder.WriteByte(b)
		} else {
			builder.WriteString(fmt.Sprintf("\\x%02x", b)) // 不可打印字符转为 \x 形式
		}
	}
	builder.WriteString("'")
	return builder.String()
}

func LevelDBReader(dbPath string) (*orderedmap.OrderedMap, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil || db == nil {
		logger.Errorf("【leveldb】数据库打开失败:%v", err)
		return nil, err
	}
	iter := db.NewIterator(nil, nil)
	result := orderedmap.New()
	result.SetEscapeHTML(false)
	var tmp []*orderedmap.OrderedMap
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		info := orderedmap.New()
		info.SetEscapeHTML(false)
		info.Set("键", formatBytes(key))
		info.Set("值", formatBytes(value))
		tmp = append(tmp, info)
	}
	iter.Release()
	defer db.Close()
	result.Set("leveldbReader", tmp)
	return result, err
}
