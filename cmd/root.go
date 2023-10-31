package cmd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mohibeyki/storage-bench/reader"
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
		runner := runner.Runner{}

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

		// Input size is in 64KiB
		runner.Size *= 64 * 1024

		var s3 bool
		if s3, err = cmd.PersistentFlags().GetBool("s3"); err != nil {
			_ = fmt.Errorf("could not parse size from arguments. [%s]", err)
			panic(err)
		}

		if s3 {
			session, err := session.NewSession(&aws.Config{
				Region:      aws.String(viper.GetString("region")),
				Credentials: credentials.NewStaticCredentials(viper.GetString("accessKey"), viper.GetString("secretKey"), "")},
			)
			if err != nil {
				panic(err)
			}

			runner.Writer = &writer.S3Writer{Bucket: viper.GetString("bucket"), Session: session}
		} else {
			runner.Writer = &writer.FSWriter{}
		}

		if zero, err := cmd.PersistentFlags().GetBool("zero"); err != nil {
			_ = fmt.Errorf("could not parse size from arguments. [%s]", err)
			panic(err)
		} else {
			if zero {
				runner.ReaderType = reflect.TypeOf(reader.ZeroReader{})
			} else {
				runner.ReaderType = reflect.TypeOf(reader.RandomReader{})
			}
		}

		_ = runner.Run()

		fmt.Printf("finished writing [%d] files\n", runner.Files)
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
	rootCmd.PersistentFlags().Uint64P("size", "s", 64, "size of each file in 64KiB")
	rootCmd.PersistentFlags().Uint64P("threads", "t", 8, "number of threads")
	rootCmd.PersistentFlags().BoolP("zero", "z", false, "write all zeroes or use random data")
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
