package main

import (
	"database/sql"
	// "encoding/json"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// alter table cs_apply_module_cache add public_type int DEFAULT NULL COMMENT '发布的类型，0测试，1正式，null：弱依赖';
// 配置文件数据结构
type Cssdbconf struct {
	Dbstr string
}

// Post请求结构
type ReqData struct {
	Org_no       string `json:"org_no"`
	Account_name string `json:"account_name"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/runsql", runsqlHandler)

	infoLog.Println("Listening...")
	errorLog.Fatal(http.ListenAndServe(":9999", mux))
	// 获取应用配置
	//myconf := getToml()
	//DB, _ := sql.Open("mysql", "sxadmin:sx@123@tcp(192.168.199.100:3306)/cs_s_run")
	//createDb(myconf)

}

func runsqlHandler(w http.ResponseWriter, r *http.Request) {
	 myconf := getToml()
	/*
		//jsonStr, _ := json.Marshal(reqdata)
		//infoLog.Println("req json: ", string(jsonStr))
		dbarr, err := GetDb(myconf)
		if err != nil {
			return
		}

		tmpl, err := template.ParseFiles("./static/home.tmpl")
		if err != nil {
			infoLog.Println("create template failed, err:", err)
			return
		}
		// 利用给定数据渲染模板, 并将结果写入w
		tmpl.Execute(w, dbarr)
	*/
	r.ParseForm()
	// infoLog.Println("r.Formstr : ", r.Form["qstr"][0])
	sqlstr := r.Form["qstr"][0]
	sqlid, err := insertSql(myconf, sqlstr)
	if err != nil {
		infoLog.Println("insert t_sql failed, err:",sqlid, err)
		return
	}
/*
	for name, value := range r.Form {
		 infoLog.Println(name,":", value[0])
		if strings.HasPrefix(name, "cs_s_run") {
			err := runSql(myconf, sqlstr, name)
			if err != nil {
				insertSqllog(myconf, name, sqlid, 0)
			} else {
				insertSqllog(myconf, name, sqlid, 1)
			}
		} 
	}
 */
	w.Write([]byte(r.Method + "\n"))

	return

}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	myconf := getToml()

	dbarr, err := GetDb(myconf)
	if err != nil {
		return
	}

	tmpl, err := template.ParseFiles("./static/home.tmpl")
	if err != nil {
		infoLog.Println("create template failed, err:", err)
		return
	}
	// 利用给定数据渲染模板, 并将结果写入w
	tmpl.Execute(w, dbarr)

	return

}

func getToml() (c Cssdbconf) {

	var conf Cssdbconf
	if _, err := toml.DecodeFile("./conf/cssdb.toml", &conf); err != nil {
		// handle error
		errorLog.Fatal(err)
	}
	return conf
}

// 获得所有的创世数据库名
func GetDb(c Cssdbconf) (dba []string, e error) {
	var dbname string
	dbarr := []string{"cs_s_run"}
	DB, _ := sql.Open("mysql", c.Dbstr)
	defer DB.Close()

	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		infoLog.Println("open database failed...")
	}

	infoLog.Println("connnect cs_s_run database success...")

	sqlStr := "show databases"
	rows, err := DB.Query(sqlStr)
	if err != nil {
		errorLog.Println("show databases failed...")
	}
	//循环显示所有的数据
	for rows.Next() {
		rows.Scan(&dbname)
		if strings.HasPrefix(dbname, "cs_s_run_") {
			dbarr = append(dbarr, dbname)
		}
	}

	infoLog.Println("--", dbarr)
	return dbarr, err
}

// 写入cs_s_update数据库的t_sql表,返回最新的ID
func insertSql(c Cssdbconf, s string) (sid int64, err error) {

	connstr := c.Dbstr + "cs_s_update"
	DB, _ := sql.Open("mysql", connstr)
	defer DB.Close()
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	sql := "insert into t_sql(esql) values (?)"
	r, err := DB.Exec(sql, s)
	if err != nil {
		infoLog.Println("exec failed,", err)
	}

	//查询最后一天用户ID，判断是否插入成功
	var id int64
	id, err = r.LastInsertId()
	if err != nil {
		infoLog.Println("exec failed,", err)
	}
	return id, err
}

// 写入cs_s_update数据库的t_sqllog表
func insertSqllog(c Cssdbconf, s string, sqlid int64, isok int) {

	connstr := c.Dbstr + "cs_s_update"
	DB, _ := sql.Open("mysql", connstr)
	defer DB.Close()
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	sql := "insert into t_sqllog(sqlid,dbname,isDone) values (?,?,?)"
	_, err := DB.Exec(sql, sqlid, s, isok)
	if err != nil {
		infoLog.Println("exec failed,", err)
	}

}

// 执行SQL并 写入cs_s_update数据库的t_sqllog表
func runSql(c Cssdbconf, s string, dbname string) (err error) {

	connstr := c.Dbstr + dbname
	DB, _ := sql.Open("mysql", connstr)
	defer DB.Close()
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	_, err = DB.Exec(s)
	if err != nil {
		infoLog.Println("exec failed,", err)
	}
	return err
}
