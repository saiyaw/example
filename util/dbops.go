package util

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"time"
)

type DBOps struct {
	Db         *sql.DB
	DriverName string
	Dbname     string
	User       string
	Password   string
}

type TableDemo struct {
	Id  int
	Val int
}

func (ops *DBOps) Init() {
	ops.Db = nil
	ops.DriverName = "postgres"
	ops.Dbname = "demo"
	ops.User = "postgres"
	ops.Password = "123456"
}

func (ops *DBOps) Open() (err error) {
	ops.Init()
	connstr := "user=" + ops.User + " password=" + ops.Password + " dbname=" + ops.Dbname + " sslmode=disable"
	ops.Db, err = sql.Open(ops.DriverName, connstr)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (ops *DBOps) Ping() (err error) {
	return ops.Db.Ping()
}

func (ops *DBOps) GetTableDemo() ([]TableDemo, error) {
	rows, err := ops.Db.Query("SELECT * FROM demo")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var td []TableDemo
	for rows.Next() {
		var t TableDemo
		err = rows.Scan(&t.Id, &t.Val)
		if err != nil {
			log.Println(err)
		} else {
			td = append(td, t)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (ops *DBOps)UpdateSymbolTime(symbol string, t time.Time)error{
	
	_, err := ops.Db.Query("UPDATE symbol SET updatetime = $1 WHERE value = $2", t.Format("2006-01-02 15:04:05 MST"), symbol)
	if err != nil{
		log.Println(err.Error())
	}
	
	return err
	
}

func (ops *DBOps)GetSymbols()([]Stock, error){
	
	rows, err := ops.Db.Query("SELECT id,value,name,updatetime FROM symbol WHERE valid = true")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var td []Stock
	for rows.Next() {
		var t Stock
		err = rows.Scan(&t.Id, &t.Symbol, &t.Name, &t.Updatetime)
		if err != nil {
			log.Println(err)
		} else {
			td = append(td, t)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (ops *DBOps) AddOneSymbol(symbol string, name string, valid bool) error {
	txn, err := ops.Db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer txn.Rollback()

	stmt, err := txn.Prepare(pq.CopyIn("symbol", "value", "name", "valid", "updatetime"))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = stmt.Exec(symbol, name, valid, time.Now())
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (ops *DBOps) CopyDataIntoTable(t []Stock) {
	txn, err := ops.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer txn.Rollback()

	stmt, err := txn.Prepare(pq.CopyIn("stock", "dealtime", "price", "volume", "amount", "climax"))
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range t {
		_, err = stmt.Exec(item.Dealtime, item.Price, item.Volume, item.Amount, item.Climax)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (ops *DBOps) Close() (err error) {
	return ops.Db.Close()
}
