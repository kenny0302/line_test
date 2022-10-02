package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	db "main/db"
	proto "main/proto"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	path, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic("Can Not Read File：" + err.Error())
	}
	router.POST("/callback", GetBotData)
	router.GET("/list", ListUser)
	router.POST("/push", PushMessage)
	return router
}

func GetBotData(c *gin.Context) {
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	secret := viper.GetString("line.secret")
	token := viper.GetString("line.token")
	bot, err := linebot.New(secret, token)
	if err != nil {
		log.Fatal(err)
		c.String(500, "Can Not Get Config!")
	}
	var data proto.User
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, "Not OK! "+err.Error())
		} else {
			c.String(500, "Not OK! "+err.Error())
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if len(event.Source.UserID) != 0 {
					client := &http.Client{}
					url := "https://api.line.me/v2/bot/profile/" + event.Source.UserID
					req, err := http.NewRequest("GET", url, nil)
					if err != nil {
						log.Println(err)
						c.String(500, err.Error())
					}
					req.Header.Add("Authorization", "Bearer "+token)
					resp, err := client.Do(req)
					if err != nil {
						log.Println(err)
						c.String(500, err.Error())
					}
					defer resp.Body.Close()
					sitemap, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatal(err)
						c.String(500, err.Error())
					}
					json.Unmarshal(sitemap, &data)
				}
				//reply same message
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
					c.String(500, err.Error())
				}
				data.Message = message.Text
			}
		}
	}
	jsondata, _ := json.Marshal(data)
	err = Set(host, port, jsondata)
	if err != nil {
		c.String(500, err.Error())
	}
	c.String(200, "Set OK!")
}

func Set(host, port string, data []byte) error {

	var Newdata proto.User
	json.Unmarshal(data, &Newdata)

	out, err := db.GetUser(host, port, "line", "user", bson.D{{"userid", Newdata.UserId}})
	if err != nil {
		return err
	}
	if len(out) <= 0 {
		err = db.SetUser(host, port, "line", "user", bson.M{"userid": Newdata.UserId, "displayname": Newdata.DisplayName, "pictureurl": Newdata.PictureUrl})
		if err != nil {
			return err
		}
	}
	time := time.Now().Unix()
	err = db.SetMessage(host, port, "line", "message", bson.M{"userid": Newdata.UserId, "message": Newdata.Message, "time": strconv.FormatInt(time, 10)})
	if err != nil {
		return err
	}
	return nil
}

func ListUser(c *gin.Context) {
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	out, err := db.GetUser(host, port, "line", "user", bson.D{{}})
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, out)
	}
}

func PushMessage(c *gin.Context) {
	secret := viper.GetString("line.secret")
	token := viper.GetString("line.token")
	userid := c.PostForm("UserId")
	bot, err := linebot.New(secret, token)
	if err != nil {
		log.Fatal(err)
		c.String(500, "Can Not Get Config!")
	}
	_, err = bot.PushMessage(userid, linebot.NewTextMessage("Testing:20221003")).Do()
	if err != nil {
		log.Print(err)
		c.String(500, err.Error())
	}
	c.String(200, "Spark OK!")
}
