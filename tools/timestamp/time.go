package timestamp

import (
	"ForensicsTool/utils"
	"context"
	"fmt"
	"log"
	"strconv"
)

type TimeStampParser struct {
	ctx context.Context
}

func NewTimeStampParser() *TimeStampParser {
	return &TimeStampParser{}
}

func (p *TimeStampParser) InitCtx(ctx context.Context) {
	p.ctx = ctx
}

func (p *TimeStampParser) ParseTimeStamp(timeStamp, oriZone, newZone string) string {
	log.Println("timeStamp", timeStamp)
	log.Println("oriZone", oriZone)
	log.Println("newZone", newZone)
	result := ""
	floatTs, err := strconv.ParseFloat(timeStamp, 64)
	if err != nil {
		floatTs = 0
	}
	i64Ts, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		i64Ts = 0
	}
	log.Println("floatTs", floatTs)
	log.Println("i64Ts", i64Ts)
	result += fmt.Sprintf("iOS数据库中的时间戳：%s \n", utils.IosTimestampToDatetime(floatTs, oriZone, newZone))
	result += fmt.Sprintf("AppleAbsoluteTime：%s \n", utils.AppleTimestampToDatetime(floatTs, oriZone, newZone))
	result += fmt.Sprintf("18位时间戳：%s \n", utils.NineTimestampToDatetime(i64Ts, oriZone, newZone))
	result += fmt.Sprintf("Chrome/Webkit：%s \n", utils.ChromeTimestampToDatetime(i64Ts, oriZone, newZone))
	result += fmt.Sprintf("WindowsFileTime：%s \n", utils.WindowsFileTimeToDatetime(uint64(i64Ts), oriZone, newZone))
	result += fmt.Sprintf("UNIX时间戳：%s \n", utils.DefaultTimestampToDatetime(i64Ts, oriZone, newZone))
	log.Println(result)
	return result
}
