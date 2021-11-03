package main

import "log"

func Success(msg string) {
	log.Println("[SUCCESS]", msg)
}

func Error(msg string, err error) {
	log.Printf("[ERROR] %v: %v", msg, err)
}

func Warm(msg string, err error) {
	log.Println("[WARM]", msg)
}

func ToLine(lid, id, types string) {
	log.Printf("[MESSAGE] | %33s | <-- | %18s | %7s |", lid, id, types)
}

func ToDiscord(lid, id, types string) {
	log.Printf("[MESSAGE] | %33s | --> | %18s | %7s |", lid, id, types)
}
