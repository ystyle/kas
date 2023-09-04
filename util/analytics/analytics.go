package analytics

import (
	fmt "fmt"
	"github.com/ystyle/google-analytics"
	"math/rand"
	"os"
	"runtime"
)

var (
	secret      string
	measurement string
	version     string
)

func Analytics(clientID, format string) {
	if secret == "" || measurement == "" {
		return
	}
	analytics.SetKeys(secret, measurement) // // required
	payload := analytics.Payload{
		ClientID: clientID, // required
		UserID:   getClientID(),
		Events: []analytics.Event{
			{
				Name: "kas", // required
				Params: map[string]interface{}{
					"os":      runtime.GOOS,
					"arch":    runtime.GOARCH,
					"format":  format,
					"version": version,
				},
			},
		},
	}
	analytics.Send(payload)
}

func getClientID() string {
	clientID := fmt.Sprintf("%d", rand.Uint32())
	config, err := os.UserConfigDir()
	if err != nil {
		return clientID
	}
	filepath := fmt.Sprintf("%s/kas/config", config)
	if exist, _ := isExists(filepath); exist {
		bs, err := os.ReadFile(filepath)
		if err != nil {
			return clientID
		}
		clientID = string(bs)
	} else {
		err := os.MkdirAll(fmt.Sprintf("%s/kas", config), 0700)
		if err != nil {
			return clientID
		}
		_ = os.WriteFile(filepath, []byte(clientID), 0700)
	}
	return clientID
}

func isExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
