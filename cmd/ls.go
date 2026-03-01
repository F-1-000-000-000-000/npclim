package cmd
import (
	"fmt"
	"io/ioutil"
    "path/filepath"
    "regexp"
    "strings"
	"os"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lsCmd = &cobra.Command {
	Use:   "ls [directory]",
	Short: "List proxy hosts",
	Long: `List current proxy hosts and optionally show server/proxypass information.
If no directory is supplied and no default directory is set in the config file, 
the current directory is used.`,
	RunE: listProxies,
}

func init() {
	//-------------------------------------------------------------------------
	// Description:	Set flags
	//-------------------------------------------------------------------------
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolP("long", "l", false, "List server and proxy information")
}


type ProxyInfo struct {
	Filename   string
	ServerName string
	ProxyPass  string
}


func listProxies(cmd *cobra.Command, args []string) error {
	//-------------------------------------------------------------------------
	// Description:	Get proxyhosts and list filenames + server + proxy
	//-------------------------------------------------------------------------
	long, _ := cmd.Flags().GetBool("long")
	var dir = "."

	if len(args) == 1 { 
		dir = args[0] 
	} else {
		defaultDir := viper.GetString("proxy-location")
		if defaultDir != "" { dir = defaultDir }
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Couldn't read %s: %w", dir, err)
	}

	fmt.Printf("Proxies found in %s:\n\n", dir)

	var proxyInfos []ProxyInfo
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".conf") {
			continue
		}

		info, err := parseProxyFile(filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Printf("%v (could not read)\n", file.Name())
			continue
		}
		proxyInfos = append(proxyInfos, info)
	}


	if !long {
		for _, info := range proxyInfos {
			fmt.Println(info.Filename)
		}
		return nil
	}

	maxFileLen := 0
	maxServerLen := 0
	for _, info := range proxyInfos {
		if len(info.Filename) > maxFileLen {maxFileLen = len(info.Filename)}
		if len(info.ServerName) > maxServerLen {maxServerLen = len(info.ServerName)}
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	for _, info := range proxyInfos {
		fmt.Printf("%-*s  %*s -> %s\n", 
			maxFileLen, 
			info.Filename, 
			maxServerLen,
			info.ServerName,
			cyan(info.ProxyPass),
		)
	}
	return nil
}


func parseProxyFile(path string) (ProxyInfo, error) {
	//-------------------------------------------------------------------------
	// Description:	Read individual proxy file
	//-------------------------------------------------------------------------
    fileContent, err := os.ReadFile(path)
    if err != nil {
        return ProxyInfo{}, err
    }

    info := ProxyInfo{Filename: filepath.Base(path)}
    text := string(fileContent)

    if match := regexp.MustCompile(`server_name\s+([^;]+);`).FindStringSubmatch(text); match != nil {
        info.ServerName = strings.TrimSpace(match[1])
    }
    if match := regexp.MustCompile(`proxy_pass\s+([^;]+);`).FindStringSubmatch(text); match != nil {
        info.ProxyPass = strings.TrimSpace(match[1])
    }

    return info, nil
}