package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

func main() {
	fmt.Println("Hello, Hello, Hello")
	for _, v := range os.Environ() {
		fmt.Println(v)
	}

	ctx := context.Background()
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	))
	client := github.NewClient(tc)

	repo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	alerts, resp, err := client.CodeScanning.ListAlertsForRepo(ctx, repo[0], repo[1], &github.AlertListOptions{})
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failure: %d\n", resp.StatusCode)
	}

	for _, alert := range alerts {
		fmt.Printf("alert: %v\n", alert)
	}
}
