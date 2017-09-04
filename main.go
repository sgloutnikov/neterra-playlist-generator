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
    "bytes"
    "log"
    "github.com/spf13/viper"
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
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.ReadInConfig()

    if !login(viper.GetString("login.username"), viper.GetString("login.password")) {
        log.Fatal("Could not login. Check credentials.")
    }

    fmt.Println("Generating playlist in " + viper.GetString("playlist.path"))
    channels := getChannels()
    generatePlaylist(viper.GetString("playlist.path"), channels)
}

func generatePlaylist(path string, channels Channels) {
    file, _ := os.Create(path)
    defer file.Close()
    writer := bufio.NewWriter(file)
    writer.WriteString("#EXTM3U\n")

    // Write out available channels
    var buffer bytes.Buffer
    for _, channel := range channels {
        buffer.WriteString("#EXTINF:-1, ")
        buffer.WriteString(channel.Media_name)
        buffer.WriteString("\n")
        buffer.WriteString(channel.Play_link)
        buffer.WriteString("\n")

        writer.WriteString(buffer.String())
        buffer.Reset()
    }

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
        setPlayLink(playLink.Str, &channel)

        channels = append(channels, channel)
    }

    return channels
}

func login(user, pass string) bool {
    // Init
    request, _ := http.NewRequest("GET", "http://www.neterra.tv/user/login_page", nil)
    client.Do(request)

    // Login POST
    payload := url.Values{}
    payload.Set("login_username", user)
    payload.Set("login_password", pass)
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