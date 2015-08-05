package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Stock struct {
	Id         int64
	Symbol     string
	Name       string
	Dealtime   time.Time
	Price      float64
	Volume     int64
	Amount     float64
	Climax     string
	Updatetime time.Time
}

func DownloadOneStockPerDay(symbol string, date string) error {
	u := url.URL{Scheme: "http", Host: "market.finance.sina.com.cn", Path: "/downxls.php"}
	q := u.Query()
	q.Set("date", date)
	q.Set("symbol", symbol)
	u.RawQuery = q.Encode()

	log.Println(u.String())

	td := GetStockInOneDay(u.String())
	if len(td) > 0 {

		var ops DBOps
		ops.Init()
		err := ops.Open()
		if err != nil {
			log.Println("open the database failed.")
			log.Fatal(err)
		}
		defer ops.Close()

		ops.CopyDataIntoTable(td)
	} else {
		log.Printf("%s:%s#No DATA\n", symbol, date)
	}

	return nil
}

func DownloadOneStock(symbol string) {
	t := time.Date(1994, time.January, 1, 0, 0, 0, 0, time.Local)

	d := time.Since(t)
	for d.Hours()/24.0 > 1 {
		t = t.AddDate(0, 0, 1)
		d = time.Since(t)
		log.Println(t)
		DownloadOneStockPerDay(symbol, t.Format("2006-01-02"))
		time.Sleep(10 * time.Millisecond)
	}
	var ops DBOps
	ops.Init()
	err := ops.Open()
	if err != nil {
		log.Println("open the database failed.")
		log.Fatal(err)
	}
	defer ops.Close()
	ops.UpdateSymbolTime(symbol, t)
}

func GenerateSymbol() {
	var ops DBOps
	ops.Init()
	err := ops.Open()
	if err != nil {
		log.Println("open the database failed.")
		log.Fatal(err)
	}
	defer ops.Close()

	for i := 0; i < 5000; i++ {
		s := 600000 + i
		symbol := fmt.Sprintf("sh%d", s)
		time.Sleep(10 * time.Millisecond)
		name, valid := GetNameBySymbol(symbol)
		ops.AddOneSymbol(symbol, name, valid)
	}

	for i := 0; i < 2000; i++ {
		symbol := fmt.Sprintf("sz%06d", i)
		time.Sleep(10 * time.Millisecond)
		name, valid := GetNameBySymbol(symbol)
		ops.AddOneSymbol(symbol, name, valid)
	}

}

func GetNameBySymbol(s string) (string, bool) {
	u := url.URL{Scheme: "http", Host: "hq.sinajs.cn", Path: "/"}
	q := u.Query()
	q.Set("list", s)
	u.RawQuery = q.Encode()
	log.Println(u.String())

	res, err := http.Get(u.String())
	if err != nil {
		log.Println(err.Error())
		return "", false
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return "", false
	}
	//	log.Println(GBKtoUTF8(string(body)))
	content := string(body)
	if strings.Contains(content, ",") != true {
		log.Println("Invalid symbol")
		return "", false
	}
	name := content[strings.Index(content, "\"")+1 : strings.Index(content, ",")]
	log.Println(GBKtoUTF8(name))

	//	log.Printf("length:%d\n", len(name))

	return GBKtoUTF8(name), true

}

func GetStockInOneDay(urlstr string) []Stock {
	res, err := http.Get(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//	log.Println(GBKtoUTF8(string(body)))
	contents := strings.Split(string(body), "\n")

	u, _ := url.Parse(urlstr)
	q := u.Query()
	date := q["date"][0]

	var td []Stock
	for _, line := range contents[1:] {

		content := strings.Split(line, "\t")
		//		log.Printf("%d:%s\n",len(content),content)
		if len(content) < 6 {
			//			log.Println("invalid format :" + stringutil.GBKtoUTF8(line))
			continue
		}

		var t Stock
		t.Dealtime, _ = time.Parse("2006-01-02 15:04:05 MST", date+" "+content[0]+" CST")
		//		log.Println(t.Dealtime.String())
		t.Price, err = strconv.ParseFloat(content[1], 64)
		t.Volume, _ = strconv.ParseInt(content[3], 0, 64)
		t.Amount, _ = strconv.ParseFloat(content[4], 64)
		t.Climax = GBKtoUTF8(content[5])
		//		log.Println(t)
		td = append(td, t)
	}

	return td
}
