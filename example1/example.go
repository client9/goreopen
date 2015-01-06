package main

/* Simple logrotate logger
 */
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/client9/reopen"
)

func main() {
	f, err := reopen.NewFileWriter("/tmp/example.log")
	if err != nil {
		log.Fatalf("Unable to set output log: %s", err)
	}
	log.SetOutput(f)

	// channel is number of signals needed to catch  (more or less)
	// we only are working with one here, SIGUP
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)
	go func() {
		for {
			<-sighup
			fmt.Printf("Got a sighup\n")
			f.Reopen()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.URL.Path)
		fmt.Fprintf(w, "%s\n", r.URL.Path)
	})
	log.Fatal(http.ListenAndServe("127.0.0.1:8123", nil))
}
