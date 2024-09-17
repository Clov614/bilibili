// Package bilibili
// @Author Clover
// @Data 2024/9/1 下午10:24:00
// @Desc 链接解析
package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var firstUrl = "https://www.bilibili.com/video/"

var defaultUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15"

var (
	parseErr         = errors.New("parse error")
	invalidVideoInfo = errors.New("invalid video info")
)

type UrlParser struct {
	ua string
}

type VideoInfo struct {
	Bvid     string `json:"bvid"`     // 稿件bvid
	Aid      int    `json:"aid"`      // 稿件avid
	Videos   int    `json:"videos"`   // 稿件分P总数
	Tid      int    `json:"tid"`      // 分区tid
	Tname    string `json:"tname"`    // 子分区名称
	Pic      string `json:"pic"`      // 稿件封面图片url
	Title    string `json:"title"`    // 稿件标题
	PubDate  int    `json:"pubdate"`  // 稿件发布时间
	Duration int    `json:"duration"` // 稿件总时长(所有分P) 单位为秒
	View     int    `json:"view"`     // 播放数
	Danmaku  int    `json:"danmaku"`  // 弹幕数
	Reply    int    `json:"reply"`    // 评论数
	Favorite int    `json:"favorite"` // 收藏数
	Coin     int    `json:"coin"`     // 投币数
	Share    int    `json:"share"`    // 分享数
	Like     int    `json:"like"`     // 点赞数

	UpName string `json:"upname"` // up主昵称
}

func NewUrlDecoder() *UrlParser {
	return &UrlParser{ua: defaultUA}
}

// ParseByBvid 根据 bvid 解析视频信息
func (p *UrlParser) ParseByBvid(bvid string) (*VideoInfo, error) {
	return p.Parse(firstUrl + bvid)
}

// Parse 根据 url 解析视频信息
func (p *UrlParser) Parse(url string) (*VideoInfo, error) {
	data, err := p.doGet(url)
	if err != nil {
		return nil, fmt.Errorf("%w :%w", parseErr, err)
	}
	videoInfo, err := parseVideoInfo(data)
	if err != nil {
		return nil, fmt.Errorf("%w :%w", parseErr, err)
	}
	return videoInfo, nil
}

func parseVideoInfo(htmldata string) (*VideoInfo, error) {
	// 正则表达式匹配 VideoInfo 结构中的字段
	re := regexp.MustCompile(`"(?i)(bvid|aid|videos|tid|tname|pic|title|pubdate|duration|view|danmaku|reply|fav(orite)?|coin|share|like)":\s*("[^"]*"|\d+)`)

	// 用于将匹配结果构造成 JSON 字符串
	var jsonFields []string
	matches := re.FindAllStringSubmatch(htmldata, -1)
	// key去重
	var keyUnique = map[string]bool{}
	for _, match := range matches {
		key := match[1]
		value := match[3]

		// 将 fav 替换为 favorite
		if key == "fav" {
			key = "favorite"
		}
		if keyUnique[key] {
			continue
		}
		keyUnique[key] = true
		// 构造 JSON 字符串的一部分
		jsonFields = append(jsonFields, fmt.Sprintf(`"%s":%s`, key, value))
	}

	// 将匹配的字段拼接成 JSON 格式
	jsonString := "{" + strings.Join(jsonFields, ",") + "}"

	// 将 JSON 字符串解析为 VideoInfo 结构
	var videoInfo VideoInfo
	err := json.Unmarshal([]byte(jsonString), &videoInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	re2 := regexp.MustCompile(fmt.Sprintf(`"bvid":\s*"%s"[\s\S]*?"owner":[\s\S]*?("name"):\s*"([^"]*)"`, videoInfo.Bvid))
	submatch := re2.FindAllStringSubmatch(htmldata, -1)
	if submatch != nil && len(submatch) > 0 {
		if len(submatch[0]) >= 3 {
			videoInfo.UpName = fmt.Sprintf("%s", submatch[0][2])
		}

	}
	// 校验是否为异常状态
	if videoInfo.Bvid == "" || (videoInfo.Tname == "" && videoInfo.View == 0) {
		return nil, fmt.Errorf("videoInfo is nil: %w", invalidVideoInfo)
	}
	return &videoInfo, nil
}

func (p *UrlParser) doGet(url string) (output string, err error) {
	var req *http.Request
	var resp *http.Response

	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", p.ua)

	resp, err = c.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get error: %w", err)
	}
	defer resp.Body.Close()
	//if resp == nil {
	//	log.Fatal("resp is nil")
	//}
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body err: %w", err)
	}
	output = string(readAll)
	return
}
