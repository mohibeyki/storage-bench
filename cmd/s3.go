package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Runs benchmark on a s3 bucket",
	Long: `Runs the same benchmark in an s3 bucket
authentication detail are read from the config file,
same path is used in the bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3 called")
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
