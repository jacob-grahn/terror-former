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
				Model          string `json:"model"`
				ResponseFormat struct {
					Type string `json:"type"`
				} `json:"response_format"`
				Messages []struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"messages"`
			}

			type Response struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}

			// Create the payload
			payload := Payload{
				Model: "gpt-3.5-turbo-1106",
				ResponseFormat: struct {
					Type string `json:"type"`
				}{
					Type: "json_object",
				},
				Messages: []struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				}{
					{
						Role:    "system",
						Content: "You always respond with JSON. You are a helpful senior infrastructure engineer who is an expert using Terraform to write infrastructure as code.",
					},
					{
						Role:    "user",
						Content: "Think about it step by step, and create a list of terraform modules that will be required to provision infrastructure that meets the following requirements. Keep in mind you will also need a core module that provisions one or more of the other modules. Format your response like {\"module_1\": \"description\", \"module_2\": \"description\", ...}\" \n --- \n" + string(content),
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

			// Unmarshal the response
			var respContentParsed Response
			err = json.Unmarshal(respContent, &respContentParsed)
			if err != nil {
				fmt.Println("Error unmarshaling response:", err)
				return
			}

			// Print the content
			fmt.Println("Content from LLM:")
			fmt.Println(respContentParsed.Choices[0].Message.Content)
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
