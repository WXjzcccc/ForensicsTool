package cracker

import (
	"ForensicsTool/utils"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func getHash(text string) string {
	return utils.MD5HashHex("mm" + text)
}

func (f *ForensicsCracker) crackUin(target []string, start int, end int, c chan string, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i <= end; i++ {
		CrackState = fmt.Sprintf("爆破范围:%d-%d\n", start, end)
		select {
		case <-ctx.Done():
			runtime.Goexit()
		default:
			_hash := getHash(fmt.Sprintf("%d", i))
			Current++
			f.sendEvent()
			for _, v := range target {
				if _hash == v {
					c <- fmt.Sprintf("[√]<%s>找到了UIN: %d\t", _hash, i)
					DoneCount++
					if DoneCount == len(target) {
						Current = TotalCount
						f.sendProgress()
						c <- "[+]全部爆破完成\n"
						close(c)
					}
				}
			}
		}

	}
	// 完成一个区域的爆破
	CalCount++
	if CalCount == Counts {
		c <- "[+]爆破结束\n"
		close(c)
	}
}

func (f *ForensicsCracker) CrackWXUin(target []string) *CrackResult {
	/*
		@param target:	用户数据库所在目录，是一个md5哈希值
	*/
	f.startTime = time.Now()
	f.crackCtx, f.crackCancel = context.WithCancel(context.Background())
	defer f.crackCancel()
	f.crackResult = &CrackResult{}
	c := make(chan string, 1)
	var wg sync.WaitGroup
	go func() {
		for i := range c {
			f.crackResult.Result += i + "\n"
			CrackState = f.crackResult.Result
			f.crackResult.Time = time.Since(f.startTime).String()
			f.crackResult.Error = ""
			f.sendState()
		}
		f.crackCancel()
		runtime.Goexit()
	}()
	TotalCount = 4_294_967_296
	Current = 0
	Frequency = TotalCount / 1_000

	for i := -2147483648; i <= 2147483647; i += 100000000 {
		wg.Add(1)
		if i+99999999 < 2147483647 {
			go f.crackUin(target, i, i+99999999, c, f.crackCtx, &wg)
		} else {
			go f.crackUin(target, i, 2147483647, c, f.crackCtx, &wg)
		}
		Counts++
	}
	wg.Wait()
	Current = TotalCount
	f.sendProgress()
	if f.crackResult.Result == "" {
		f.crackResult.Time = time.Since(f.startTime).String()
		f.crackResult.Error = "未成功爆破"
	}
	return f.crackResult
}
