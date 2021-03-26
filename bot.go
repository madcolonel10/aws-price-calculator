package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

type App struct {
	Router *router.Router
}

type HttpError struct {
	Msg        string
	StatusCode int
	Body       string
}

func (h HttpError) Error() string {
	return fmt.Sprintf("Err while Making Rest Call, Status: %s, StatusCode: %d, Msg: %s", h.Msg, h.StatusCode, h.Body)
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

	type Response struct {
		Data struct {
			MessageId string `json:"id"`
			RoomId    string `json:"roomId"`
		} `json:"data"`
	}

	var response Response
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Printf("error while trying to unmarshal webhook payload: %s\n", err)
		return
	}

	msgId := response.Data.MessageId
	data, err := getMessge(msgId)
	if err != nil {
		fmt.Printf("error in http:%s\n", err)
		return
	}
	fmt.Printf("data from message api: %s\n", data)
	instruction := getInstruction(data)
	fmt.Printf("instruction for bot is: %s\n", instruction)

	instructionTokens := strings.Split(instruction, " ")
	fmt.Printf("instruction tokens: %v\n", instructionTokens)

	instructionType := strings.TrimSpace(instructionTokens[0])

	switch instructionType {
	case "get":
		fmt.Println("get instuction")

	default:
		fmt.Println("invalid instruction")
		return
	}
}

func getEstimates(params []string) {
	//0 - serviceName
	//1 - env for which we will estimate the price
	//2 - time window 1 quarter is default
	defaultValues := []string{"jenkins", "dev", "1"}
	copy(defaultValues, params[1:]) //starting from index 1 as index 0 will be "estimate" keyword
	serviceName := defaultValues[0]
	env := defaultValues[1]
	nQuarters := defaultValues[2]

	fmt.Printf("serviceName: %s env: %s nQuarters: %s\n", serviceName, env, nQuarters)

	data, err := downloadInfraDataFromS3()
	if err != nil {
		return
	}
	fmt.Printf("data:%s\n", data)

	type HostConfiguration struct {
		VCpu   int
		Memory int
	}

	//service will be jenkins, bitbucket..
	//each service needs a bunch of hosts to run
	type CodeService struct {
		Hosts []HostConfiguration
	}

	var m map[string]map[string]CodeService = make(map[string]map[string]CodeService)
	err = yaml.Unmarshal([]byte(data), &m)

	fmt.Println(m["dev"]["jenkins"].Hosts[0].Memory)
}

func getInstruction(messageResponse string) string {
	var response struct {
		Text string `json:"text"`
	}
	err := json.Unmarshal([]byte(messageResponse), &response)
	if err != nil {
		fmt.Printf("error while trying to unmarshal instruction for bot: %s\n", err)
		return ""
	}
	str := "cops "
	return response.Text[len(str):]
}

func publishMessage(message string, roomId string) (string, error) {

	var body struct {
		RoomId string `json:"roomId"`
		Text   string `json:"text"`
	}

	body.RoomId = roomId
	body.Text = message

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	fmt.Printf("message to be published:%s roomId:%s\n", message, roomId)
	client := &http.Client{}
	url := "https://webexapis.com/v1/messages/"
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(bodyJson)))

	req.Header.Set("Authorization", "Bearer "+os.Getenv("BOT_ACCESS_TOKEN"))
	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err while making http call: %s", err)
		return "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", HttpError{resp.Status, resp.StatusCode, string(data)}
	}

	return string(data), nil

}

func getMessge(messageId string) (string, error) {
	client := &http.Client{}
	url := "https://webexapis.com/v1/messages/" + messageId
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", "Bearer "+os.Getenv("BOT_ACCESS_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err while making http call: %s", err)
		return "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", HttpError{resp.Status, resp.StatusCode, string(data)}
	}
	return string(data), nil
}

func InitializeApp() *App {
	fmt.Println("Initializing App")
	app := &App{Router: router.New()}
	app.setupRoutes()
	return app
}

func downloadInfraDataFromS3() (string, error) {
	url := os.Getenv("S3_URL")
	fmt.Printf("s3 url:%s\n", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err while making http call: %s", err)
		return "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", HttpError{resp.Status, resp.StatusCode, string(data)}
	}
	return string(data), nil
}

func main() {
	//app := InitializeApp()
	//app.Run()
	//getEstimates([]string{"estimate"})
	_, err := publishMessage("testMessage", "Y2lzY29zcGFyazovL3VzL1JPT00vNDA1NzBhZjAtZDRjYS0xMWVhLTk2YzctYjExZmFhNTI1Mjcx")
	if err != nil {
		fmt.Println(err)
	}
}
