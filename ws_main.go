package main

import (
	"bufio"
	"io/ioutil"
	"net/smtp"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

const oldFile string = "./old.txt"

const newFile string = "./new.txt"

func main() {
	old := comp(oldFile)
	write(newFile)
	new := comp(newFile)
	if old != new {
		sendGmail(new)
	}
	write(oldFile)
}

func write(file string) {
	//jsonファイルの読み込み
	bytes, jsonerr := ioutil.ReadFile("./context.json")
	if jsonerr != nil {
		panic(jsonerr)
	}
	// []byte型からjson型へ変換
	json, _ := simplejson.NewJson(bytes)
	doc, err := goquery.NewDocument(json.Get("url").MustString())
	if err != nil {
		panic(err)
	}
	res, _ := doc.Find("textarea").Html()
	err2 := ioutil.WriteFile(file, []byte(res), 0664)
	if err2 != nil {
		panic(err2)
	}
	//fmt.Println(res)
}

func comp(filePath string) (rs string) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	//戻り値用のスライスを宣言
	var rsSlice []string
	var cnt int = 0
	for scanner.Scan() {
		if cnt > 0 {
			break
		}
		scText := scanner.Text()
		rsSlice = append(rsSlice, scText)
		cnt++
	}
	defer file.Close()
	if len(rsSlice) > 0 {
		return rsSlice[0]
	}
	return ""
}

func sendGmail(message string) {
	//jsonファイルの読み込み
	bytes, err := ioutil.ReadFile("./context.json")
	if err != nil {
		panic(err)
	}
	// []byte型からjson型へ変換
	json, _ := simplejson.NewJson(bytes)
	auth := smtp.PlainAuth(
		"",
		json.Get("address").MustString(),
		json.Get("passwd").MustString(),
		"smtp.gmail.com",
	)

	err2 := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		json.Get("address").MustString(),
		[]string{"trigger@recipe.ifttt.com"},
		[]byte(
			"To:"+json.Get("address").MustString()+"\r\n"+
				"Subject:message\r\n"+
				"\r\n"+
				message),
	)
	if err2 != nil {
		panic(err)
	}
}
