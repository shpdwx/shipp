/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()

		if err := process(ctx); err != nil {
			fmt.Printf("run failed: %v", err)
			return
		}

		fmt.Println("ok!")

	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func process(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	select {
	case <-time.After(5 * time.Second):
		return errors.New("timeout")
	case <-ctx.Done():
		return ctx.Err()
	}

	time.Sleep(6 * time.Second)
	fmt.Println("it is time")

	return nil
}
