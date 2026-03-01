package cmd
import (
    "fmt"
    "os"
	"path/filepath"
	"github.com/spf13/viper"
    "github.com/spf13/cobra"
)
var rmCmd = &cobra.Command{
	Use:   "rm <proxy>",
	Short: "Remove a proxy host",
	Long: `Remove a proxy configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: removeProxy,
}

func init() {
	//-------------------------------------------------------------------------
	// Description:	Set flags
	//-------------------------------------------------------------------------
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().StringP("proxy-location", "l", "./", "Directory containing proxy hosts")
}

func removeProxy(cmd *cobra.Command, args []string) error {
	//-------------------------------------------------------------------------
	// Description:	Remove proxy conf file
	//-------------------------------------------------------------------------
	proxyLocation := viper.GetString("proxy-location")
	proxyName := args[0]
	filePath := filepath.Join(proxyLocation, proxyName + ".conf")
	
	// file not found
	if _, err := os.Stat(filePath); os.IsNotExist(err) { return fmt.Errorf("no proxy host found for %s", args[0]) }

	err := os.Remove(filePath)
	if err != nil { return fmt.Errorf("error removing proxy host: %s", err) }
	fmt.Println("Succesfully removed ", filePath)

	return nil
}