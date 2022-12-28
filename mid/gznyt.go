package mid

import (
	"bytes"
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// NYTimes ...
// https://github.com/NYTimes/gziphandler

func gz(htm *bytes.Buffer) {
	withoutGz := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		//io.WriteString(w, htm.String())
		htm.WriteTo(w)
	})

	withGz := gziphandler.GzipHandler(withoutGz)

	http.Handle("/", withGz)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
