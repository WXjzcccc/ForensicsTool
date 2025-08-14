package cracker

import (
	"ForensicsTool/utils"
	"context"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var macsDefault = []string{"130", "131", "132", "133", "135", "136", "137", "138", "139", "140", "141", "144", "145", "146", "147", "148",
	"149", "150", "151", "152", "153", "155", "156", "157", "158", "159", "166", "167", "171", "172", "173", "174",
	"175", "176", "177", "178", "180", "181", "182", "183", "184", "185", "186", "187", "188", "189", "190", "191",
	"192", "193", "195", "196", "197", "198", "199"}

func getHeadAndTail(phone string) (string, string) {
	hash := utils.SHA256HashHex(phone)
	return hash[:5], hash[len(hash)-5:]
}

func crackPhone(region, mac, head, tail string, length int, c chan string, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range int64(math.Pow(10, float64(length))) {
		CrackState = fmt.Sprintf("爆破号段:%s\n", mac)
		select {
		case <-ctx.Done():
			runtime.Goexit()
		default:
			fmtStr := "%s%s%0" + strconv.Itoa(length) + "d"
			phone := fmt.Sprintf(fmtStr, region, mac, i)
			hashHead, hashTail := getHeadAndTail(phone)
			if hashHead == head && hashTail == tail {
				c <- fmt.Sprintf("[√]已找到匹配的手机号：%s\n", phone)
			}
		}
	}
	c <- fmt.Sprintf("[+]%s开头的手机号爆破结束\n", mac)
	DoneCount += 1
	if DoneCount == Counts {
		// 避免爆破完成后阻塞，虽然阻塞也没什么关系，后面没有逻辑就是了
		close(c)
	}
}

func (f *ForensicsCracker) CrackAirDrop(head, tail, region string, macs []string, length int) *CrackResult {
	/*
		@param	head:	前5位哈希值
		@param	tail:	后5位哈希值
		@param	region:	区号，默认86
		@param	macs:	号段
		@param	length:	除区号和号段后的号码长度，默认8
	*/
	f.startTime = time.Now()
	f.crackCtx, f.crackCancel = context.WithCancel(context.Background())
	defer f.crackCancel()
	f.crackResult = &CrackResult{}
	c := make(chan string, 1)
	if region == "" {
		region = "86"
	}
	if len(macs) == 0 {
		macs = macsDefault
	}
	if length == 0 {
		length = 8
	}
	var wg sync.WaitGroup
	go func() {
		for i := range c {
			f.crackResult.Result += i + "\n"
			CrackState = f.crackResult.Result
			f.crackResult.Time = time.Since(f.startTime).String()
			f.crackResult.Error = ""
		}
		f.crackCancel()
		runtime.Goexit()
	}()
	for _, mac := range macs {
		wg.Add(1)
		go crackPhone(region, mac, head, tail, length, c, f.crackCtx, &wg)
	}
	wg.Wait()
	if f.crackResult.Result == "" {
		f.crackResult.Time = time.Since(f.startTime).String()
		f.crackResult.Error = "未成功爆破"
	}
	return f.crackResult
}
