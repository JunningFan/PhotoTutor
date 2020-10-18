package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()
	resp, _ := client.EnableTrace().R().Get("http://localhost:9200/picture/_search")

	fmt.Println(resp)

}
