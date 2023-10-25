package cmd

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/mohibeyki/storage-bench/runner"
	"github.com/mohibeyki/storage-bench/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "storage-bench",
	Short: "A storage benchmark tool",
	Long: `Storage Benchmark is a simple storage benchmark tool (duh)
It supports local filesystem and s3 storage`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var total atomic.Uint64
		runner := runner.Runner{Total: &total}

		if runner.Path, err = cmd.PersistentFlags().GetString("path"); err != nil {
			_ = fmt.Errorf("could not parse Path from arguments. [%s]", err)
			panic(err)
		}

		if runner.Files, err = cmd.PersistentFlags().GetUint64("files"); err != nil {
			_ = fmt.Errorf("could not parse files from arguments. [%s]", err)
			panic(err)
		}

		if runner.Threads, err = cmd.PersistentFlags().GetUint64("threads"); err != nil {
			_ = fmt.Errorf("could not parse threads from arguments. [%s]", err)
			panic(err)
		}

		if runner.Size, err = cmd.PersistentFlags().GetUint64("size"); err != nil {
			_ = fmt.Errorf("could not parse size from arguments. [%s]", err)
			panic(err)
		}

		// Input size is in KiB
		runner.Size *= 1024

		var s3 bool
		if s3, err = cmd.PersistentFlags().GetBool("s3"); err != nil {
			_ = fmt.Errorf("could not parse size from arguments. [%s]", err)
			panic(err)
		}

		if s3 {
			runner.Writer = &writer.S3Writer{Bucket: viper.GetString("bucket"), Region: viper.GetString("region"), AccessKey: viper.GetString("accessKey"), SecretKey: viper.GetString("secretKey")}
		} else {
			runner.Writer = &writer.FSWriter{}
		}

		_ = runner.Run()

		fmt.Printf("finished writing [%d] files\n", total.Load())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("path", "p", "/tmp/bench", "path to store files")
	rootCmd.PersistentFlags().Uint64P("files", "f", 1024, "number of files to create")
	rootCmd.PersistentFlags().Uint64P("size", "s", 1024, "size of each file in KiB, use numbers divisible by 64")
	rootCmd.PersistentFlags().Uint64P("threads", "t", 8, "number of threads")
	rootCmd.PersistentFlags().Bool("s3", false, "use s3 backend instead of file system")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.storage-bench")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
