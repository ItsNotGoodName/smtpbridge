/*
Copyright © 2022 ItsNotGoodName

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/pkg/interrupt"
	"github.com/ItsNotGoodName/smtpbridge/server"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	serverConfig *config.Config
	watch        *bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smtpbridge",
	Short: "Bridge email to other messaging services.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("cmd.rootCmd.Run: using config file:", viper.ConfigFileUsed())

		restartCh := make(chan struct{})
		if *watch {
			viper.OnConfigChange(func(in fsnotify.Event) {
				log.Println("cmd.rootCmd.Run: config file changed, restarting")
				restartCh <- struct{}{}
			})
			viper.WatchConfig()
		}

		osCtx := interrupt.Context()

		// Start
		serverConfig.Load()
		ctx, cancel := context.WithCancel(osCtx)
		serverCh := server.Start(ctx, serverConfig)

		// Signals
		for {
			select {
			case <-restartCh:
				// Restart
				cancel()
				<-serverCh
				serverConfig = config.New()
				serverConfig.Load()
				ctx, cancel = context.WithCancel(osCtx)
				serverCh = server.Start(ctx, serverConfig)
			case <-serverCh:
				// Shutdown
				cancel()
				return
			}
		}
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
	cobra.OnInitialize(initConfig)

	serverConfig = config.New()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smtpbridge.yaml)")

	watch = rootCmd.Flags().Bool("watch", false, "restart when config file changes")

	rootCmd.Flags().Bool("http-disable", serverConfig.HTTP.Disable, "disable http server")
	viper.BindPFlag("http.disable", rootCmd.Flags().Lookup("http-disable"))

	rootCmd.Flags().String("http-host", serverConfig.HTTP.Host, "http host address to listen on")
	viper.BindPFlag("http.host", rootCmd.Flags().Lookup("http-host"))

	rootCmd.Flags().Uint16("http-port", serverConfig.HTTP.Port, "http port to listen on")
	viper.BindPFlag("http.port", rootCmd.Flags().Lookup("http-port"))

	rootCmd.Flags().Bool("smtp-disable", serverConfig.SMTP.Disable, "disable smtp server")
	viper.BindPFlag("smtp.disable", rootCmd.Flags().Lookup("smtp-disable"))

	rootCmd.Flags().String("smtp-host", serverConfig.SMTP.Host, "smtp host address to listen on")
	viper.BindPFlag("smtp.host", rootCmd.Flags().Lookup("smtp-host"))

	rootCmd.Flags().Uint16("smtp-port", serverConfig.SMTP.Port, "smtp port to listen on")
	viper.BindPFlag("smtp.port", rootCmd.Flags().Lookup("smtp-port"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".smtpbridge" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".smtpbridge")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
