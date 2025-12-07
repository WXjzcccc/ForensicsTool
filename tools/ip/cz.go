package ip

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/ipipdotnet/ipdb-go"
)

type CZIP struct {
	CityName      string `json:"city_name"`
	ContinentCode string `json:"continent_code"`
	CountryCode   string `json:"country_code"`
	CountryName   string `json:"country_name"`
	DistrictName  string `json:"district_name"`
	ISPDomain     string `json:"isp_domain"`
	OwnerDomain   string `json:"owner_domain"`
	RegionName    string `json:"region_name"`
	Error         string `json:"error"`
}

type IP struct {
	ctx     context.Context
	dbPath  string
	db      *ipdb.City
	version string
}

const tag_api = "https://api.github.com/repos/nmgliangwei/qqwry.ipdb/tags"
const db_url = "https://raw.gitmirror.com/nmgliangwei/qqwry.ipdb/main/qqwry.ipdb"
const (
	IPNewVersionEvent    = "ForensicsTool::IP::NewVersion"
	IPUpdateStartEvent   = "ForensicsTool::IP::UpdateStart"
	IPUpdateSuccessEvent = "ForensicsTool::IP::UpdateSuccess"
	IPUpdateErrorEvent   = "ForensicsTool::IP::UpdateError"
)

func NewIP() *IP {
	return &IP{dbPath: "qqwry.ipdb"}
}

func (i *IP) InitCtx(ctx context.Context) {
	i.ctx = ctx
}

// LoadDB 加载IP库，失败返回错误
func (i *IP) LoadDB() error {
	db, err := ipdb.NewCity(i.dbPath)
	if err != nil {
		logger.Errorf("【IP】加载纯真IP库<%s>失败：%v", i.dbPath, err)
		return err
	}

	i.db = db
	i.version = i.db.BuildTime().Format("2006-01-02")
	logger.Debugf("【IP】当前纯真IP库版本<%s>", i.version)
	return nil
}

// CheckUpdate 检查IP库是否有更新
func (i *IP) CheckUpdate() bool {
	req, _ := http.NewRequest("GET", tag_api, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("【IP】查询tag失败")
		return false
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("【IP】服务器返回错误: %s", resp.Status)
		return false
	}
	defer resp.Body.Close()
	var tags []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false
	}

	if len(tags) == 0 {
		logger.Error("【IP】获取tag失败")
		return false
	}
	tag := tags[0].Name

	if tag != i.version {
		logger.Infof("【IP】纯真IP库存在新版本：%s", tag)
		// runtime.EventsEmit(i.ctx, IPUpdateEvent, tag)
		return true
	}

	return false
}

// UpdateDB 更新IP库，失败返回错误
func (i *IP) UpdateDB() error {
	// runtime.EventsEmit(i.ctx, IPUpdateStartEvent, "开始下载")
	resp, err := http.Get(db_url)
	if err != nil {
		logger.Error("【IP】下载纯真IP库失败")
		// runtime.EventsEmit(i.ctx, IPUpdateErrorEvent, "下载失败")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("【IP】服务器返回错误: %s", resp.Status)
		// runtime.EventsEmit(i.ctx, IPUpdateErrorEvent, fmt.Sprintf("服务器返回错误: %s", resp.Status))
		return fmt.Errorf("服务器返回错误: %s", resp.Status)
	}
	file, err := os.Create(i.dbPath) //下载到软件根目录
	if err != nil {
		logger.Errorf("【IP】创建文件失败：%v", err)
		// runtime.EventsEmit(i.ctx, IPUpdateErrorEvent, fmt.Sprintf("【IP】创建文件失败：%v", err))
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		logger.Errorf("【IP】写入文件失败：%v", err)
		// runtime.EventsEmit(i.ctx, IPUpdateErrorEvent, fmt.Sprintf("【IP】创建文件失败：%v", err))
		return err
	}
	logger.Info("【IP】纯真IP库下载成功")
	// runtime.EventsEmit(i.ctx, IPUpdateSuccessEvent, "下载成功")
	return nil
}

// CheckDB 检查纯真IP库是否可用
func (i *IP) CheckDB() error {
	file, err := os.Stat(i.dbPath)
	if err != nil {
		logger.Errorf("【IP】检查纯真IP库<%s>失败：%v", i.dbPath, err)
		return err
	}
	if file.IsDir() {
		logger.Errorf("【IP】纯真IP库<%s>不是文件", i.dbPath)
		return err
	}
	if file.Size() == 0 {
		logger.Errorf("【IP】纯真IP库<%s>为空", i.dbPath)
		return err
	}
	return nil
}

func (i *IP) Search(ip string) *CZIP {
	if i.db == nil {
		logger.Error("【IP】数据库未初始化")
		return &CZIP{Error: "数据库未初始化"}
	}
	ipMap, err := i.db.FindMap(ip, "CN")
	if err != nil {
		logger.Errorf("【IP】查询纯真IP库失败：%v", err)
		return &CZIP{Error: err.Error()}
	}
	czip := &CZIP{
		CityName:      ipMap["city_name"],
		ContinentCode: ipMap["continent_code"],
		CountryCode:   ipMap["country_code"],
		CountryName:   ipMap["country_name"],
		RegionName:    ipMap["region_name"],
		ISPDomain:     ipMap["isp_domain"],
		DistrictName:  ipMap["district_name"],
		OwnerDomain:   ipMap["owner_domain"],
	}
	return czip
}

func (i *IP) SearchAll(ips []string) map[string]*CZIP {
	results := make(map[string]*CZIP)
	var wg sync.WaitGroup
	var mu sync.Mutex
	maxConcurrent := 10
	semaphore := make(chan struct{}, maxConcurrent)

	for _, ip := range ips {
		wg.Add(1)
		go func(ipAddr string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			result := i.Search(ipAddr)
			mu.Lock()
			results[ip] = result
			mu.Unlock()
		}(ip)
	}
	wg.Wait()
	return results
}
