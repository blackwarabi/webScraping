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
	log.Print("処理開始")
	old := readFile(oldFile)
	writeFile(newFile)
	new := readFile(newFile)
	if old != new {
		sendGmail(new)
	}
	writeFile(oldFile)
	log.Print("処理終了")
}

//現在の更新情報を引数のファイルへ書き込み
func writeFile(file string) {
	log.Print("writeFileの処理開始")
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
	log.Print(res)
	log.Print("writeFileの処理終了")
}

//引数のファイルを読み込み、先頭の文字列を返す
func readFile(filePath string) (rs string) {
	log.Print("readFileの処理開始")
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	scText := scanner.Text()
	defer file.Close()
	log.Print("readFileの処理終了")
	return scText
}

//メール送信
func sendGmail(message string) {
	log.Print("sendGmailの処理開始")
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
	log.Print("sendGmailの処理終了")
}

//設定ファイルの読み込み
func loadConfigFile() []byte {
	log.Print("loadConfigFileの処理開始")
	//jsonファイルの読み込み
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("loadConfigFileの処理終了")
	return bytes
}
