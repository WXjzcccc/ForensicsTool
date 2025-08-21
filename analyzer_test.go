package main

import (
	"ForensicsTool/tools/timestamp"
	"log"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	//ext := &extractor.InfoExtractor{}
	//extResult := ext.ExtractXShell("C:\\Users\\Administrator\\Documents\\NetSarang Computer\\7\\Xshell\\Sessions", "AdministratorS-1-5-21-2019645882-1436685051-1453824686-500")
	//log.Println(extResult)
	//extResult2 := ext.ExtractNavicat("C:\\Users\\Administrator\\Desktop\\TODO\\NTUSER.DAT")
	//log.Println(extResult2)
	//a, err := tool.AnalyzeMobaXterm("C:\\Users\\Administrator\\Desktop\\TODO\\moba\\NTUSER.DAT", "12345678")
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(a)
	//b, err := tool.AnalyzeDbeaver("C:\\Users\\Administrator\\Desktop\\TODO\\计算机镜像2")
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(b)
	//log.Println(ext.ExtractFinalShell("C:\\Users\\Administrator\\Desktop\\workspace\\rootfs_\\conn"))
	//log.Println(ext.ExtractHawk2("C:\\Users\\Administrator\\Desktop\\TODO\\Hawk2.xml", "rp6Fz0G8E6p7KzfJS0qg3Agyoek/GoRg12kAptH9VDo="))
	//reg := &winreg.Reg{}
	//ret, err := reg.AnalyzeWinReg("C:\\Users\\Administrator\\Desktop\\TODO\\test")
	//if err != nil {
	//	log.Println(err)
	//}
	//for s, i := range ret {
	//	log.Println(s, i)
	//}

	//crack := &cracker.ForensicsCracker{}
	//ticker := time.NewTicker(time.Second * 10)
	//go func() {
	//	for _ = range ticker.C {
	//		log.Println(crack.GetState())
	//	}
	//}()
	////log.Println(crack.CrackAirDrop("5e4cf", "8f170", "86", []string{"139"}, 8))
	//log.Println(crack.CrackWXUin([]string{"8c6823a0ac2e55ed3a3832bb7610d5e1"}))
	//dec := &database.DecryptDatabase{}
	//
	//log.Println(dec.DecryptSystemDataSQLite("G:\\ZD20250709-001\\ZD20250709-001.mf", "MobileForensic@2013@"))

	tt := &timestamp.TimeStampParser{}
	log.Println(tt.ParseTimeStamp("1755761706", "UTC", "Asia/Shanghai"))
}
