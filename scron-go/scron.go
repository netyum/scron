package main

import (
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	syscmd "scron-go/system"
	cron "scron-go/conparser"
	"time"
	"strings"
	"os"
	"strconv"
	"path"
	"net"
)

var db *sql.DB
var stmt *sql.Stmt
var ip string

func getString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func gethostname() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return name
}

func gethostbyname(addr string) string {
	ns, err := net.LookupHost(addr)
	if err != nil {
		fmt.Println("err: %s", err.Error())
		return ""
	}

	for _, n := range ns {
		return n
		break
	}
	return ""
}

func main() {
	hostname := gethostname()
	localip := gethostbyname(hostname)

	var configfile *string = flag.String("c", "/etc/scron.conf", "specify a config file name with path")
	var waringfile *string = flag.String("w", "/etc/scron_waring.conf", "specify a waring config file name with path")
	var environment *string = flag.String("e", "pro", "specify an environment name. dev/test/pro pro was default")

	flag.Parse()
	waring_config, err := ini.InsensitiveLoad(*waringfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	phones, err := waring_config.Section(*environment).GetKey("phones")
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	warning_script, err := waring_config.Section(*environment).GetKey("warning_script")
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	phonesString := phones.String()
	warningScriptString := warning_script.String()

	//获取配置信息
	config, err := ini.InsensitiveLoad(*configfile)
	if err != nil {
		syscmd.Warning(warningScriptString, fmt.Sprintf("scron读取配置文件：%s出错 来自主机: %s, ip: %s", *configfile, hostname, localip), phonesString)
		os.Exit(4)
	}

	db_host, err := config.Section(*environment).GetKey("db_host")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	db_name, err := config.Section(*environment).GetKey("db_name")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	db_user, err := config.Section(*environment).GetKey("db_user")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	db_pass, err := config.Section(*environment).GetKey("db_pass")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	db_port, err := config.Section(*environment).GetKey("db_port")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	log_path, err := config.Section(*environment).GetKey("log_path")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	cron_table_name, err := config.Section(*environment).GetKey("cron_table_name")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	host_table_name, err := config.Section(*environment).GetKey("host_table_name")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	iskill, err := config.Section(*environment).GetKey("iskill")
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	var iskill_bool bool
	iskill_bool, err = iskill.Bool()
	if err != nil {
		iskill_bool = false
	}

	syscmd.Mkdir(log_path.String())
	date_path := time.Now().Format("2006/01/02")
	syscmd.Mkdir(strings.Join([]string{log_path.String(), date_path}, "/"))


	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
	db, err = sql.Open("mysql", dataSourceName)
	defer db.Close()

	if err != nil {
		msg := fmt.Sprintf("scron连接MySQL出错从主机：%s IP：%s 连接到 %s", hostname, localip, db_host)
		syscmd.Warning(warningScriptString, msg, phonesString)
		os.Exit(6)
	}

	sql1 := fmt.Sprintf("SELECT * FROM `%s` %s, `%s` %s ", cron_table_name.String(), "c", host_table_name.String(), "h")
	sql2 := sql1 + " WHERE c.host=h.host and c.host=?"

	stmt, _ := db.Prepare(sql2)
	rows, err := stmt.Query(localip) //localip

	defer stmt.Close()
	defer rows.Close()
	if err != nil {
		fmt.Println("select data error: %v\n", err)
		os.Exit(6)
	}

	var cronId int
	var task_point *string
	var active int
	var mhdmd_point *string
	var command_point *string
	var params_point *string
	var process int
	var isQueue int
	var runAt_point *string
	var host_point *string
	var logFile_point *string
	var timeout int
	var user_point *string
	var created_at_point *string
	var updated_at_point *string
	var errorLogUpdateSize_point *string
	var LogUpdateSize_point *string
	var receiverPhone_point *string
	var receiverWx_point *string
	var hostId int
	var hostName_point *string
	var hostip_point *string
	var isEnable int
	for rows.Next() {
		err = rows.Scan(&cronId, &task_point, &active, &mhdmd_point, &command_point, &params_point, &process,
			&isQueue, &runAt_point, &host_point, &logFile_point, &timeout, &user_point, &created_at_point, &updated_at_point,
				&errorLogUpdateSize_point, &LogUpdateSize_point, &receiverPhone_point, &receiverWx_point,
				&hostId, &hostName_point, &hostip_point, &isEnable)
		if err != nil {
			fmt.Println(err)
			os.Exit(7)
		}
		mhdmd := getString(mhdmd_point)
		task := getString(task_point)
		user :=  getString(user_point)
		if user == "" {
			user = "root"
		}
		command := getString(command_point)
		command = strings.TrimSpace(command)

		params := getString(params_point)
		params = strings.TrimSpace(params)

		host_name := getString(hostName_point)
		last_size := getString(errorLogUpdateSize_point)
		last_size = strings.TrimSpace(last_size)
		var last_size_int int
		if last_size == "" {
			last_size_int = 0
		} else {
			last_size_int, err = strconv.Atoi(last_size)
			if err != nil {
				last_size_int = 0
			}
		}
		receiverPhone := getString(receiverPhone_point)
		receiverPhone = strings.TrimSpace(receiverPhone)
		if receiverPhone == "" {
			receiverPhone = phonesString
		}

		runAt := getString(runAt_point)
		runAt_int, err := strconv.Atoi(runAt)

		logFile := getString(logFile_point)
		logFile = strings.TrimSpace(logFile)
		if logFile == "" {
			logFile = strconv.Itoa(cronId) + ".log"
		}

		if err != nil {
			runAt_int = 0
		}

		//check timeout
		if timeout != 0 {
			msg := fmt.Sprintf("taskId %d %s 在 %s 上运行时间已经超过了设置时间 %d 分钟。",  cronId, task, host_name, timeout)
			syscmd.CheckTimeout(timeout, warningScriptString, command, params, msg, receiverPhone, iskill_bool)

		}

		if active == 1 && cron.Parse(mhdmd) && cron.IsRun(int64(runAt_int)) {
			fileExt := path.Ext(logFile)
			l := len(logFile) - len(fileExt)
			fileMain := logFile[:l]

			logErrorFile := fileMain+"_error"+fileExt

			num, err := syscmd.ProcessNum(command, params)

			if err != nil {
				continue
			}

			if num >= process {
				continue
			}

			log := strings.Join([]string{log_path.String(), date_path, logFile}, "/")
			error_log := strings.Join([]string{log_path.String(), date_path, logErrorFile}, "/")

			syscmd.Run(user, command, params, log, error_log)

			error_log_msg := fmt.Sprintf("taskId %d %s 在 %s 上运行的错误日志递增。请及时查看最新错误日志。", cronId, task, host_name)
			current_size, err := syscmd.CheckErrorLog(warningScriptString, error_log_msg, error_log, receiverPhone, last_size_int)
			if err != nil {
				fmt.Println(err)
			}

			if current_size != last_size_int {
				update_sql := fmt.Sprintf("UPDATE `%s` SET ", cron_table_name.String()) + " errorLogUpdatedSize = ? WHERE cronId= ?"
				stmt, _ := db.Prepare(update_sql)
				_, err := stmt.Exec(strconv.Itoa(current_size), cronId)
				defer stmt.Close()
				if err != nil {
					fmt.Println("update data error: %v\n", err)
					continue
				}
			}

			update_sql2 := fmt.Sprintf("UPDATE `%s` SET ", cron_table_name) +  "runAt = ? WHERE cronId=?"
			stmt, _ := db.Prepare(update_sql2)
			_, err = stmt.Exec(time.Now().Unix(), cronId)
			defer stmt.Close()
			if err != nil {
				fmt.Println("update data error2: %v\n", err)
			}
		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		os.Exit(8)
	}
}