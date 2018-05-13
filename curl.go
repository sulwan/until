package until

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PostData(post_data interface{}, sinapay_url string) (result map[string]interface{}, err error) {

	excludeSlice := []string{"sign", "sign_type", "sign_version"}
	splice := QueryStr(post_data, true, excludeSlice)

	fmt.Println(splice)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(sinapay_url, "application/x-www-form-urlencoded;charset=utf-8", strings.NewReader(splice))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()
	ret, rerr := ioutil.ReadAll(res.Body)
	if rerr != nil {
		log.Fatal(err)
		return
	}
	encodeStr, err := url.QueryUnescape(string(ret))
	result = make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(encodeStr)))
	decoder.UseNumber() // 此处能够保证bigint的精度
	decoder.Decode(&result)
	return
}
