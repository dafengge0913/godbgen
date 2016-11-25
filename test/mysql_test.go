package test

import (
	"database/sql"
	"fmt"
	"github.com/dafengge0913/gocfg"
	"github.com/dafengge0913/godbgen/test/output"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func getDB(t *testing.T) (*sql.DB, error) {
	cfg, err := gocfg.ParseIni("demo_config.ini")
	if err != nil {
		return nil, fmt.Errorf("read config file error: %v", err)
	}
	db, err := sql.Open("mysql", cfg.GetString("conn_url")+cfg.GetString("db_name"))
	if err != nil {
		return nil, fmt.Errorf("connect db error: %v", err)
	}
	return db, nil
}

func TestMysqlInsert(t *testing.T) {

	db, err := getDB(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	di := &model.DemoItem{
		ItemId: 1003,
		Num:    11,
	}
	di, err = di.Insert(db)

	if err != nil {
		t.Errorf("Insert error:%v", err)
		return
	}

	t.Logf("Insert result :%v", di)
}

func TestMysqlUpdate(t *testing.T) {
	db, err := getDB(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	di := &model.DemoItem{
		Id:     1,
		ItemId: 1001,
		Num:    16,
	}

	rows, err := di.Update(db)

	if err != nil {
		t.Errorf("Update error:%v", err)
		return
	}

	t.Logf("Update rows :%d", rows)

}

func TestMysqlLoad(t *testing.T) {
	db, err := getDB(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	di := &model.DemoItem{
		Id:     1,
		ItemId: 1001,
	}

	di, err = di.Load(db)

	if err != nil {
		t.Errorf("Load error:%v", err)
		return
	}

	t.Logf("Load result :%v", di)

}

func TestMysqlDelete(t *testing.T) {
	db, err := getDB(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	di := &model.DemoItem{
		Id:     5,
		ItemId: 1002,
	}

	rows, err := di.Delete(db)

	if err != nil {
		t.Errorf("Delete error:%v", err)
		return
	}

	t.Logf("Delete rows :%d", rows)

}

func TestMysqlLoadAll(t *testing.T) {
	db, err := getDB(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	dis, err := model.LoadAllDemoItem(db)

	if err != nil {
		t.Errorf("LoadAll error:%v", err)
		return
	}

	for _, di := range dis {
		t.Logf("LoadAll result :%v", di)
	}
}
