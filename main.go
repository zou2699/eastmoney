package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"regexp"
)

const gzzfRegex = `id="gz_gszzl">(.*?)<`
const nameRegex = `<div style="float: left">(.*?)<span>`
const idRegex = `</span><span class="ui-num">(.*?)</span></div>`
const dateRegex = `id="gz_gztime">[(?](.*?)[)?]</span>`

type Msg struct {
	Name   string
	Id     string
	Gzzf   string
	GzDate string
}

type Config struct {
	Id []string `yaml:"id"`
}

func main() {

	var conf Config
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("%+v", conf)

	for _, id := range conf.Id {
		bUrl := "http://fund.eastmoney.com/" + id + ".html"
		msg := parse(bUrl)
		fmt.Printf("%+v\n", msg)
	}

	//bUrl := "http://fund.eastmoney.com/110022.html"
	//msg := parse(bUrl)
	//fmt.Printf("%+v\n", msg)

}

func parse(bUrl string) Msg {
	result := Msg{}

	client := http.Client{}
	request, err := http.NewRequest("GET", bUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Add("Accept-Charset", "UTF-8,*;q=0.5")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:60.0) Gecko/20100101 Firefox/60.0")
	request.Header.Add("referer", "http://fund.eastmoney.com")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(string(respBody))

	s := reg(nameRegex, respBody)
	result.Name = s

	s = reg(idRegex, respBody)
	result.Id = s

	s = reg(dateRegex, respBody)
	result.GzDate = s

	s = reg(gzzfRegex, respBody)
	result.Gzzf = s

	return result

}

func reg(regexString string, content []byte) string {
	Reg := regexp.MustCompile(regexString)
	match := Reg.FindAllSubmatch(content, -1)
	for _, m := range match {
		//fmt.Println("基金ID: ", string(m[1]))
		return string(m[1])
	}
	return ""
}
