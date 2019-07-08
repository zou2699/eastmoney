package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

const gzzfRegex = `id="gz_gszzl">(.*?)<`
const nameRegex = `<div style="float: left">(.*?)<span>`
const idRegex = `</span><span class="ui-num">(.*?)</span></div>`
const dateRegex = `id="gz_gztime">[(?](.*?)[)?]</span>`

type Msg struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	GzZf   string `json:"gz_zf"`
	GzDate string `json:"gz_date"`
}

type Config struct {
	Id []string `yaml:"id"`
}

func main() {
	router := gin.Default()

	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	router.Use(cors.New(config))

	r := router.Group("/")

	{
		r.GET("/", func(c *gin.Context) {
			var conf Config
			var msgList []Msg
			var wg sync.WaitGroup
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
				wg.Add(1)
				go func(id string) {
					bUrl := "http://fund.eastmoney.com/" + id + ".html"
					msg := parse(bUrl)
					msgList = append(msgList, msg)
					wg.Done()
				}(id)
			}
			wg.Wait()
			c.JSON(http.StatusOK, msgList)
		})

		r.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			bUrl := "http://fund.eastmoney.com/" + id + ".html"
			msg := parse(bUrl)
			if msg.Id == "" {
				msg.Name = "Not Found"
				c.JSON(http.StatusOK, msg)
				return
			}
			c.JSON(http.StatusOK, msg)
		})
	}

	router.Run(":8000")
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
		return result
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return result
	}

	//fmt.Println(string(respBody))

	s := reg(nameRegex, respBody)
	result.Name = s

	s = reg(idRegex, respBody)
	result.Id = s

	s = reg(dateRegex, respBody)
	result.GzDate = s

	s = reg(gzzfRegex, respBody)
	result.GzZf = s

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
