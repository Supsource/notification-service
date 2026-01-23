package main

import (
	"log"
	"notification-service/internal/model"
	"notification-service/internal/queue"
	"notification-service/internal/repository"
)

const MaxRetries = 3

func processNotification(
	id string,
	repo repository.NotificationRepository,
	producer *queue.Producer,
	dlqProducer *queue.Producer,
) {
	n, err := repo.GetByID(id)
	if err != nil {
		log.Println("notification not found:", err)
		return
	}

	// mark processing here
	if err := repo.UpdateStatus(id, model.StatusProcessing, nil); err != nil {
		return
	}

	if err := sendNotification(n); err != nil {
		repo.IncrementRetry(n.ID)

		if n.RetryCount+1 >= MaxRetries {
			repo.UpdateStatus(n.ID, model.StatusFailed, strPtr(err.Error()))
			dlqProducer.Enqueue(n.ID)
			return
		}

		repo.UpdateStatus(n.ID, model.StatusPending, strPtr(err.Error()))
		producer.Enqueue(n.ID)
		return
	}

	repo.UpdateStatus(id, model.StatusSent, nil)
}

func sendNotification(n *model.Notification) error {
	log.Printf("Sending notification %s to %s: %s", n.Type, n.UserID, n.Title)
	// will add actual logic here (email, push, etc.)
	return nil
}

func strPtr(s string) *string {
	return &s
}
