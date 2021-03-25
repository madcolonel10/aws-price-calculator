package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type App struct {
	Router *router.Router
}

func (app *App) Run() {
	fasthttp.ListenAndServe(":5000", app.Router.Handler)
}

func (app *App) setupRoutes() {
	v1BaseEndpoint := "/api/v1/bot"
	app.Router.GET(v1BaseEndpoint, echoEndpoint)
	app.Router.POST(v1BaseEndpoint+"/message", botLogic)
}

func echoEndpoint(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("v1 version of bot api")
}

func botLogic(ctx *fasthttp.RequestCtx) {
	fmt.Println("bot logic starts here")
	body := ctx.Request.Body()
	fmt.Println(string(body))
	msgId := "Y2lzY29zcGFyazovL3VzL01FU1NBR0UvODA0ZTI4NDAtOGRhOC0xMWViLTlkYTEtMzMxZGY5ZmRkZmJm"
	data, err := getMessge(msgId)
	if err != nil {
		fmt.Printf("error in http:%s\n", err)
	}
	fmt.Printf("data from message api: %s\n", data)
}

func getMessge(messageId string) (string, error) {
	client := &http.Client{}
	url := "https://webexapis.com/v1/messages/" + messageId
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err while making http call: %s", err)
		return "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", err
	}
	return string(data), nil
}

func InitializeApp() *App {
	fmt.Println("Initializing App")
	app := &App{Router: router.New()}
	app.setupRoutes()
	return app
}

func main() {
	app := InitializeApp()
	fmt.Println("yo yo")
	app.Run()
}
