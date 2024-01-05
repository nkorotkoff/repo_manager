package admin

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"repo_manager/internal/config"
	"repo_manager/internal/structs"
	"repo_manager/internal/telegram_bot"
)

func Init(config *config.Config) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", index)
	r.Get("/get-data", getData)
	r.Post("/save-data", saveData)
	err := http.ListenAndServe(config.AdminHost+":"+config.AdminPort, r)
	if err != nil {
		fmt.Print("Error to serve router")
		return
	}

}

func index(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("../internal/admin/views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getData(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading data.json: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func saveData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var formData []structs.DevEnvironment
	if err := json.Unmarshal(body, &formData); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	file, err := os.Create("data.json")
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encodedData, err := json.MarshalIndent(formData, "", "    ")
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	_, err = file.Write(encodedData)
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	telegram_bot.ParseDevEnvironments()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message": "Data saved successfully"}`)
}
