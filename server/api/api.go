package api

import (
	"encoding/json"
	"fmt"
	"go-rest-api/logger"
	"go-rest-api/types/datastore"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	yayerror "go-rest-api/errors"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	APILogger = logger.GetLogger("api")
)

const maxRetryCount int = 5

type APIServer struct {
	port       int
	httpserver *http.Server
	logger     *logrus.Entry
	db         datastore.DataStore

	_stopChannel chan bool
	_serving     bool
}

func NewAPIServer(db datastore.DataStore, port int) *APIServer {
	httpserver := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  90 * time.Second,
	}
	apiServer := APIServer{
		port:       port,
		httpserver: httpserver,
		db:         db,
		logger:     APILogger,
	}

	return &apiServer
}

func (server *APIServer) Serve() bool {
	router := mux.NewRouter()
	router.Use(server.loggingMiddleware)
	router.HandleFunc("/orders", server.OrderHandler).Methods("POST")
	router.HandleFunc("/orders", server.OrderHandler).Methods("GET").Queries("page", "{page}").Queries("limit", "{limit}")
	router.HandleFunc("/orders/{id}", server.OrderHandler).Methods("PUT")
	router.PathPrefix("/").Handler(http.HandlerFunc(server.uncatchRequest))
	// router.NotFoundHandler = http.HandlerFunc(server.notFound)
	server.httpserver.Handler = router

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		for {
			endTime := time.Now().Add(15 * time.Second)
			for retryCount := 0; time.Now().Before(endTime) && retryCount < maxRetryCount; retryCount++ {
				APILogger.Infof("api server try Listening on port: '%v'", server.port)
				if err := server.httpserver.ListenAndServe(); err != nil {
					APILogger.Errorln("listen error:", err)
				}
				time.Sleep(2 * time.Second)
			}

			// if within 5 second of retry end time exit the program, else retry
			if time.Now().Add(5 * time.Second).Before(endTime) {
				APILogger.Errorf("Cannot listen on port '%s' in 15s, Now exiting...", server.port)
				os.Exit(yayerror.EXITCODE_ERVER_STARTUP_ERROR)
				break
			}
		}
	}()
	time.Sleep(200 * time.Millisecond)

	return true
}
func (server *APIServer) uncatchRequest(w http.ResponseWriter, r *http.Request) {
	server.logger.Warnf("Uncatch Request: [%v] %v", r.Method, r.URL)
	w.WriteHeader(400)
	outgoingJSON, err := json.Marshal(yayerror.API400)
	if err != nil {
		server.logger.Errorf("Cant marshal json: %+v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(outgoingJSON))
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (server *APIServer) loggingMiddleware(next http.Handler) http.Handler {
	var apiRequestLogger = logger.GetLogger("apiRequest")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}
		o.Header().Set("Content-Type", "application/json")

		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		apiRequestLogger.Infof("%s - - [%s] %q %d %d %q %q",
			addr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
			o.status,
			o.written,
			r.Referer(),
			r.UserAgent())
		next.ServeHTTP(o, r)
	})
}
