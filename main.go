package main

import (
	"github.com/dafengge0913/gocfg"
	"github.com/dafengge0913/godbgen/generator"
	"github.com/dafengge0913/golog"
)

func main() {

	log := golog.NewLogger(golog.LEVEL_DEBUG, nil)
	cfg, err := gocfg.ParseIni("test/demo_config.ini")
	if err != nil {
		log.Error("read config file error: %v", err)
		return
	}

	//get arguments from config file or command line or others
	gen := generator.NewMysqlGen(log, cfg.GetString("conn_url"), cfg.GetString("db_name"), "model", "test/output")

	if err = gen.Generate(); err != nil {
		log.Error("Generate error: %v", err)
		return
	}

}
