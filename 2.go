package main

import (
	"encoding/hex"
	"fmt"
	"github.com/shiyanhui/dht"
	"net/http"
	_ "net/http/pprof"
	"log"
	"database/sql"
	"time"
	_ "github.com/Go-SQL-Driver/MySQL"
	"strconv"
	//"reflect"
	//"encoding/json"
	"unicode"
)

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func initMysql() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/news?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertParam(db *sql.DB, torrent string, filename string, filesize int) {
	//defer db.Close()
	var insertSql string = "insert into `torrent`(`torrent`, `filename`, `init_time`, `filesize`) values(?, ?, ?, ?)";
	//fmt.Println("insert sql", reflect.TypeOf(torrent), reflect.TypeOf(filename), time.Now().Format("2006-01-02 15:04:05"), reflect.TypeOf(strconv.Itoa(filesize)))
	result, _ := db.Exec(insertSql, torrent, filename, time.Now().Format("2006-01-02 15:04:05"), strconv.Itoa(filesize))
	result.RowsAffected()
}

func queryExist(db *sql.DB, torrent string) (count int) {
	//defer db.Close()
	//row, err := db.Query("select count(1) from user where id = ?", idParam)
	row, err := db.Query("select count(1) from `torrent` where `torrent` = ?", torrent)
	if err != nil {
		log.Fatal(err)
	}
	//c, _ := row.RowsAffected()
	//log.Println("add affected rows:", c)
	return checkCount(row)
}
func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
var db = initMysql()

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	w := dht.NewWire(65536, 1024, 256)
	go func() {
		for resp := range w.Response() {
			metadata, err := dht.Decode(resp.MetadataInfo)
			if err != nil {
				continue
			}
			info := metadata.(map[string]interface{})

			if _, ok := info["name"]; !ok {
				continue
			}

			bt := bitTorrent{
				InfoHash: "magnet:?xt=urn:btih:" + hex.EncodeToString(resp.InfoHash),
				Name:     info["name"].(string),
			}

			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]file, len(files))

				for i, item := range files {
					f := item.(map[string]interface{})
					bt.Files[i] = file{
						Path:   f["path"].([]interface{}),
						Length: f["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}
			if !IsChineseChar(bt.Name) {
				fmt.Println("过滤了非中文" + bt.Name)
			} else {
				//bt是对象，需要写入到table
				go func() {
					if 0 == queryExist(db, bt.InfoHash) {
						fmt.Println("insert name" + bt.Name)
						insertParam(db, bt.InfoHash, bt.Name, bt.Length)
					} else {
						fmt.Println("过滤了" + bt.Name)
					}
				}()
			}
			//data, err := json.Marshal(bt)
			//if err == nil {
			//	fmt.Printf("%s\n\n", data)
			//}
		}
	}()
	go w.Run()

	config := dht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)

	d.Run()
}
