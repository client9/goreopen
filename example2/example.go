package main

/* Similar to previous example but uses a BufferedFileWriter
 * When buf is full OR every 30 seconds, the buffer is flushed to disk
 *
 * care is done to make sure transient or partial log messages are not written.
 *
 * Note the signal handler catches SIGTERM to flush out and existing buffers
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
	bf := reopen.NewBufferedFileWriter(f)
	log.SetOutput(bf)

	sighup := make(chan os.Signal, 2)
	signal.Notify(sighup, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		for {
			s := <-sighup
			switch s {
			case syscall.SIGHUP:
				fmt.Printf("Got a sighup\n")
				bf.Reopen()
			case syscall.SIGTERM:
				fmt.Printf("Got SIGTERM\n")
				// make sure any remaining logs are flushed out
				bf.Close()
				os.Exit(0)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.URL.Path)
		fmt.Fprintf(w, "%s\n", r.URL.Path)
	})
	log.Fatal(http.ListenAndServe("127.0.0.1:8123", nil))
}
