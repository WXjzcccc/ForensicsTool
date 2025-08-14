package cracker

import (
	"context"
	"fmt"
	"time"
)

// 助记词顺序恢复是单独的，就不放进来了

type ForensicsCracker struct {
	ctx         context.Context
	startTime   time.Time
	endTime     time.Time
	crackCancel context.CancelFunc
	crackCtx    context.Context
	crackResult *CrackResult
}

type CrackResult struct {
	Result string `json:"result"`
	Time   string `json:"time"`
	Error  string `json:"error"`
}

var CrackState string
var CalCount int = 0
var DoneCount int = 0
var Counts int = 0

func NewForensicsCracker() *ForensicsCracker {
	return &ForensicsCracker{}
}

func (f *ForensicsCracker) InitCtx(ctx context.Context) {
	f.ctx = ctx
}

func (f *ForensicsCracker) CancelCrack() {
	/*
		取消任务
	*/
	f.crackCancel()
	f.endTime = time.Now()
}

func (f *ForensicsCracker) GetState() string {
	return fmt.Sprintf("状态：%s\n耗时：%v\n信息：{\n\t结果：%v\n\t耗时：%v\n\t错误：%v\n}",
		CrackState, time.Since(f.startTime), f.crackResult.Result, f.crackResult.Time, f.crackResult.Error)
}
