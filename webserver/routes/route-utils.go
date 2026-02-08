package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/louie-jones-strong/go-shared/logger"
)

func RenderPage(w http.ResponseWriter, templateName string, data any) {

	t, err := template.ParseGlob("templates/components/*.jinja")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err = t.ParseFiles(fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

func UnmarshalRequestBody[T any](r *http.Request) (T, error) {
	var value T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return value, fmt.Errorf("Unable to read body: %v", err)
	}

	var rawMessage json.RawMessage
	rawMessage = body
	logger.Debug("RawMessage: %s\n", rawMessage)

	err = json.Unmarshal(body, &value)
	if err != nil {
		return value, fmt.Errorf("Invalid JSON body: %v", err)
	}

	return value, nil
}

func WriteJSONResponse[T any](w http.ResponseWriter, obj T) error {
	start := time.Now()

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode object as JSON and send it
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		return err
	}

	logger.Debug("Writing JSON Response took: %v", time.Since(start))
	return nil
}

type RouteInfo struct {
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

func GetRoutes(r *chi.Mux) ([]RouteInfo, error) {
	routes := make([]RouteInfo, 0)

	walkFunc := func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {

		routes = append(routes, RouteInfo{
			Method: method,
			Path:   route,
		})
		return nil
	}

	err := chi.Walk(r, walkFunc)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func PrintRoutes(r *chi.Mux) error {
	routes, err := GetRoutes(r)
	if err != nil {
		return err
	}

	for _, route := range routes {
		_, err := fmt.Printf("%s %s\n", route.Method, route.Path)
		if err != nil {
			return err
		}
	}

	return nil
}
