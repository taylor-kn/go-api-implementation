package main
import (
	"fmt"
	"net/http"
	"log"
	"github.com/go-chi/chi"
	"encoding/json"
	"strconv"
)

type Service struct {
	Service_id string
	Service_name string
	Alerts []Alert
}

type Alert struct {
	AlertID string `json:"alert_id"`
	Model string `json:"model"`
	AlertType string `json:"alert_type"`
	AlertTS string `json:"alert_ts"`
	Severity string `json:"severity"`
	TeamSlack string `json:"team_slack"`
}

type Request struct {
	Alert_id string `json:"alert_id"`
	Service_id string `json:"service_id"`
	Service_name string `json:"service_name"`
	Model string `json:"model"`
	AlertType string `json:"alert_type"`
	AlertTS string `json:"alert_ts"`
	Severity string `json:"severity"`
	TeamSlack string `json:"team_slack"`
}

type Response struct {
	Alert_id string `json:"alert_id"`
	Error string `json:"error"`
}

var Services = make(map[string]Service)

func main() {
	// Router setup
	r := chi.NewRouter()
	// Route requests
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from home"))
	})
	r.Post("/alerts", WriteAlert)
	r.Get("/alerts", ReadAlerts)
	// Server start
	srv := &http.Server{
	Addr: fmt.Sprintf(":%s", "8080"),
	Handler: r,
	}
	log.Println("Server started...")
	if err := srv.ListenAndServe(); err != nil {
	log.Fatal(fmt.Sprintf("%+v", err))
	}
}

func WriteAlert(w http.ResponseWriter, r *http.Request) {
    var req Request

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	newAlert := Alert {
		AlertID: req.Alert_id,
		Model: req.Model,
		AlertType: req.AlertType,
		AlertTS: req.AlertTS,
		Severity: req.Severity,
		TeamSlack: req.TeamSlack,
	}

	service, exists := Services[req.Service_id]
    if exists {
		service.Alerts = append(service.Alerts, newAlert)
		Services[req.Service_id] = service
    } else {
        newService := Service{
            Service_id: req.Service_id,
            Service_name: req.Service_name,
			Alerts: []Alert{newAlert},
        }
        Services[req.Service_id] = newService
    }

    res := Response{req.Alert_id, ""}
    jsonResponse, _ := json.Marshal(res)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func ReadAlerts(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var queryParams = make(map[string]string)
	for k, v := range values {
		queryParams[k] = v[0]
	}

	//Filter by start and end time query parameters
	startTime, _ := strconv.Atoi(queryParams["start_ts"])
	endTime, _ := strconv.Atoi(queryParams["end_ts"])

	filteredAlerts := []Alert{}

	for _, alert := range Services[queryParams["service_id"]].Alerts {
		alertTS, _ := strconv.Atoi(alert.AlertTS)
		if(alertTS >= startTime && alertTS <= endTime) {
			filteredAlerts = append(filteredAlerts, alert)
		}
	}

	filteredService := Service{
		Service_id: queryParams["service_id"],
		Service_name: Services[queryParams["service_id"]].Service_name,
		Alerts: filteredAlerts,
	}

    jsonResponse, _ := json.Marshal(filteredService)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}