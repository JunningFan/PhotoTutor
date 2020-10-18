package elsClient

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"phototutor/backend/util"
)


func PutElsObj(url string, v interface{})  {
	if jbytes, err:= json.Marshal(v); err!= nil {
		log.Printf("Els Marshal Err: %s\n",err.Error())
	} else if resp, err := http.Post(util.ELS_BASE + url, "application/json", bytes.NewReader(jbytes));  err != nil {
		log.Printf("Els Sync Err: %s\n",err.Error())
	} else if err := resp.Body.Close(); err!= nil  {
		log.Printf("Els Close Fp Err: %s\n",err.Error())
	}
}