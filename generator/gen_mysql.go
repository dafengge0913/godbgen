package generator

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/dafengge0913/golog"
	"github.com/dafengge0913/gotypes/set"
	_ "github.com/go-sql-driver/mysql"
	"go/format"
	"html/template"
	"os"
	"strings"
)

type mysqlGen struct {
	log            *golog.Logger
	dataSourceName string
	databaseName   string
	packageName    string
	outputPath     string
}

func NewMysqlGen(log *golog.Logger, dataSourceName, databaseName, packageName, outputPath string) *mysqlGen {
	return &mysqlGen{
		log:            log,
		dataSourceName: dataSourceName,
		databaseName:   databaseName,
		packageName:    packageName,
		outputPath:     outputPath,
	}
}

func (g *mysqlGen) Generate() error {
	db, err := sql.Open("mysql", g.dataSourceName)
	if err != nil {
		return fmt.Errorf("connect to db error: %v", err)
	}
	defer db.Close()

	tables := g.getAllTable(db, g.databaseName)
	for _, tName := range tables {
		if err := g.genModel(db, tName); err != nil {
			return err
		}
	}
	return nil
}

func (g *mysqlGen) getAllTable(db *sql.DB, dbName string) []string {
	names := make([]string, 0)
	if rows, err := db.Query("select table_name from information_schema.tables where table_schema = ?", dbName); err != nil {
		g.log.Error("getAllTable error: %v", err)
		return names
	} else {
		defer rows.Close()
		for rows.Next() {
			var tName string
			if err := rows.Scan(&tName); err != nil {
				g.log.Error("%v", err)
				return names
			}
			g.log.Debug("find table [%s]", tName)
			names = append(names, tName)
		}
	}
	return names
}

type tableStruct struct {
	colName string
	colType string
	colKey  string
}

func (g *mysqlGen) getTableStruct(db *sql.DB, tableName string) ([]*tableStruct, error) {
	tss := make([]*tableStruct, 0)
	if rows, err := db.Query("select column_name, column_type,column_key from information_schema.columns where table_schema = ? and table_name = ? ", g.databaseName, tableName); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			ts := &tableStruct{}
			if err := rows.Scan(&ts.colName, &ts.colType, &ts.colKey); err != nil {
				return nil, err
			}
			tss = append(tss, ts)
		}
	}
	return tss, nil
}

type modelField struct {
	FieldName string
	FieldType string
}

func (g *mysqlGen) genModel(db *sql.DB, tableName string) error {
	data := make(map[string]interface{})
	modelName := camelFormat(tableName)
	data["ModelName"] = modelName
	data["PackageName"] = g.packageName
	tmpl, err := template.ParseFiles("tmpl/model.tmpl")
	if err != nil {
		return fmt.Errorf("open template error:%v ", err)
	}
	tss, err := g.getTableStruct(db, tableName)
	if err != nil {
		return err
	}

	importSet := set.NewSet()
	mfs := make([]*modelField, len(tss), len(tss))
	for i, ts := range tss {

		goType := g.sqlType2GoType(ts.colType)
		if pkg, ok := g.genImports(goType); ok {
			importSet.Add(pkg)
		}

		mfs[i] = &modelField{
			FieldName: camelFormat(ts.colName),
			FieldType: goType,
		}
	}

	data["ModelFields"] = mfs
	data["Imports"] = importSet.List()

	file, err := os.OpenFile(g.outputPath+string(os.PathSeparator)+modelName+".go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("write model file error:%v ", err)
	}
	defer file.Close()
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("template execute error:%v ", err)
	}

	fmtCode, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format code error:%v ", err)
	}
	_, err = file.Write(fmtCode)
	return err
}

func camelFormat(name string) string {
	words := strings.Split(name, "_")
	result := ""
	for _, word := range words {
		if len(word) < 1 {
			continue
		}
		result += strings.ToUpper(string(word[0])) + word[1:]
	}
	return result
}

func (g *mysqlGen) sqlType2GoType(colType string) string {
	if colType == "tinyint(1)" {
		return "bool"
	}
	idx := strings.Index(colType, "(")
	var sqlType string
	if idx != -1 {
		sqlType = colType[:idx]
	} else {
		sqlType = colType
	}
	switch sqlType {
	case "int":
		return "int"
	case "tinyint":
		return "int8"
	case "smallint":
		return "int16"
	case "integer":
		return "int32"
	case "bigint":
		return "int64"
	case "float", "decimal":
		return "float32"
	case "double":
		return "float64"
	case "money":
		return "string"
	case "text", "varchar", "char":
		return "string"
	case "date", "time", "datetime":
		return "time.Time"
	case "timestamp":
		return "time.Duration"
	}
	return "string"
}

func (g *mysqlGen) genImports(goType string) (string, bool) {
	switch goType {
	case "time.Time", "time.Duration":
		return "time", true
	default:
		return "", false
	}
}
