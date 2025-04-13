package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"

	"vet_note/db"
)

func getPatientList(w http.ResponseWriter, r *http.Request) {
	patients, err := db.GetPatientList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(patients)
}

func handlePatient(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getPatient(w, r)
		return
	}

	if r.Method == http.MethodPost {
		updatePatient(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		deletePatient(w, r)
		return
	}
}

func getPatient(w http.ResponseWriter, r *http.Request) {
	// Extract patientId from URL: expected URL pattern /v1/patient/{patientId}
	id := strings.TrimPrefix(r.URL.Path, "/v1/patient/")
	patient, err := db.GetPatient(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(patient)
}

func updatePatient(w http.ResponseWriter, r *http.Request) {
	var viewPatient db.ViewPatient
	if err := json.NewDecoder(r.Body).Decode(&viewPatient); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("Invalid request body: %s", err))
		return
	}

	if _, ok := viewPatient.ID.Get(); ok {
		// Update existing patient
		err := db.UpdatePatient(viewPatient.AsSetter())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	} else {
		// Create new patient
		id, err := db.CreatePatient(viewPatient.AsSetter())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
		viewPatient.ID.Set(int32(id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(viewPatient)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updatePatients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	var viewPatients []db.ViewPatient
	if err := json.NewDecoder(r.Body).Decode(&viewPatients); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("Invalid request body: %s", err))
		return
	}

	// Update existing patient
	for _, viewPatient := range viewPatients {
		err := db.UpdatePatient(viewPatient.AsSetter())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func deletePatient(w http.ResponseWriter, r *http.Request) {
	// Extract patientId from URL: expected path format: /v1/patient/{patientId}
	id := strings.TrimPrefix(r.URL.Path, "/v1/patient/")
	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("parse error: %s", err))
		return
	}
	err = db.DeletePatient(int32(intID))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("db error: %s", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleProcedure(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getProcedure(w, r)
		return
	}
	if r.Method == http.MethodPost {
		updateProcedure(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		deleteProcedure(w, r)
		return
	}
}

func getProcedure(w http.ResponseWriter, r *http.Request) {
	// Extract procedureId from URL: expected path format: /v1/procedure/{procedureId}
	id := strings.TrimPrefix(r.URL.Path, "/v1/procedure/")
	procedure, err := db.GetProcedure(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(db.FmtError("Database error: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(procedure)
}

func updateProcedure(w http.ResponseWriter, r *http.Request) {
	var viewProcedure db.ViewProcedure

	if err := json.NewDecoder(r.Body).Decode(&viewProcedure); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("Invalid request body: %s", err))
		return
	}

	// Log the updated procedure information.
	slog.Debug("updateProcedure", "procedure", viewProcedure)

	if _, ok := viewProcedure.ID.Get(); ok {
		err := db.UpdateProcedure(viewProcedure.AsSetter())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	} else {
		err := db.CreateProcedure(viewProcedure.AsSetter())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func deleteProcedure(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/procedure/")
	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("parse error: %s", err))
		return
	}
	err = db.DeleteProcedure(int32(intID))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("db error: %s", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getPatientTypes(w http.ResponseWriter, r *http.Request) {
	types, err := db.GetPatientTypes()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read settings: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(types)
}

func getProcedureTypes(w http.ResponseWriter, r *http.Request) {
	types, err := db.GetProcedureTypes()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read settings: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(types)
}

func getPatientFolders(w http.ResponseWriter, r *http.Request) {
	folders, err := db.GetPatientFolders()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read settings: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(folders)
}

func handleSetting(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		updateSetting(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		deleteSetting(w, r)
		return
	}
}

func updateSetting(w http.ResponseWriter, r *http.Request) {
	var viewSetting db.ViewSetting
	if err := json.NewDecoder(r.Body).Decode(&viewSetting); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("Invalid request body: %s", err))
		return
	}
	if viewSetting.ID.IsUnset() {
		if err := db.CreateSetting(viewSetting); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(db.FmtError("Failed to update setting: %s", err))
			return
		}
	} else {
		if err := db.UpdateSetting(viewSetting); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(db.FmtError("Failed to update setting: %s", err))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("Invalid request method: %s", r.Method))
		return
	}

	var viewSettings []db.ViewSetting
	if err := json.NewDecoder(r.Body).Decode(&viewSettings); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("Invalid request body: %s", err))
		return
	}
	for _, viewSetting := range viewSettings {
		if viewSetting.ID.IsUnset() {
			if err := db.CreateSetting(viewSetting); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(db.FmtError("Failed to update setting: %s", err))
				return
			}
		} else {
			if err := db.UpdateSetting(viewSetting); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(db.FmtError("Failed to update setting: %s", err))
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

func deleteSetting(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/setting/")
	intID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("parse error: %s", err))
		return
	}
	err = db.DeleteSetting(int32(intID))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(db.FmtError("db error: %s", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Handle preflight requests.
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// slog.Debug("request", "addr", r.RemoteAddr, "method", r.Method, "url", r.URL)
		h.ServeHTTP(w, r)
	})
}

//go:embed ui/static/*
var embeddedFiles embed.FS

var (
	port         = flag.Int("port", 8000, "the port to listen on")
	databaseFile = flag.String("db", "", "path to the database")
	useCors      = flag.Bool("cors", false, "enable cors")
	dbLog        = flag.Bool("dbLog", false, "enable database logging")
)

func main() {
	flag.Parse()
	if *databaseFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	programLevel := new(slog.LevelVar)
	programLevel.Set(slog.LevelDebug)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	err := db.InitializeDB(*databaseFile, *dbLog)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to initialize db")
	}

	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".mjs", "application/javascript")
	_ = mime.AddExtensionType(".cjs", "application/javascript")

	// Create a new ServeMux.
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/patient-list/", getPatientList)
	mux.HandleFunc("/v1/patient/", handlePatient)
	mux.HandleFunc("/v1/patients/", updatePatients)
	mux.HandleFunc("/v1/procedure/", handleProcedure)
	mux.HandleFunc("/v1/patient-types/", getPatientTypes)
	mux.HandleFunc("/v1/procedure-types", getProcedureTypes)
	mux.HandleFunc("/v1/patient-folder", getPatientFolders)
	mux.HandleFunc("/v1/setting", handleSetting)
	mux.HandleFunc("/v1/settings", updateSettings)

	subFS, err := fs.Sub(embeddedFiles, "ui/static")
	if err != nil {
		fmt.Printf("%s", err)
		panic("failed to initialize www")
	}

	fileServer := http.FileServer(http.FS(subFS))
	mux.Handle("/", fileServer)

	var handler http.Handler = mux
	if *useCors {
		handler = corsMiddleware(handler)
	}
	handler = loggingMiddleware(handler)

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	slog.Error("server", "err", http.ListenAndServe(addr, handler))
}
