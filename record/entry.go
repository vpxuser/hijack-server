package record

import (
	"encoding/json"
	"fmt"
	"github.com/vpxuser/proxy"
	"hijack-server/setting"
	"io"
	"os"
	"time"
)

const dir_path = "report"

func init() {
	//defer setting.Template.Close()
	if setting.Cfg.Report {
		err := os.MkdirAll(dir_path, 0755)
		if err != nil {
			proxy.Fatal(err)
		}

		filename := time.Now().Format("2006-01-02_15_04_05")
		path := fmt.Sprintf("%s/%s.html", dir_path, filename)
		file, err := os.Create(path)
		if err != nil {
			proxy.Fatal(err)
		}

		_, err = io.Copy(file, setting.Template)
		if err != nil {
			proxy.Fatal(err)
		}
		setting.Template.Close()

		proxy.Infof("已生成测试报告，路径：%s", path)
		go insert(file)
	}
}

type Entry struct {
	CreateTime int64   `json:"create_time"`
	Detail     *Detail `json:"detail"`
	Plugin     string  `json:"plugin"`
	Target     Target  `json:"target"`
}

type Detail struct {
	Addr     string            `json:"addr"`
	Payload  string            `json:"payload"`
	Snapshot [][]string        `json:"snapshot"`
	Extra    map[string]string `json:"extra"`
}

type Target struct {
	Url string `json:"url"`
}

func NewEntry() *Entry {
	return &Entry{
		CreateTime: time.Now().UnixMilli(),
		Detail: &Detail{
			Extra: make(map[string]string),
		},
	}
}

var report = make(chan *Entry)

func Push(entry *Entry) {
	report <- entry
}

const format = `<script class='web-vulns'>webVulns.push(%s)</script>`

func insert(file *os.File) {
	defer file.Close()
	for entry := range report {
		record, err := json.Marshal(entry)
		if err != nil {
			proxy.Error(err)
		}
		line := fmt.Sprintf(format, record)
		_, err = file.WriteString(line + "\n")
		if err != nil {
			proxy.Error(err)
		}
	}
}
