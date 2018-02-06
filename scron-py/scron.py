#!/usr/bin/env python
# -*- coding: utf-8 -*-
# 
# filename: scron.py
# author: syang
#
 
import time
import os, sys
import getopt
import ConfigParser
import subprocess
import socket

reload(sys) 
sys.setdefaultencoding('utf-8')


hostname = socket.gethostname()
localip = socket.gethostbyname(hostname)
# 报警函数
def waring(msg, phones) :
    ymdhms = time.strftime("%Y-%m-%d %H:%M:%S")
    timemsg = "[%s] %s" % (ymdhms, msg)
    info(msg)

    waring_cmd = "/usr/bin/php /sp_edaijia/www/v2/protected/yiic sms ScriptsWaring --msg=\"%s\" --time=0 --phone=\"%s\"" % (timemsg, phones)
    os.system(waring_cmd)
    # 报警脚本，自己定义

def info(msg) :
    ymdhms = time.strftime("%Y-%m-%d %H:%M:%S")
    timemsg = "[%s] %s" % (ymdhms, msg)
    print timemsg

def usage() :
    print "usage: scron -c /etc/scron.conf -e pro"
    print "      -h was show help"
    print "      -c was specify a config file name with path"
    print "      -w was specify a waring config file with path. /etc/scron_waring.conf was default"
    print "      -e was specify an environment name. dev/test/pro pro was default"

try :
    import MySQLdb
    import MySQLdb.cursors
except ImportError, e:
    print "ImportError %s" % (e.args[0], )
    waring("scron导入MySQL库出错在主机: %s, ip: %s" % (hostname, localip))
    quit()


class System :
    def __init__(self):
        pass

    def kill(self, pid) :
        cmd = "/bin/kill %s " % (pid, )
        os.system(cmd)

    def mkdir(self, path) :
        os.system("/bin/mkdir -p %s > /dev/null 2>&1" % (path, ))

    def process_num(self, command, params) :
        if params.strip() == "" :
            cmd = "ps -ef | grep -v grep | grep -- \"%s$\" | wc -l" % (command, )
        else :
            print params
            cmd = "ps -ef | grep -v grep | grep -- \"%s\" | grep -- \"%s$\" | wc -l" % (command, params)
        print cmd
        handler = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE)
        num = int(handler.communicate()[0])
        return num

    def run(self, user, command, params, log, error_log) :
        if user != "" and user != "root" :
          cmd = "setsid su %s -c 'setsid %s %s >> %s 2>>%s &'" % (user, command, params, log, error_log)
        else :
          cmd = "setsid %s %s >> %s 2>>%s &" % (command, params, log, error_log)
        info( "%s %s" % (command, params))
        os.system(cmd)

    def check_timeout(self, command, params, msg, iskill, receiver) :
        global phones
        if params.strip() == "" :
            cmd = "ps -eo pid,etime,cmd | grep -v grep | grep -- \"%s\" | awk '{print $2\"|\"$1}'" % (command, )
        else :
            cmd = "ps -eo pid,etime,cmd | grep -v grep | grep -- \"%s\" | grep -- \"%s\" | awk '{print $2\"|\"$1}'" % (command, params)
    
        handler = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE)
        time_list = handler.communicate()[0].strip().split("\n")
        for v in time_list :
            vs = v.strip().split("|")
            if len(vs) != 2 :
                continue

            run_time = self.get_run_minute(vs[0])
            if run_time > timeout :
                newmsg = "%s 进程id: %s" % (msg, vs[1])
                waring_phones = phones
                if receiver != "" :
                    waring_phones = receiver

                if int(iskill) == 1 :
                    killmsg = "%s 已经被杀掉" % (newmsg, )
                    waring(killmsg, waring_phones)
                    syscmd.kill(vs[1])
                else :
                    waring(newmsg, waring_phones)

    def check_error_log(self, msg, error_log, last_size, receiver) :
        global phones
        waring_phones = phones
        if receiver != "" :
            waring_phones = receiver

        cmd = "ls -l %s | awk '{print $5}'" % (error_log, )
        handler = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE)
        current_size = str(handler.communicate()[0].strip())
        if current_size == "" :
            current_size = 0
        if current_size == "0" :
            return "0"
        if last_size.strip() == "" :
            last_size = "0"
        if str(current_size) != str(last_size) :
            waring(msg, waring_phones)
        return current_size

    def get_run_minute(sef, str) :
        min = 0
        if str.find("-") != -1 :
            bits = str.split("-")
            if len(bits) == 2 :
                min = 60 * 24 * int(bits[0])
                str = bits[1]
    
        if str.find(":") != -1 :
            bits = str.split(":")
            if len(bits) == 3 :
                min = min + (60 * int(bits[0]))
                min = min + int(bits[1])
            elif len(bits) == 2 :
                min = min + int(bits[0])
        return min




class CronParser :
    def __init__(self) :
        self.data = {}
        self.data["day"] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31]
        self.data["week"] = [0, 1, 2, 3, 4, 5, 6]
        self.data["hour"] = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23]
        self.data["month"] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
        self.data["minute"] = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59]

    def is_run(self, last_ran_time) :
        timestamp = time.time()
        if timestamp - last_ran_time < 59 :
            return False
        return True

    def parse(self, string) :
        bits = string.split(" ")

        if len(bits) != 5 :
            return False

        return self.__parse(bits)

    def __parse(self, bits) :
        minute, hour, day, month, week = bits

        if self.__analysis(month, "month") != True :
            return False

        if self.__analysis(week, "week") != True :
            return False

        if self.__analysis(day, "day") != True :
            return False

        if self.__analysis(hour, "hour") != True :
            return False

        if self.__analysis(minute, "minute") != True :
            return False

        return True

    def __analysis(self, str, date_type) :
        step = 0 
        unit = []
        result_unit = []
        pos = str.find('/')
        if pos != -1 :
            str, step = str.split("/")
            if str != "*" :
                return False

        if str == "*" :
            try :
                unit = self.data[date_type]
            except KeyError, e:
                return False

        else :
            sections = []
            pos = str.find(",")
            if pos != -1 :
                sections = str.split(",")
            else :
                sections.append(str)

            if len(sections) > 0 :
                for v in sections :
                    pos = v.find("-")
                    if pos != -1 :
                        start = 0
                        end = 0
                        start, end = v.split("-")
                        start = int(start)
                        end = int(end)
                        if start > end :
                            end, start = start, end

                        for i in range(start, end):
                            unit.append(i)
                    else :
                        unit.append(int(v))
        step = int(step)
        unit = list(set(unit))
        if step > 0 :
            if len(unit) > 0 :
                i = 0
                for v in unit :
                    if i % step == 0 :
                         result_unit.append(int(v))
                    i = i + 1
        else :
            result_unit = unit
        now_date = self.__get_now_number(date_type)
        retval = False
        for v in result_unit :
            if v == now_date :
                retval = True
        return retval


    def __get_now_number(self, date_type, timestamp = 0) :
        if timestamp == 0 :
            timestamp = time.time()

        if date_type == "month" :
            format = "%m"
        elif date_type == "week" :
            format = "%w"
        elif date_type == "day" :
            format = "%d"
        elif date_type == "hour" :
            format = "%H"
        elif date_type == "minute" :
            format = "%M"
        else :
            return -1

        return int(time.strftime(format, time.localtime(timestamp)))
            
if __name__ == "__main__" :
    config = ConfigParser.ConfigParser()
    syscmd = System()

    config_file = "/etc/scron.conf"
    waring_config_file = "/etc/scron_waring.conf"
    env = "pro"
    try :
        opts, args = getopt.getopt(sys.argv[1:], "hc:e:w:")
    except getopt.GetoptError, e :
        print e
        sys.exit(-1)

    for op, value in opts :
        if op == "-c" :
            config_file = value
        elif op == "-e" :
            env = value
        elif op == "-w" :
            waring_config_file = value
        elif op == "-h" :
            usage()
            sys.exit(-2)
    try :
        config.readfp(open(waring_config_file))
    except IOError, e :
        info("scron报警配置文件出错")
        sys.exit(-6)

    try :
        phones = config.get(env, "phones")
    except ConfigParser.NoSectionError, e :
        info("scron读取报警配置环境出错")
        sys.exit(-7)
    except ConfigParser.NoOptionError, e :
        info("scron读取报警配置有误")
        sys.exit(-8)

    try :
        config.readfp(open(config_file))
    except IOError, e :
        print e
        waring("scron读取配置文件：%s出错 来自主机: %s, ip: %s" % (config_file, hostname, localip), phones)
        sys.exit(-4)

    try :
        db_host = config.get(env, "db_host")
        db_name = config.get(env, "db_name")
        db_user = config.get(env, "db_user")
        db_pass = config.get(env, "db_pass")
        db_port = config.get(env, "db_port")
        log_path = config.get(env, "log_path")
        cron_table_name = config.get(env, "cron_table_name")
        host_table_name = config.get(env, "host_table_name")
        iskill = config.get(env, "iskill")

    except ConfigParser.NoSectionError, e :
        print e
        waring("scron读取配置环境: %s出错 来自主机: %s, ip: %s" % (env, hostname, localip), phones)
        sys.exit(-5)
    except ConfigParser.NoOptionError, e :
        print e
        waring("scron读取配置有误 来自主机: %s, ip: %s" % (hostname, localip), phones)
        sys.exit(-6)

    syscmd.mkdir(log_path)
    date_path = time.strftime("%Y/%m/%d")
    syscmd.mkdir( "/".join([log_path, date_path]) )

    ## open db
    try :
        conn = MySQLdb.connect(host=db_host, user=db_user, passwd=db_pass, db=db_name, port=int(db_port), cursorclass=MySQLdb.cursors.DictCursor)
        cursor=conn.cursor()
        sql = "SET NAMES UTF8"
        cursor.execute(sql)

        sql = "".join(["SELECT * FROM `%s` %s, `%s` %s" % (cron_table_name, "c", host_table_name, "h"), " WHERE c.host=h.host and c.host=%s"])
        cursor.execute(sql, [localip])

        for row in cursor.fetchall() :
            print sql
            cron = CronParser()
            user = row["user"]

            if user != None :
                user = user.strip()
            else :
                user = "root"

            command = row["command"].strip()
            params = row["params"].strip()
            print params
            timeout = int(row["timeout"])
            cron_id = row["cronId"];
            task = row["task"];
            active = int(row["active"])
            host_name = row["host_name"]
            last_size = str(row["errorLogUpdatedSize"])
            receiverPhone = row["receiverPhone"]
            if receiverPhone != None :
                receiverPhone = str(receiverPhone).strip()
            else :
                receiverPhone = ""
                
            #check timeout
            if timeout != 0 :
                msg = "taskId %s %s 在 %s 上运行时间已经超过了设置时间 %s 分钟。" % (cron_id, task, host_name, timeout)
                syscmd.check_timeout(command, params, msg, iskill, receiverPhone)
            if active ==  1 and cron.parse(row["mhdmd"]) == True and cron.is_run(float(row["runAt"])) == True :

                if row["logFile"].strip() == "" :
                    log_file = "".join([str(row["cronId"]), ".log"]);
                else :
                    log_file = row["logFile"];

                main_file, ext_file = os.path.splitext(log_file)

                log_error_file = "".join([main_file, "_error", ext_file])

                num = syscmd.process_num(command, params)
                print command
                print num
                print row["process"]
                if num >= int(row["process"]) :
                    continue

                log = "/".join([log_path, date_path, log_file]);
                error_log = "/".join([log_path, date_path, log_error_file]);
                print "run"
                syscmd.run(user, command, params, log, error_log)

                error_log_msg = "taskId %s %s 在 %s 上运行的错误日志递增。请及时查看最新错误日志。" % (cron_id, task, host_name)
                current_size = syscmd.check_error_log(error_log_msg, error_log, last_size, receiverPhone)
                if (current_size != last_size) :
                    sql = "".join(["UPDATE `%s` SET " % (cron_table_name, ), "errorLogUpdatedSize = %s WHERE cronId=%s"])
                    cursor.execute(sql, [str(current_size), row["cronId"]])
                    conn.commit()

                sql = "".join(["UPDATE `%s` SET " % (cron_table_name, ), "runAt = %s WHERE cronId=%s"])
                cursor.execute(sql, [str(time.time()), row["cronId"]])
                conn.commit()
    except MySQLdb.Error,e:
        print "Mysql Error %d %s" % (e.args[0], e.args[1])
        msg = "scron连接MySQL出错从主机：%s IP：%s 连接到 %s" % (hostname, localip, db_host)
        waring(msg, phones)

    try :
        cursor.close()
        conn.close()
    except NameError, e:
        print e
