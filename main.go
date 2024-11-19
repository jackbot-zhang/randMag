package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine

// TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.
var magList []string

func main() {
	magList = []string{}
	var err error
	engine, err = xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()
	err = engine.Sync(mag{})
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	txt := "请选择命令：\n 1.读取不删除\n 2.读取删除\n 3.标准输入\n 4.文件输入\n 5.退出\n"
	fmt.Print(txt)
	for {
		input, _ := reader.ReadString('\n') // 读取用户输入
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			printTen()
		case "2":
			printAndUpdate()
		case "3":
			fmt.Println("请输入磁力链接:")
			inputMode(reader)
			fmt.Print(txt)
		case "4":
			fmt.Println("请输入文件名:")
			inputFile(reader)
			time.Sleep(1 * time.Second)
			fmt.Printf("\n文件插入完成,%s", txt)
		case "5":
			return
		case "":
			fmt.Print(txt)
		default:
			fmt.Println("拜托，别乱选")
		}

	}

}
func inputFile(reader *bufio.Reader) {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	file, err := os.Open(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bt, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	arr := strings.Split(string(bt), "\n")
	magnetRegex := `magnet:\?xt=urn:btih:[0-9a-zA-Z]{32,40}`
	re := regexp.MustCompile(magnetRegex)
	cnt := 0
	for _, s := range arr {
		// 查找哈希值
		match := re.FindAllString(s, -1)
		if match != nil {
			for _, magnet := range match {
				_, err = engine.InsertOne(mag{Url: magnet})
				if err == nil {
					cnt++
				}
			}
		}
	}
	fmt.Printf("\n一共成功插入%d条数据\n", cnt)
}
func printTen() {
	list, err := selectTen()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, m := range list {
		fmt.Println(m.Url)
	}
}
func printAndUpdate() {
	list, err := selectTen()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, m := range list {
		fmt.Println(m.Url)
		engine.ID(m.ID).MustCols("view_cnt").Update(mag{ViewCnt: 1})
	}
}
func selectTen() ([]*mag, error) {
	var x []*mag
	err := engine.Table(mag{}).Where("view_cnt is null or  view_cnt is null = 0").Limit(10).Find(&x)
	return x, err
}
func inputMode(reader *bufio.Reader) {
	magnetRegex := `magnet:\?xt=urn:btih:[0-9a-zA-Z]{32,40}`
	re := regexp.MustCompile(magnetRegex)
	continuousInputMode := false
	cnt := 0
	var t1 *time.Time
	for {
		input, _ := reader.ReadString('\n') // 读取用户输入
		t2 := time.Now()
		if t1 != nil {
			if t2.Sub(*t1) < 100*time.Millisecond {
				continuousInputMode = true
				tmp := time.Now()
				t1 = &tmp
			} else {
				continuousInputMode = false
				t1 = nil
				fmt.Println("输入结束,共插入", cnt, "条数据。请按回车键")
				cnt = 0
			}
		}

		//
		input = strings.TrimSpace(input)
		if !continuousInputMode && input == "back" {
			return
		} else if !continuousInputMode && input == "" {
			fmt.Println("请输入磁力链接,或输入back返回:")
		} else {
			//第一次进来更新首次时间
			if t1 == nil {
				tmp := time.Now()
				t1 = &tmp
			}
			dealLine(input, re, &cnt)
		}
	}
}

func dealLine(line string, re *regexp.Regexp, cnt *int) {
	match := re.FindAllString(line, -1)
	if match != nil {
		for _, magnet := range match {
			_, err := engine.InsertOne(mag{Url: magnet})
			if err != nil {
				fmt.Println(err.Error())
			}
			if err == nil {
				*cnt++
			}
		}
	}

}

type mag struct {
	ID      int    `xorm:"id pk autoincr"`
	Url     string `xorm:"url unique"`
	ViewCnt int    `xorm:"view_cnt"`
}
