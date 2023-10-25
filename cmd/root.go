package cmd

import (
	"fmt"
	"github.com/mohibeyki/storage-bench/runner"
	"github.com/mohibeyki/storage-bench/writer"
	"github.com/spf13/cobra"
	"os"
	"sync/atomic"
)

var rootCmd = &cobra.Command{
	Use:   "storage-bench",
	Short: "A storage benchmark tool",
	Long: `Storage Benchmark is a simple storage benchmark tool (duh)
It supports local filesystem and s3 storage`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		var files uint64
		var size uint64
		var threads uint64
		var err error

		if path, err = cmd.PersistentFlags().GetString("path"); err != nil {
			_ = fmt.Errorf("could not parse Path from arguments. [%s]\n", err)
			return
		}

		if files, err = cmd.PersistentFlags().GetUint64("files"); err != nil {
			_ = fmt.Errorf("could not parse files from arguments. [%s]\n", err)
			return
		}

		if size, err = cmd.PersistentFlags().GetUint64("size"); err != nil {
			_ = fmt.Errorf("could not parse size from arguments. [%s]\n", err)
			return
		}

		if threads, err = cmd.PersistentFlags().GetUint64("threads"); err != nil {
			_ = fmt.Errorf("could not parse threads from arguments. [%s]\n", err)
			return
		}

		var total atomic.Uint64
		var fsWriter writer.FSWriter
		runner := runner.Runner{Path: path, Files: files, Size: size, Threads: threads, Total: &total, Writer: &fsWriter}
		_ = runner.Run()

		fmt.Printf("finished writing [%d] files", total.Load())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", "/tmp/bench", "path to store files")
	rootCmd.PersistentFlags().Uint64P("files", "f", 1024, "number of files to create")
	rootCmd.PersistentFlags().Uint64P("size", "s", 1024, "size of each file in KiB")
	rootCmd.PersistentFlags().Uint64P("threads", "t", 8, "number of threads")
}
