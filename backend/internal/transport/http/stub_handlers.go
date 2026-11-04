package http

import (
	"net/http"

	"github.com/dinkmachinery/dink-core/internal/platform/response"
)

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request, partition string) {
	response.Error(w, r, http.StatusNotImplemented, response.ProblemDetails{
		Title:  "Not Implemented",
		Code:   "ERR_PARTITION_PENDING",
		Detail: "This endpoint will be activated during partition " + partition + ".",
	})
}

func (s *Server) ListCountryConfigs(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-11")
}

func (s *Server) GetCountryConfigByCode(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-11")
}

func (s *Server) UpdateCountryConfig(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-11")
}

func (s *Server) GetAdminUiConfig(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-11")
}

func (s *Server) SubmitBookingInspection(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-07")
}

func (s *Server) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-06")
}

func (s *Server) GetMyWallet(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-09")
}

func (s *Server) ListWalletTransactions(w http.ResponseWriter, r *http.Request) {
	s.notImplemented(w, r, "05-09")
}
