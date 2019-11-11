package queue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func getData(url string) error {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return err
	}
	client := http.DefaultClient
	client.Timeout = time.Duration(timeout)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func interfaceToJSON(input interface{}) string {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, strings.Replace(err.Error(), `"`, `\"`, -1))
	}
	return string(data)
}
