package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	records []*Record
)

type Record struct {
	// 郵便番号(ハイフンなし)
	ZipCode string `json:"zipCode"`
	// 都道府県名
	Address1 string `json:"address1"`
	// 市区町村名
	Address2 string `json:"address2"`
	// 町域名
	Address3 string `json:"address3"`
	// 都道府県名カナ
	Kana1 string `json:"kana1"`
	// 市区町村名カナ
	Kana2 string `json:"kana2"`
	// 町域名カナ
	Kana3 string `json:"kana3"`
}

func loadKenAll(filename string) ([]*Record, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))

	ret := make([]*Record, 0)
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		record := &Record{}
		for i, v := range r {
			switch i {
			case 2:
				// 郵便番号
				record.ZipCode = v
			case 3:
				// 都道府県名(カナ)
				record.Kana1 = v
			case 4:
				// 市区町村名(カナ)
				record.Kana2 = v
			case 5:
				// 町域名(カナ)
				record.Kana3 = v
			case 6:
				// 都道府県名
				record.Address1 = v
			case 7:
				// 市区町村名
				record.Address2 = v
			case 8:
				// 町域名
				record.Address3 = v
			}
		}
		// 各種フィルタを通す
		filterBlankAddress3(record)
		filterComeAfterAddress3(record)
		filterItienAddress3(record)
		filterSpecialWordAddress3(record)
		filterRangeSpecAddress3(record)

		ret = append(ret, record)
	}
	// 郵便番号重複行を取り除く
	ret = removeDuplicateRec(ret)

	return ret, nil
}

// 指定rune以降消す
func trimAfter(s string, r rune) string {
	idx := strings.IndexRune(s, r)
	if idx >= 0 {
		return s[:idx]
	}
	return s
}

// 郵便番号重複行を取り除いて住所3を開始カッコ以降消す
// 面倒すぎてやる気がでない...いっそ消してしまおう...
// 郵便番号で並んでいる前提なのですでにだめだが目をつぶる...
func removeDuplicateRec(records []*Record) []*Record {
	ret := make([]*Record, 0)
	var prev *Record

	for _, r := range records {
		if prev == nil {
			prev = r
			ret = append(ret, r)
			continue
		}
		if r.ZipCode == prev.ZipCode {
			prev.Address3 = trimAfter(prev.Address3, '（')
			prev.Kana3 = trimAfter(prev.Kana3, '(')
			continue
		}
		ret = append(ret, r)
		prev = r
	}

	return ret
}

// 町域名が「以下に掲載がない場合」の場合は、ブランクに置換
func filterBlankAddress3(r *Record) *Record {
	if r.Address3 == "以下に掲載がない場合" {
		r.Address3 = ""
		r.Kana3 = ""
	}
	return r
}

// 「（次のビルを除く）」を取り除く
func filterSpecialWordAddress3(r *Record) *Record {
	if strings.Contains(r.Address3, "（次のビルを除く）") {
		r.Address3 = strings.Replace(r.Address3, "（次のビルを除く）", "", 1)
		r.Kana3 = strings.Replace(r.Kana3, "(ﾂｷﾞﾉﾋﾞﾙｦﾉｿﾞｸ)", "", 1)
	}
	return r
}

// 町域名に「の次に番地がくる場合」を含む場合は、ブランクに置換
func filterComeAfterAddress3(r *Record) *Record {
	if strings.HasSuffix(r.Address3, "の次に番地がくる場合") {
		r.Address3 = ""
		r.Kana3 = ""
	}
	return r
}

// 町域名「（または町・村）一円」の場合は、ブランクに置換
// ただし「一円」が地名(滋賀県犬上郡多賀町一円)である場合は置換しない
func filterItienAddress3(r *Record) *Record {
	if strings.HasSuffix(r.Address3, "一円") &&
		r.Address1 != "滋賀県" &&
		r.Address2 != "犬上郡多賀町" {
		r.Address3 = ""
		r.Kana3 = ""
	}
	return r
}

// 範囲指定している場合は住所3を開始カッコ以降消す
// いっそ消してしまおう
func filterRangeSpecAddress3(r *Record) *Record {
	if strings.Contains(r.Address3, "～") ||
		strings.Contains(r.Address3, "以上") ||
		strings.Contains(r.Address3, "以下") ||
		strings.Contains(r.Address3, "以外") ||
		strings.Contains(r.Address3, "その他") ||
		strings.Contains(r.Address3, "（全域）") ||
		strings.Contains(r.Address3, "（丁目）") ||
		strings.Contains(r.Address3, "（群）") ||
		strings.Contains(r.Address3, "（番地）") ||
		strings.Contains(r.Address3, "（○○屋敷）") ||
		strings.Contains(r.Address3, "を含む）") ||
		strings.Contains(r.Address3, "を除く") {
		r.Address3 = trimAfter(r.Address3, '（')
		r.Kana3 = trimAfter(r.Kana3, '(')
	}
	return r
}

// 前方一致でlimit分探す
func findKenAll(records []*Record, q string, limit int) []*Record {
	ret := make([]*Record, 0)
	for _, v := range records {
		if strings.HasPrefix(v.ZipCode, q) {
			ret = append(ret, v)
			limit--
			if limit <= 0 {
				return ret
			}
		}
	}
	return ret
}
