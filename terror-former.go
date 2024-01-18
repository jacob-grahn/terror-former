package main

import (
	"fmt"
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

	var cmdPlan = &cobra.Command{
		Use:   "plan",
		Short: "Describes what will be done",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Plan command placeholder")
		},
	}

	var cmdApply = &cobra.Command{
		Use:   "apply",
		Short: "Applies the changes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Apply command placeholder")
		},
	}

	var cmdDestroy = &cobra.Command{
		Use:   "destroy",
		Short: "Removes the changes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Destroy command placeholder")
		},
	}

	rootCmd.AddCommand(cmdPlan, cmdApply, cmdDestroy)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
