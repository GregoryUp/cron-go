package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func fillUsers() {
	urlApi, _ := os.LookupEnv("URL")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(urlApi + "/access-keys")

	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(resp.Body)

	resp.Body.Close()

	fmt.Println(string(body))

	return

	db, err := sql.Open("mysql", "root:rootroot@/outline_db")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err != nil {
		panic(err)
	}

	var accessKeys []AccessKey
	var jsonData map[string]json.RawMessage
	
	err = json.Unmarshal(body, &jsonData)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonData["accessKeys"], &accessKeys)

	if err != nil {
		panic(err)
	}

	stmt, _ := db.Prepare("INSERT INTO users SET vpn_key_id = ?, name = ?, access_key = ?, date_until = NOW()")

	defer stmt.Close()

	for _, key := range accessKeys {
		stmt.Exec(key.Id, key.Name, key.AccessUrl)
	}
}