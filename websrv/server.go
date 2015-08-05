package websrv

import (
	"fmt"
	"github.com/grayzone/example/util"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type SrvOp struct {
	Port           string
	TemplateFolder string
}

func (s *SrvOp) pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pwd)
	return pwd
}

func (s *SrvOp) staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func (s *SrvOp) renderHandler(w http.ResponseWriter, r *http.Request, templatename string) {
	t, err := template.ParseFiles(s.pwd() + s.TemplateFolder + templatename)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func (s *SrvOp) viewHandler(w http.ResponseWriter, r *http.Request) {
	s.renderHandler(w, r, "index.html")
}

func (s *SrvOp) dbpingHandler(w http.ResponseWriter, r *http.Request) {
	var ops util.DBOps
	ops.Init()
	err := ops.Open()
	if err != nil {
		fmt.Fprint(w, "open the database failed.")
		log.Fatal(err)
	}
	defer ops.Close()
	err = ops.Ping()
	if err != nil {
		fmt.Fprint(w, "ping the database failed.")
		log.Fatal(err)
	}
	fmt.Fprint(w, "ok")
}

func (s *SrvOp) dbGetDemoHandler(w http.ResponseWriter, r *http.Request) {

	//	symbol := r.FormValue("symbol")
	util.DownloadOneStock("sh600005")
	fmt.Fprint(w, "ok")
}

func (s *SrvOp) getstockHandler(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	symbol := r.FormValue("symbol")
	t, _ := time.Parse("2006-01-02", date)
	for i := 0; i < 100; i++ {
		t := t.AddDate(0, 0, i)
		go util.DownloadOneStockPerDay(symbol, t.Format("2006-01-02"))
	}
	fmt.Fprint(w, "ok")
}

func (s *SrvOp) dataCopyHandler(w http.ResponseWriter, r *http.Request) {

	u := url.URL{Scheme: "http", Host: "market.finance.sina.com.cn", Path: "/downxls.php"}
	q := u.Query()
	q.Set("date", "2014-05-01")
	q.Set("symbol", "000001")
	u.RawQuery = q.Encode()

	fmt.Fprint(w, u.String())
}

func (s *SrvOp) InitRoute() {
	http.HandleFunc(s.TemplateFolder, s.staticHandler)

	http.HandleFunc("/", s.viewHandler)

	http.HandleFunc("/dbping", s.dbpingHandler)
	http.HandleFunc("/dbgetdemo", s.dbGetDemoHandler)
	http.HandleFunc("/getstock", s.getstockHandler)
	http.HandleFunc("/datacopy", s.dataCopyHandler)
}

func (s *SrvOp) Start() {

	http.ListenAndServe(s.Port, nil)
}
