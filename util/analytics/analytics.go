package analytics

import (
	"fmt"
	"github.com/ystyle/google-analytics"
	"math/rand"
	"runtime"
	"time"
)

var (
	secret      string
	measurement string
	version     string
)

func Analytics(client, format string) {
	if secret == "" || measurement == "" {
		return
	}
	t := time.Now().Unix()
	analytics.SetKeys(secret, measurement) // // required
	payload := analytics.Payload{
		ClientID: fmt.Sprintf("%d.%d", rand.Int31(), t), // required
		UserID:   client,
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
