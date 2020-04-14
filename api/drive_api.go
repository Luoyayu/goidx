package api

import (
	"encoding/json"
	"log"
)

type SDrive struct {
	Kind string `json:"kind"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetSharedDrives(auth string) []*SDrive {
	opName := "GetSharedDrives"

	resp := Do("GET", opName, "~_~_gdindex/drives", "", auth)
	defer resp.Body.Close()

	a := struct {
		Kind   string    `json:"kind"`
		Drives []*SDrive `json:"drives"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		log.Fatal(opName, err)
	} else {
		return a.Drives
	}
	return nil
}
