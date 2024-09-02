// Package bilibili
// @Author Clover
// @Data 2024/9/1 下午10:32:00
// @Desc
package bilibili

import (
	"fmt"
	"strings"
	"testing"
)

var exampleurl = "https://www.bilibili.com/video/"

func Test_doGet(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name               string
		args               args
		wantContainStrList []string
	}{
		{
			name: "test01",
			args: args{
				url: exampleurl,
			},
			wantContainStrList: []string{"BV12i421r7WW", "aid", "\"reply\":", "\"view\":"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, _ := NewUrlDecoder().doGet(tt.args.url)
			for _, wantContainStr := range tt.wantContainStrList {
				if strings.Contains(output, wantContainStr) == false {
					t.Errorf("case: %s doGet() not contain %v", tt.name, wantContainStr)
				}
			}
		})
	}
}

func TestUrlParser_Parse(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *VideoInfo
		wantErr bool
	}{
		{
			name: "test01",
			args: args{
				url: exampleurl + "BV12i421r7WW",
			},
			want: &VideoInfo{
				Title:    "第四章究竟讲了什么故事",
				View:     930000,
				Coin:     12000,
				Like:     29000,
				Favorite: 12000,
			},
		},
		{
			name: "test02",
			args: args{
				url: exampleurl + "BV1bZ421L7eu",
			},
			want: &VideoInfo{
				Title:    "牡蛎也能用来造桥",
				View:     390000,
				Coin:     3954,
				Like:     20000,
				Favorite: 4034,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewUrlDecoder()
			got, err := p.Parse(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("case: %s parse is nil", tt.name)
			}
			switch true {
			case !strings.Contains(got.Title, tt.want.Title):
				t.Errorf("case: %s got.title(%s) Contain (%s)", tt.name, got.Title, tt.want.Title)
			case got.View < tt.want.View:
				t.Errorf("case: %s got.view(%d)<(%d)", tt.name, got.View, tt.want.View)
			case got.Coin < tt.want.Coin:
				t.Errorf("case: %s got.coin(%d)<(%d)", tt.name, got.Coin, tt.want.Coin)
			case got.Like < tt.want.Like:
				t.Errorf("case: %s got.like(%d)<(%d)", tt.name, got.Like, tt.want.Like)
			case got.Favorite < tt.want.Favorite:
				t.Errorf("case: %s got.fav(%d)<(%d)", tt.name, got.Favorite, tt.want.Favorite)
			}
			t.Log(fmt.Sprintf("videoInfo: %+v", got))
		})
	}
}
