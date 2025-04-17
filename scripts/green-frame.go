//go:build ignoretest
// +build ignoretest

package main

import (
  "bytes"
  "fmt"
  "net/http"
)

func sendPostRequest(data string) error {
  url := "http://localhost:17000/"
  req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "text/plain")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  return nil
}

func main() {
  if err := sendPostRequest("white"); err != nil {
    fmt.Printf("Ошибка команды 'white': %v\n", err)
    return
  }

  if err := sendPostRequest("bgrect 0.25 0.25 0.75 0.75"); err != nil {
    fmt.Printf("Ошибка команды 'bgrect': %v\n", err)
    return
  }

  if err := sendPostRequest("green"); err != nil {
    fmt.Printf("Ошибка команды 'green': %v\n", err)
    return
  }

  if err := sendPostRequest("figure 0.6 0.6"); err != nil {
    fmt.Printf("Ошибка команды 'figure': %v\n", err)
    return
  }

  if err := sendPostRequest("update"); err != nil {
    fmt.Printf("Ошибка команды 'update': %v\n", err)
    return
  }
}
