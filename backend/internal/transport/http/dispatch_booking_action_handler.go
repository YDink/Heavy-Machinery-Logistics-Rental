package http

import (
	"net/http"
	"strings"
)

// DispatchBookingAction dispatches colon-action syntax POST /bookings/{id}:{action} to specific handlers
func (s *Server) DispatchBookingAction(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	parts := strings.SplitN(action, ":", 2)
	if len(parts) == 2 {
		id := parts[0]
		verb := parts[1]
		r.SetPathValue("id", id)
		switch verb {
		case "confirm-payment":
			s.ConfirmBookingPayment(w, r)
			return
		case "decline-payment":
			s.DeclineBookingPayment(w, r)
			return
		case "dispute":
			s.DisputeBookingPayment(w, r)
			return
		case "cancel":
			s.CancelBooking(w, r)
			return
		case "handover-check-in":
			s.HandoverCheckIn(w, r)
			return
		case "return-check-out":
			s.ReturnCheckOut(w, r)
			return
		}
	}
	s.notImplemented(w, r, "05-05")
}
