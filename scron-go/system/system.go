//                       __                    __
// _      ______ ___  __/ /_  ____  ____  ____/ /
//| | /| / / __ '/ / / / __ \/ __ \/ __ \/ __  /
//| |/ |/ / /_/ / /_/ / / / / /_/ / /_/ / /_/ /
//|__/|__/\__,_/\__, /_/ /_/\____/\____/\__,_/
//             /____/
//

package system

import (
	"strconv"
	"os/exec"
	"strings"
	"fmt"
	"time"
	"os"
)

var cmdwithpath string
var err error

func init() {
	cmdwithpath, err = exec.LookPath("bash")
	if err != nil {
		fmt.Println("not find bash.")
		os.Exit(5)
	}
}

//杀掉进程
func Kill(pid int) error {
	cmdString := fmt.Sprintf("kill %s", strconv.Itoa(pid))
	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

//建立目录
func Mkdir(path string) error {
	cmdString := fmt.Sprintf(`mkdir -p %s > /dev/null 2>&1`, path)
	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

//进程数
func ProcessNum(command, params string) (int, error) {
	var cmdString string
	if strings.TrimSpace(params) == "" {
		cmdString = fmt.Sprintf( `ps -ef | grep -v "grep" | grep -- "%s" | wc -l`, command)

	} else {
		cmdString = fmt.Sprintf( `ps -ef | grep -v "grep" | grep -- "%s" | grep -- "%s$" | wc -l`, command, params)
	}
	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return -2, err
	}

	s := string(out)
	ss := strings.Split(s, "\n")
	for i := 0; i < len(ss); i++ {
		ti, err := strconv.Atoi(strings.TrimSpace(ss[i]))
		if err != nil {
			return -3, err
		} else {
			return ti, nil
		}
		break
	}
	return 0, nil
}

//运行
func Run(user, command, params, log, error_log string) error {
	if command == "" || params == "" {
		return nil
	}
	var cmdString string
	if user != "" && user != "root" {
		cmdString = fmt.Sprintf(cmdwithpath, "-c", "su %s -c '%s %s %s >> %s 2>>%s &'", user, cmdwithpath, command, params, log, error_log)
	} else {
		cmdString = fmt.Sprintf("setsid %s %s >> %s 2>>%s &", command, params, log, error_log)
	}
	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

//检测超时
func CheckTimeout(timeout int, warningCmd, command, params, msg, receiver string, isKill bool) error {
	var cmdString string

	if strings.TrimSpace(params) == "" {
		cmdString = fmt.Sprintf(`ps -eo pid,etime,cmd | grep -v grep | grep -- "%s" | awk '{print $2\"|\"$1}'`, command)
	} else {
		cmdString = fmt.Sprintf(`ps -eo pid,etime,cmd | grep -v grep | grep -- "%s" | grep -- "%s$" | awk '{print $2\"|\"$1}'`, command, params)
	}

	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	s := string(out)
	ss := strings.Split(s, "\n")
	for i := 0; i < len(ss); i++ {
		vs := strings.Split(strings.TrimSpace(ss[i]), "|")
		if len(vs) != 2 {
			continue
		}

		run_time := get_run_minute(vs[0])
		if run_time > timeout {
			newmsg := msg + `进程id: ` + vs[1]

			if isKill == true {
				pid, err := strconv.Atoi(vs[1])
				if err == nil {
					Kill(pid)
				}
			} else {
				Warning(warningCmd, newmsg, receiver)
			}
		}
	}
	return nil
}

//检测错误
func CheckErrorLog(warningCmd, msg, error_log, receiver string, last_size int) (int, error) {
	cmdString := fmt.Sprintf(`ls -l %s | awk '{print $5}'`, error_log)
	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	s := string(out)
	ss := strings.Split(s, "\n")
	for i := 0; i < len(ss); i++ {
		var current_size int
		current_size_s := strings.TrimSpace(ss[i])
		if current_size_s == "" {
			current_size = 0
		} else {
			current_size, err = strconv.Atoi(current_size_s)
			if err != nil {
				current_size = 0
			}
		}
		if current_size == 0 {
			return current_size, nil
		}

		if last_size != current_size {
			Warning(warningCmd, msg, receiver)
		}
		return current_size, nil
		break
	}
	return 0, nil
}

// 报警函数
func Warning(cmdString, msg, phones string) error {
	ymdhms := time.Now().Format("2006/01/02 15/04/05")
	timemsg := fmt.Sprintf("[%s] %s", ymdhms, msg)

	if cmdString == "" || msg == "" {
		return nil
	}

	cmdString = fmt.Sprintf(`%s "%s" "%s"`, cmdString, timemsg, phones)
	cmd := exec.Command(cmdwithpath, "-c", cmdString)

	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

//信息输出
func Info(msg string) {
	ymdhms := time.Now().Format("2006/01/02 15/04/05")
	timemsg := fmt.Sprintf("[%s] %s", ymdhms, msg)
	fmt.Println(timemsg)
}

func get_run_minute(str string) int {
	min := 0
	if strings.Index(str, "-") != -1 {
		bits := strings.Split(str, "-")
		if len(bits) == 2 {
			mInt, err := strconv.Atoi(bits[0])
			if err == nil {
				min = 60 * 24 * mInt
				str = bits[1]
			}
		}
	}

	if strings.Index(str, ":") != -1 {
		bits := strings.Split(str, ":")
		if len(bits) == 3 {
			mInt, err := strconv.Atoi(bits[0])
			if err == nil {
				min = min + (60 * mInt)
				mmInt, err := strconv.Atoi(bits[1])
				if err == nil {
					min = min + mmInt
				}
			} else if len(bits) == 2 {
				min = min + mInt
			}
		}
	}
	return min
}
