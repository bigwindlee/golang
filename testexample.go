package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
)

// Define a type that will store home page sizes:
type HomePageSize struct {
  URL string
  Size int
}

func main(){
  urls := []string{
    "http://www.taobao.com",
    "http://www.qq.com",
    "http://www.baidu.com",
    "http://www.jd.com",
  }

  results := make(chan HomePageSize)

  for _, url := range urls{
    // Notice that this is an unnamed function that is immediately invoked.
    // This is a common pattern with goroutines.
    go func(url string){
      res, err := http.Get(url)
      if err != nil{
        panic(err)
      }
      defer res.Body.Close()

      bs, err := ioutil.ReadAll(res.Body)
      if err != nil {
        panic(err)
      }

      results <- HomePageSize{
        URL: url,
        Size: len(bs),
      }
    }(url)
  }

  var biggest HomePageSize

  for range urls{
    result := <-results
    if result.Size > biggest.Size{
      biggest = result
    }
  }

  fmt.Println("The biggest home page:", biggest.URL, biggest.Size)
}
