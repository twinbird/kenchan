package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	AccessControlAllowOrigin string
	ResponseLimit            int
	AcceptPort               int
	CsvFilePath              string
	UpdateCsvFileUrl         string
	RefreshKen               bool
)

func init() {
	flag.StringVar(&AccessControlAllowOrigin, "a", "", "レスポンスに付与するAccess-Control-Allow-Originを指定します")
	flag.StringVar(&CsvFilePath, "f", "KEN_ALL.CSV", "利用するcsvファイルパスを指定します")
	flag.IntVar(&ResponseLimit, "l", 20, "検索結果件数のリミットを指定します")
	flag.IntVar(&AcceptPort, "p", 8080, "受付ポートを指定します")
	flag.StringVar(&UpdateCsvFileUrl, "u", "https://www.post.japanpost.jp/zipcode/dl/kogaki/zip/ken_all.zip", "csvファイル取得に利用するURLを指定します")
	flag.BoolVar(&RefreshKen, "r", false, "csvファイルを更新する際に指定します")
	flag.Parse()

	if RefreshKen == true || fileExists(CsvFilePath) == false {
		if err := downloadKenAll(UpdateCsvFileUrl, CsvFilePath); err != nil {
			log.Fatal(err)
		}
	}

	var err error
	records, err = loadKenAll(CsvFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	q = strings.Replace(q, "-", "", 1)
	found := findKenAll(records, q, ResponseLimit)

	ret, err := json.MarshalIndent(found, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if AccessControlAllowOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", AccessControlAllowOrigin)
	}
	w.Write(ret)
}

func main() {
	p := strconv.Itoa(AcceptPort)

	http.HandleFunc("/search", searchHandler)

	log.Println("Accepting port" + p)
	log.Fatal(http.ListenAndServe(":"+p, nil))
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
