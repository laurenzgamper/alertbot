package main

import (
	"encoding/json"
	"net/http"
	"io"
	"log"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
	"github.com/hako/durafmt"
)

func decodeWebhook(body io.ReadCloser) notify.WebhookMessage {
	var webhook notify.WebhookMessage

	decoder := json.NewDecoder(body)
	defer func() {
		err := body.Close();
		if err != nil {
			log.Print("can't close reponse body");
		}
	}()

	err := decoder.Decode(&webhook);
	if  err != nil {
		log.Print("failed to decode webhook message");
	}
	return webhook
}

func AlertToMessage(a types.Alert) string {
	var status, duration string
	switch a.Status() {
	case model.AlertFiring:
		status = fmt.Sprintf("🔥 *%s* 🔥", strings.ToUpper(string(a.Status())))
		duration = fmt.Sprintf("*Started*: %s ago", durafmt.Parse(time.Since(a.StartsAt)))
	case model.AlertResolved:
		status = fmt.Sprintf("*%s*", strings.ToUpper(string(a.Status())))
		duration = fmt.Sprintf(
			"*Ended*: %s ago\n*Duration*: %s",
			durafmt.Parse(time.Since(a.EndsAt)),
			durafmt.Parse(a.EndsAt.Sub(a.StartsAt)),
		)
	}

	return fmt.Sprintf(
		"%s\n*%s* (%s)\n%s\n%s\n",
		status,
		a.Labels["alertname"],
		a.Annotations["summary"],
		a.Annotations["description"],
		duration,
	)
}

// HandleWebhook returns a HandlerFunc that sends messages to bot via a channel
func HandleWebhook(messages chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("got message")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		webhook := decodeWebhook(r.Body)

		for _, webAlert := range webhook.Alerts {
			labels := make(map[model.LabelName]model.LabelValue)
			for k, v := range webAlert.Labels {
				labels[model.LabelName(k)] = model.LabelValue(v)
			}

			annotations := make(map[model.LabelName]model.LabelValue)
			for k, v := range webAlert.Annotations {
				annotations[model.LabelName(k)] = model.LabelValue(v)
			}

			alert := types.Alert{
				Alert: model.Alert{
					StartsAt:     webAlert.StartsAt,
					EndsAt:       webAlert.EndsAt,
					GeneratorURL: webAlert.GeneratorURL,
					Labels:       labels,
					Annotations:  annotations,
				},
			}

			messages <- AlertToMessage(alert) + "\n"
		}

		w.WriteHeader(http.StatusOK)
	}
}

func listenWebHook(address string, messages chan<- string ) {
	http.HandleFunc("/", HandleWebhook(messages))
	err := http.ListenAndServe(address, nil)
	log.Print(err)
}