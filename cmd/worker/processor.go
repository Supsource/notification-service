package main

import (
	"log"
	"notification-service/internal/model"
	"notification-service/internal/repository"
)

func processNotification(id string, repo *repository.PostgresNotificationRepo) {
	n, err := repo.GetByID(id)
	if err != nil {
		log.Println("notifiction not found:", err)
		return
	}

	// mark processing here
	if err := repo.UpdateStatus(id, model.StatusProcessing, nil); err != nil {
		return
	}

	// simulate sending here
	err = sendNotification(n)

	if err != nil {
		repo.MarkFailed(n)
		return
	}

	repo.UpdateStatus(id, model.StatusSent, nil)
}
