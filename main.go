package main

import (
	"github.com/DockerContainerService/git-syncer/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	configFile, privateKeyFile string
	retries, routineNum        int
	debug                      bool
	version                    = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     "git-syncer [flags]",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		c, err := client.Default(configFile, retries, routineNum, privateKeyFile)
		if err != nil {
			logrus.Fatalf("err1: %v", err)
		}
		c.Run()
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVar(&privateKeyFile, "privateKeyFile", "", "private key file")
	rootCmd.PersistentFlags().IntVarP(&retries, "retries", "r", 3, "times to retry failed task")
	rootCmd.PersistentFlags().IntVarP(&routineNum, "proc", "p", 5, "numbers of worker")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
}

func main() {
	logger := logrus.StandardLogger()
	logger.SetReportCaller(true)
	logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	if err := rootCmd.Execute(); err != nil {
		logrus.Panic(err)
	}
}
