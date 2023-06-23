package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gnaydenova/notification-system/app/notifications"
	"github.com/gnaydenova/notification-system/app/notifications/dto"
)

type NotificationHandler struct {
	producer notifications.Producer
}

func NewNotificationHandler(p notifications.Producer) *NotificationHandler {
	return &NotificationHandler{producer: p}
}

func (h *NotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var n dto.Notification
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil || n.Message == "" || n.Channel == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.producer.Produce(r.Context(), n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
