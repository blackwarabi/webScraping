package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

//更新前ファイル
const oldFile string = "./outFile/old.txt"

//更新後ファイル
const newFile string = "./outFile/new.txt"

//メイン処理
func main() {
	old := readFile(oldFile)
	write(newFile)
	new := readFile(newFile)
	if old != new {
		sendGmail(new)
	}
	write(oldFile)
}

//現在の更新情報を引数のファイルへ書き込み
func write(file string) {
	//設定ファイルの読み込み
	bytes := loadConfigFile()
	// []byte型からjson型へ変換
	json, _ := simplejson.NewJson(bytes)

	doc, err := goquery.NewDocument(json.Get("url").MustString())
	if err != nil {
		log.Fatal(err)
	}
	res, _ := doc.Find("textarea").Html()
	err2 := ioutil.WriteFile(file, []byte(res), 0664)
	if err2 != nil {
		log.Fatal(err2)
	}
	//fmt.Println(res)
}

//引数のファイルを読み込み、先頭の文字列を返す
func readFile(filePath string) (rs string) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	scText := scanner.Text()
	defer file.Close()
	return scText
}

//メール送信
func sendGmail(message string) {
	//設定ファイルの読み込み
	bytes := loadConfigFile()
	// []byte型からjson型へ変換
	json, _ := simplejson.NewJson(bytes)
	auth := smtp.PlainAuth(
		"",
		json.Get("address").MustString(),
		json.Get("passwd").MustString(),
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
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
	if err != nil {
		log.Fatal(err)
	}
}

//設定ファイルの読み込み
func loadConfigFile() []byte {
	//jsonファイルの読み込み
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
