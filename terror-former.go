package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "CLI is a simple command line interface",
		Long: `CLI is a simple command line interface with 3 commands:
        plan - describes what will be done
        apply - applies the changes
        destroy - removes the changes`,
	}

	var fileName string

	var cmdPlan = &cobra.Command{
		Use:   "plan",
		Short: "Describes what will be done",
		Run: func(cmd *cobra.Command, args []string) {

			// Read the file
			content, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			// Send the contnet to LLM api  and get the result
			// Define the payload structure
			type Payload struct {
				Model    string `json:"model"`
				Messages []struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"messages"`
			}

			// Create the payload
			payload := Payload{
				Model: "gpt-3.5-turbo",
				Messages: []struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				}{
					{
						Role:    "system",
						Content: "You are a helpful assistant.",
					},
					{
						Role:    "user",
						Content: string(content),
					},
				},
			}

			// Marshal the payload into JSON
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Error marshaling payload:", err)
				return
			}

			// Create a new request
			req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(payloadBytes))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			// Set the headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			// Read the response
			respContent, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			// Print the response
			fmt.Println("Response from LLM:")
			fmt.Println(string(respContent))
		},
	}

	cmdPlan.Flags().StringVarP(&fileName, "file", "f", "terror.txt", "File to read")

	var cmdApply = &cobra.Command{
		Use:   "apply",
		Short: "Applies the changes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Apply command placeholder")
			fmt.Println("File to read:", fileName)
		},
	}

	cmdApply.Flags().StringVarP(&fileName, "file", "f", "terror.txt", "File to read")

	var cmdDestroy = &cobra.Command{
		Use:   "destroy",
		Short: "Destroys the changes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Destroy command placeholder")
			fmt.Println("File to read:", fileName)
		},
	}

	cmdDestroy.Flags().StringVarP(&fileName, "file", "f", "terror.txt", "File to read")

	rootCmd.AddCommand(cmdPlan, cmdApply, cmdDestroy)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
