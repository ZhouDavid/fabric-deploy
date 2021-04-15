package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	//服务器端口号
	serverPort string
	//钉钉告警url
	dingdingUrl string
)

//钉钉告警模板
const template = `{
     "msgtype": "text",
     "text": {
         "content": "报警%s"
     },
     "at": {
         "isAtAll": false
     }
 }`

func init() {

	// flag.StringVar(&serverPort, "p", "9201", "server port")
	// flag.StringVar(&dingdingUrl, "u", "", "alarm webhook")
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
	//读取skywalking告警消息
	body, _ := ioutil.ReadAll(r.Body)
	//拼接钉钉所需的消息格式
	msg := strings.NewReader(fmt.Sprintf(template, body))
	//发送钉钉消息
	res, err := http.Post(dingdingUrl, "application/json", msg)

	if err != nil || res.StatusCode != 200 {
		log.Print("send alram msg to dingding fatal.", err, res.StatusCode)
	}

}
func hmacSha256(ti, secret string) string {

	//
	// ti := "1618219217645"
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(ti + "\n" + secret))
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(sha, ti)
	return url.QueryEscape(sha)
}
func main() {
	t := time.Now().UnixNano() / 1e6
	ti := strconv.Itoa(int(t))
	// flag.Parse()
	sign := hmacSha256(ti, "SECd6b811498de0dea9d9cf5e3dce57dba7b2bf23f63b6f08dd949cdb09270f1f17")
	// if dingdingUrl == "" {
	fmt.Println(sign)

	dingdingUrl = "https://oapi.dingtalk.com/robot/send?access_token=8e6d7291483c60e7227cb992fc9f880f280dd7911b01fd12d68e2b8a898379ad&timestamp=" + ti + "&sign=" + sign + ""
	// log.Fatal("dingdingUrl cannot be empty usage -h get help. ")
	// }
	msg := strings.NewReader(fmt.Sprintf(template, "服务启动有问题"))
	res, err := http.Post(dingdingUrl, "application/json", msg)
	defer res.Body.Close()
	if err != nil || res.StatusCode != 200 {
		log.Print("send alram msg to dingding fatal.", err, res.StatusCode)
	}
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	// http.HandleFunc("/alarm", sendMsg)

	// //启动web服务器
	// if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
	// 	log.Fatal("server start fatal.", err)
	// }

}
