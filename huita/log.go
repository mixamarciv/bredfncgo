package main

import (
	a "app_fnc"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	s "strings"
)

var log_file string

func InitLog() {

	log_path, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if runtime.GOOS == "windows" {
		log_path = s.Replace(log_path, "\\", "/", -1)
	} else { //if runtime.GOOS == "linux" {
		log_path = "/var/log"
	}

	timestr := a.CurTimeStrShort()
	log_path = log_path + "/log_nax_exporter/" + timestr[0:8]

	a.MkdirAll(log_path)
	log_file = log_path + "/" + a.CurTimeStrShort() + ".log"

	WriteLogln("start log")
	WriteLogln("log file: " + log_file)
}

func WriteLog(data string) {
	a.FileAppendStr(log_file, data)
}

func WriteLogln(data string) {
	a.FileAppendStr(log_file, a.CurTimeStr()+" "+s.TrimRight(data, "\n\r\t ")+"\n")
}

func WriteLogErr(info string, err error) {
	a.FileAppendStr(log_file, a.CurTimeStr()+" "+info+"\n"+a.ErrStr(err))
}

func WriteLogErrAndExit(info string, err error) {
	if err == nil {
		return
	}
	a.FileAppendStr(log_file, a.CurTimeStr()+" "+info+"\n"+a.ErrStr(err))
	panic(err)
	os.Exit(1)
}

func LogPrint(data string) {
	fmt.Println(data)
	WriteLogln(data)
}

func LogPrintErrAndExit(info string, err error) {
	if err == nil {
		return
	}
	fmt.Println(info)
	fmt.Printf("%+v", err)
	WriteLogErrAndExit(info, err)
}

func LogPrintErr(info string, err error) {
	if err == nil {
		return
	}
	fmt.Println(info)
	fmt.Printf("%+v", err)
	WriteLogErr(info, err)
}

func LogPrintAndExit(info string) {
	fmt.Println(info)
	WriteLogln(info)
	os.Exit(1)
}
