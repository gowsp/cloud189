package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "file direct link sharing",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		handler, err := App().Share("/", args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/", handler)
		log.Println("start share serve at", args[0])
		http.ListenAndServe(args[0], mux)
	},
}
