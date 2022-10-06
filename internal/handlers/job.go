package handlers

import (
	"encoding/json"
	"github.com/supernova0730/job/internal/repository"
	"net/http"
)

type JobID struct {
	Code string `json:"code"`
}

type JobSchedule struct {
	Code     string `json:"code"`
	Schedule string `json:"schedule"`
}

type JobHandler struct {
	jobRepo *repository.JobRepository
}

func New(jobRepo *repository.JobRepository) *JobHandler {
	return &JobHandler{
		jobRepo: jobRepo,
	}
}

func (h *JobHandler) StartJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var b JobID
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.jobRepo.SetActive(r.Context(), b.Code, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *JobHandler) StopJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var b JobID
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.jobRepo.SetActive(r.Context(), b.Code, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *JobHandler) ChangeSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var b JobSchedule
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.jobRepo.SetSchedule(r.Context(), b.Code, b.Schedule)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *JobHandler) Start(address string) error {
	http.HandleFunc("/api/v1/job/start", h.StartJob)
	http.HandleFunc("/api/v1/job/stop", h.StopJob)
	http.HandleFunc("/api/v1/job/schedule", h.ChangeSchedule)
	return http.ListenAndServe(address, nil)
}
