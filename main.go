package main

import (
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"
    "io/ioutil"
    "github.com/tidwall/gjson"
    "os"
    "bufio"
)

var jar, _ = cookiejar.New(nil)
var client = http.Client{Jar:jar}

type Channel struct {
    Media_name string `json:"media_name"`
    Issues_id string `json:"issues_id"`
    Play_link string `json:"play_link"`
}

func setPlayLink(playLink string, channel *Channel) {
    channel.Play_link = playLink
}

type Channels []Channel

func main () {
    login()
    channels := getChannels()
    fmt.Println(channels)

}

func generatePlaylist(path string, channels Channels) {
    file, _ := os.Create(path)
    defer file.Close()
    writer := bufio.NewWriter(file)
    writer.WriteString("#EXTM3U\n")


    writer.Flush()
}

func getChannels() Channels {
    channels := Channels{}
    request, _ := http.NewRequest("POST", "http://www.neterra.tv/content/live", nil)
    response, _ := client.Do(request)
    bodyBytes, _ := ioutil.ReadAll(response.Body)
    contentJson := gjson.GetBytes(bodyBytes, "prods")

    for _, result := range contentJson.Array() {
        channel := Channel{}
        gjson.Unmarshal([]byte(gjson.Get(result.Raw, "0").Raw), &channel)

        // Get play link
        fmt.Println("Getting URL for " + channel.Media_name)
        payload := url.Values{}
        payload.Set("issue_id", channel.Issues_id)
        payload.Set("quality", "0")
        payload.Set("type", "live")
        request, _ = http.NewRequest("POST", "http://www.neterra.tv/content/get_stream",
            strings.NewReader(payload.Encode()))
        request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

        response, _ := client.Do(request)
        bodyBytes, _ := ioutil.ReadAll(response.Body)
        playLink := gjson.GetBytes(bodyBytes, "play_link")
        setPlayLink(playLink.Raw, &channel)

        channels = append(channels, channel)
    }

    return channels
}

func login() bool {
    // Init
    request, _ := http.NewRequest("GET", "http://www.neterra.tv/user/login_page", nil)
    client.Do(request)

    // Login POST
    //TODO: Read from Config
    payload := url.Values{}
    payload.Set("login_username", "")
    payload.Set("login_password", "")
    payload.Set("login", "1")
    payload.Set("login_type", "1")

    request, _ = http.NewRequest("POST", "http://www.neterra.tv/user/login_page",
        strings.NewReader(payload.Encode()))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    response, _ :=client.Do(request)
    bodyBytes, _ := ioutil.ReadAll(response.Body)
    body := string(bodyBytes)

    return strings.Contains(body, "var LOGGED = '1'")
}