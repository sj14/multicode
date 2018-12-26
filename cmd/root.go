package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sj14/multidecode/protodec"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multidecode",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("needs exactly one input to decode")
		}
		input := []byte(args[0])

		unmarshalled := &protodec.Empty{}
		for {
			applied := false

			if err := proto.Unmarshal(input, unmarshalled); err == nil {
				input = []byte(unmarshalled.String())
				logVerbose("applied proto decoding:\n%s\n\n", unmarshalled.String())
				// we can't decode more on top of proto
				break
			}

			if b, err := base64.StdEncoding.DecodeString(string(input)); err == nil {
				input = b
				logVerbose("applied base64 decoding:\n%v\n\n", string(b))
				applied = true
			}

			if b, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(string(input), "0x"))); err == nil {
				input = b
				logVerbose("applied hex decoding:\n%v\n\n", string(b))
				applied = true
			}

			if !applied {
				break
			}
		}
		logVerbose("result:\n")
		fmt.Printf("%v", string(input))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	verbose bool
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func logVerbose(format string, v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, v...)
}
