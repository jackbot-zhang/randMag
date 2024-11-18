package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	// 打开数据库连接
	db, err := sql.Open("sqlite3", "D:\\mycode\\randMag\\test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// 查询数据
	rows, err := db.Query("SELECT id, url  FROM projects where  view_cnt is null  limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var idList []int64
	for rows.Next() {
		var id int64
		var url string
		err = rows.Scan(&id, &url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(url)
		idList = append(idList, id)
	}
	// 检查查询过程中可能出现的错误
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("输入OK,更新并且退出")
	// 捕获 Ctrl+C 信号
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if strings.ToLower(input) == "ok" {
			// 更新数据
			for _, id := range idList {
				_, err = db.Exec("UPDATE projects SET view_cnt = ? WHERE id = ?", 1, id)
				if err != nil {
					log.Printf("Failed to update view_cnt for project ID %d: %v", id, err)
				}
			}
			log.Println("更新完成，程序退出")
			time.Sleep(time.Second * 2)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Reading standard input: %v", err)
	}

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
