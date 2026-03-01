package cmd
import (
	"fmt"
	"path/filepath"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit <proxy>",
	Short: "Edit a proxy host configuration using the system editor",
	Args:  cobra.ExactArgs(1),
	RunE:  editProxy,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editProxy(cmd *cobra.Command, args []string) error {
	proxyLocation := viper.GetString("proxy-location")
	proxyName := args[0]
	filePath := filepath.Join(proxyLocation, proxyName + ".conf")

	// file not found
	if _, err := os.Stat(filePath); os.IsNotExist(err) { return fmt.Errorf("no proxy host found for %s", args[0]) }

	editor := os.Getenv("EDITOR")
	if editor == "" { editor = "vim" }

	editorCmd := exec.Command(editor, filePath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	return editorCmd.Run()
}