package main

import (
	"legolas/common/models"
	"log"
)

func jobHandler(data string) {
	log.Println("Get =>")
	log.Println(data)

	job, err := models.JobFromJson([]byte(data))
	if err != nil {
		log.Printf("Unmarshalling failed: %v\n", err)
	}

	log.Println(job)
}
