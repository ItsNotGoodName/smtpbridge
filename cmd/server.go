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
	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/ItsNotGoodName/smtpbridge/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverConfig *config.Config

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		serverConfig.Load()

		// Start server
		server.Start(serverConfig)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverConfig = config.New()

	serverCmd.Flags().Bool("http", serverConfig.HTTP.Enable, "enable http server")
	viper.BindPFlag("http.enable", serverCmd.Flags().Lookup("http"))

	serverCmd.Flags().String("http-host", serverConfig.HTTP.Host, "http host address to listen on")
	viper.BindPFlag("http.host", serverCmd.Flags().Lookup("http-host"))

	serverCmd.Flags().Uint16("http-port", serverConfig.HTTP.Port, "http port to listen on")
	viper.BindPFlag("http.port", serverCmd.Flags().Lookup("http-port"))

	serverCmd.Flags().Bool("smtp", serverConfig.SMTP.Enable, "enable smtp server")
	viper.BindPFlag("smtp.enable", serverCmd.Flags().Lookup("smtp"))

	serverCmd.Flags().String("smtp-host", serverConfig.SMTP.Host, "smtp host address to listen on")
	viper.BindPFlag("smtp.host", serverCmd.Flags().Lookup("smtp-host"))

	serverCmd.Flags().Uint16("smtp-port", serverConfig.SMTP.Port, "smtp port to listen on")
	viper.BindPFlag("smtp.port", serverCmd.Flags().Lookup("smtp-port"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
