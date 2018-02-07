package main

import "io/ioutil"
import "net/http"
import "strings"
import "errors"
import "time"

type Client struct{
    client http.Client

}

func NewClient() *Client{
    return &Client{http.Client{Timeout: 5*time.Second}}
}

func (c *Client) get_if_content_type_is(url string, content_type... string) ([]byte, int, error){
    response, err := c.client.Get(url)
    if err!=nil {
        return nil, -1, err
    }
    defer response.Body.Close()

    content_type_index:=int(-1)
    if func()bool{
        for index,ct:=range content_type{
            if strings.Contains(response.Header.Get("Content-Type"), ct){
                content_type_index=index
                return false
            }
        }
        return true
    }(){
        return nil, -1, errors.New("Content-Type does not contain mime type")
    }

    content, err:=ioutil.ReadAll(response.Body)
    if err!=nil {
        return nil, -1, err
    }

    return content, content_type_index, nil
}