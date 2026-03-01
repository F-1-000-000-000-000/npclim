package cmd
import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "npclim",
	Short: "Manage proxy host configuration files",
	Long: "Manage proxy host config files for nginx (and nginx-like programs) using the CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProxies(cmd, args)
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
  rootCmd.Flags().BoolP("long", "l", false, "Show detailed proxy information")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/npclim")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Error reading config:", err)
		}
	}
}


