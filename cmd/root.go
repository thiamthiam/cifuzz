package cmd

import (
	"errors"
	"fmt"
	"os"

	"code-intelligence.com/cifuzz/pkg/cmdutils"
	"code-intelligence.com/cifuzz/pkg/storage"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "cifuzz",
	Short: "#tbd",
	// We are using our the custom ErrSilent instead to support a more specific
	// error handling
	SilenceErrors: true,
	SilenceUsage:  true,
	// We are using PersistentPreRun to call setup to make sure that all flags/arguments are available
	PersistentPreRun: setup,
}

var fs *afero.Afero

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false,
		"Show more verbose output, can be helpful for debugging problems")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func setup(cmd *cobra.Command, args []string) {
	fs = storage.WrapFileSystem()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {

		// Errors that are not ErrSilent are not expected and we want to show their full stacktrace
		if !errors.Is(err, cmdutils.ErrSilent) {
			fmt.Printf("%+v \n", err)
			fmt.Println(rootCmd.UsageString())
		}

		os.Exit(1)
	}
}