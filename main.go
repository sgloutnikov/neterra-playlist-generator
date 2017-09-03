package main

import (
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strings"
    "io/ioutil"
)

var jar, _ = cookiejar.New(nil)
var client = http.Client{Jar:jar}

func main () {
    fmt.Println("Hello world.")


    // Make Json models

    // Login
    // Get Live Stream and save name + issue_id
    // Get stream with issue_id and save to file

}



func getChannels() {

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