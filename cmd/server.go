/*
Copyright © 2021 ItsNotGoodName

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
	"os"
	"path"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/left/smtp"
	"github.com/ItsNotGoodName/smtpbridge/right/database"
	"github.com/ItsNotGoodName/smtpbridge/right/endpoint"
	"github.com/ItsNotGoodName/smtpbridge/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start SMTP server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Read config
		config := app.NewConfig()

		// Init repositories
		endpointREPO := endpoint.NewRepository(config.ConfigEndpoints)
		databaseREPO := database.NewDB(config.DBPath, config.AttachmentsPath)

		// Init services
		authSVC := service.NewMockAuth()
		bridgeSVC := service.NewBridge(endpointREPO, config.Bridges)
		messageSVC := service.NewMessage(bridgeSVC, endpointREPO, databaseREPO, databaseREPO)

		// Init smtp server
		smtpServer := smtp.New(authSVC, messageSVC, config.ConfigSMTP)

		// Start smtp server
		smtpServer.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("host", "", "host to listen on")
	viper.BindPFlag("smtp.host", serverCmd.Flags().Lookup("host"))

	serverCmd.Flags().Uint16("port", 1025, "port to listen on")
	viper.BindPFlag("smtp.port", serverCmd.Flags().Lookup("port"))

	serverCmd.Flags().Int("size", 1024*1024*25, "max size of email in bytes")
	viper.BindPFlag("smtp.size", serverCmd.Flags().Lookup("size"))

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	rootPath := path.Join(home, ".smtpbridge")
	os.MkdirAll(rootPath, os.ModePerm)
	cobra.CheckErr(err)

	serverCmd.Flags().String("db", path.Join(rootPath, "smtpbridge.db"), "database file")
	viper.BindPFlag("db", serverCmd.Flags().Lookup("db"))

	serverCmd.Flags().String("attachments", path.Join(rootPath, "attachments"), "attachments directory")
	viper.BindPFlag("attachments", serverCmd.Flags().Lookup("attachments"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
