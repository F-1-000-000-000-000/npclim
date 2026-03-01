package cmd
import (
    "bytes"
    "fmt"
    "os"
	"path/filepath"
    "text/template"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)
var newCmd = &cobra.Command{
	Use:   "new [filename]",
	Short: "Create a proxy host",
	Long: `Create a new proxy configuration file using either the default simple
template or a provided template (path supplied as a flag or in the config). Outputs
to current directory unless location is specified as a flag or in the config. Filename
is accepted as an argument (prioritized) or as a template using -f or the config file.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: newProxy,
}

func init() {
	//-------------------------------------------------------------------------
	// Description:	Set flags
	//-------------------------------------------------------------------------
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("subdomain", "s", "", "Subdomain to be used for proxy host")
	newCmd.Flags().StringP("domain", "d", "", "Base domain to be used for proxy host")
	newCmd.Flags().StringP("host", "H", "localhost", "Host to forward traffic to")
	newCmd.Flags().IntP("port", "p", 0, "Port number to forward traffic to")
	newCmd.Flags().StringP("configuration-template", "t", "", "File to be used as template")
	newCmd.Flags().StringP("proxy-location", "l", "./", "Directory containing proxy hosts")
	newCmd.Flags().StringP("filename-template", "f", "{{.Subdomain}}.{{.Domain}}.conf", "Filename template")
}


const defaultTemplate = `server {
    listen 80;
    server_name {{.Subdomain}}.{{.Domain}};

    location / {
        proxy_pass http://{{.Host}}:{{.Port}};
    }
}
`

type ProxyConfig struct {
    Subdomain string
    Domain    string
    Host	  string
	Port      int
}


func newProxy(cmd *cobra.Command, args []string) error {
	//-------------------------------------------------------------------------
	// Description:	Create a new proxy host configuration
	//-------------------------------------------------------------------------
  	viper.BindPFlags(cmd.Flags())

	// required flags:
	subdomain := viper.GetString("subdomain")
	if subdomain == "" { return fmt.Errorf("subdomain is required (set via -s or config file)")}
	domain := viper.GetString("domain")
	if domain == "" { return fmt.Errorf("domain is required (set via -d or config file)") }
	port := viper.GetInt("port")
	if port == 0 { return fmt.Errorf("port is required (set via -p or config file)") }

	// optional flags:
	host := viper.GetString("host")
	templateFile := viper.GetString("configuration-template")
	filename := viper.GetString("filename-template")
	outputDir := viper.GetString("proxy-location")

	// save contents
    proxyVariables := ProxyConfig {
        Subdomain: subdomain,
        Domain:    domain,
		Host:      host,
        Port:      port,
    }

	// determine filename
	var finalFilename string
	if len(args) > 0 {
    	finalFilename = args[0] + ".conf"
	} else {
		var buf bytes.Buffer
		filenameTemplate, err := template.New("filename").Parse(filename)
		if err != nil { return fmt.Errorf("failed to load filename template: %w", err) }
		if err := filenameTemplate.Execute(&buf, proxyVariables); err != nil {
			return fmt.Errorf("failed to render filename template: %w", err)
		}
		finalFilename = buf.String()
	}

	// create proxy host conf
	fileTemplate, err := loadTemplate(templateFile)
    if err != nil { return fmt.Errorf("failed to load config template: %w", err) }
	var proxyHostContents bytes.Buffer
    if err := fileTemplate.Execute(&proxyHostContents, proxyVariables); err != nil {
        return fmt.Errorf("failed to render config template: %w", err)
    }

	// actually write to disk
	fullFilePath := filepath.Join(outputDir, finalFilename)
	err = os.WriteFile(fullFilePath, proxyHostContents.Bytes(), 0644)
	if err != nil { return fmt.Errorf("failed to write file: %w", err) }

	fmt.Printf("Created %s\n", fullFilePath)
	return nil
}


func loadTemplate(customPath string) (*template.Template, error) {
	//-------------------------------------------------------------------------
	// Description:	Load template. Order of prio: flag > .config > default
	//-------------------------------------------------------------------------
	// prioritize -t flag
	if customPath != "" {
        fileContent, err := os.ReadFile(customPath)
        if err != nil { return nil, fmt.Errorf("could not read template file: %w", err) }
        return template.New("proxy").Parse(string(fileContent))
    } 
	// then check for template in config dir
	fileContent, err := os.ReadFile(os.ExpandEnv("$HOME/.config/npclim/template.conf"))
	if err != nil && !os.IsNotExist(err) { return nil, err }

	// file exists and can read :)
	if err == nil { return template.New("proxy").Parse(string(fileContent)) }
	
	// else fallback to hardcoded template
	return template.New("proxy").Parse(defaultTemplate)
}


