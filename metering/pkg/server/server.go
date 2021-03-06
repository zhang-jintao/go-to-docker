package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	goruntime "runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tangfeixiong/go-to-docker/metering/pb"
	"github.com/tangfeixiong/go-to-docker/metering/pkg/collector"
	"github.com/tangfeixiong/go-to-docker/metering/pkg/exporter"
)

type myServer struct {
	grpcHost               string
	httpHost               string
	packgesHome            string
	redisSentinelAddresses string
	redisAddresses         []string
	redisDB                int
	etcdAddresses          string
	mysqlAddress           string
	gnatsdAddresses        string
	kafkaAddresses         string
	zookeeperAddresses     string
	rabbitAddress          string
	priorityCMDB           []string
	priorityMQ             []string
	subSubject             string
	subQueue               string
	pubSubject             string
	cmCache                string
	unsubCh                chan string
	dispatchersignal       chan bool
}

type myExporter struct {
	myServer
	exportermanager *exporter.ExporterManager
}

type myCollector struct {
	myServer
	collectormanager *collector.CollectorManager
}

func RunExporter(meterdriver, collectorrpc string) {
	m := new(myExporter)
	m.grpcHost = ":12315"
	m.httpHost = ":12316"
	m.exportermanager = new(exporter.ExporterManager)
	m.exportermanager.Dispatchers = make(map[string]exporter.MeterDispatcher)
	m.exportermanager.MeteringNameURLs = make(map[string][]string)
	// m.exportermanager.MetricsCollectorRPC = "localhost:12305"
	// m.exportermanager.MeteringNameURLs["docker"] = []string{"unix:///var/run/docker.sock"}

	if v, ok := os.LookupEnv("METERINGEXPORTER_GRPC_PORT"); ok && 0 != len(v) {
		if strings.Contains(v, ":") {
			m.grpcHost = v
		} else {
			m.grpcHost = "localhost:" + v
		}
	}
	if v, ok := os.LookupEnv("METERINGEXPORTER_HTTP_PORT"); ok && 0 != len(v) {
		if strings.Contains(v, ":") {
			m.httpHost = v
		} else {
			m.httpHost = ":" + v
		}
	}

	if meterdriver != "" {
		if err := os.Setenv("METERING_DRIVER_LIST", meterdriver); err != nil {
			fmt.Println("Environment not set, error:", err)
		}
	}
	if v, ok := os.LookupEnv("METERING_DRIVER_LIST"); ok && 0 != len(v) {
		for _, pair := range strings.Split(v, ",") {
			nu := strings.Split(pair, "=")
			if len(nu) != 2 {
				panic("Invalid environment value: METERING_DRIVER_LIST")
			}
			m.exportermanager.MeteringNameURLs[nu[0]] = strings.Split(nu[1], ";")
		}
	}

	if collectorrpc != "" {
		if err := os.Setenv("METRICS_COLLECTOR_RPC", collectorrpc); err != nil {
			fmt.Println("Environment not set, error:", err)
		}
	}
	if v, ok := os.LookupEnv("METRICS_COLLECTOR_RPC"); ok && 0 != len(v) {
		if strings.HasPrefix(v, "grpc=") {
			m.exportermanager.MetricsCollectorRPC = v[5:]
		} else {
			fmt.Println("RPC driver not support,", v)
		}
	}

	if len(m.exportermanager.MeteringNameURLs) == 0 {
		panic("Meters were not available")
	}
	m.run()
}

func RunCollector(storage string) {
	m := new(myCollector)
	m.grpcHost = ":12305"
	m.httpHost = ":12306"
	m.collectormanager = new(collector.CollectorManager)
	m.collectormanager.MetricsStorageDuration = time.Second * 10

	if v, ok := os.LookupEnv("METERINGCOLLECTOR_GRPC_PORT"); ok && 0 != len(v) {
		if strings.Contains(v, ":") {
			m.grpcHost = v
		} else {
			m.grpcHost = "localhost:" + v
		}
	}
	if v, ok := os.LookupEnv("METERINGCOLLECTOR_HTTP_PORT"); ok && 0 != len(v) {
		if strings.Contains(v, ":") {
			m.httpHost = v
		} else {
			m.httpHost = ":" + v
		}
	}

	if storage != "" {
		if err := os.Setenv("METRICS_STORAGE_DRIVER", storage); err != nil {
			fmt.Println("Environment not set, error:", err)
		}
	}
	if v, ok := os.LookupEnv("METRICS_STORAGE_DRIVER"); ok && 0 != len(v) {
		if strings.HasPrefix(v, "elasticsearch=") {
			m.collectormanager.MetricsStorageDriver = v
		} else {
			fmt.Println("RPC driver not support,", v)
		}
	}

	m.run()
}

func (m *myExporter) run() {
	wg := sync.WaitGroup{}
	ch := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.startGRPC(ch)
	}()

	select {
	case <-ch:
		return
	case <-time.After(time.Millisecond * 500):
		fmt.Println("So gRPC running")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.startGateway()
	}()

	m.dispatchersignal = make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		/*
		   default to read from 'docker stats'?
		*/
		m.exportermanager.Dispatch(m.dispatchersignal)
	}()

	/*
	   https://github.com/kubernetes/kubernetes/blob/release-1.1/build/pause/pause.go
	*/
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Block until a signal is received.
		<-c

		// to stop stuff...
		m.dispatchersignal <- false
		goruntime.Goexit()
	}()

	wg.Wait()
}

func (m *myExporter) startGRPC(ch chan<- bool) {
	host := m.grpcHost

	s := grpc.NewServer()
	pb.RegisterExporterServiceServer(s, m)

	l, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start gRPC on host", l.Addr())
	if err := s.Serve(l); nil != err {
		panic(err)
	}
	ch <- false
}

func (m *myExporter) startGateway() {
	gRPCHost := m.grpcHost
	host := m.httpHost

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	// mux.HandleFunc("/swagger/", serveSwagger2)
	//	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
	//		io.Copy(w, strings.NewReader(healthcheckerpb.Swagger))
	//	})

	dopts := []grpc.DialOption{grpc.WithInsecure()}

	fmt.Println("Start gRPC Gateway into host", gRPCHost)
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterExporterServiceHandlerFromEndpoint(ctx, gwmux, gRPCHost, dopts); err != nil {
		fmt.Println("Failed to run HTTP server. ", err)
		return
	}

	mux.Handle("/", gwmux)
	// serveSwagger(mux)
	//	fmt.Printf("Start HTTP")
	//	if err := http.ListenAndServe(host, allowCORS(mux)); nil != err {
	//		fmt.Fprintf(os.Stderr, "Server died: %s\n", err)
	//	}

	lstn, err := net.Listen("tcp", host)
	if nil != err {
		panic(err)
	}

	fmt.Printf("http on host: %s\n", lstn.Addr())
	srv := &http.Server{
		Handler: func /*allowCORS*/ (h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if origin := r.Header.Get("Origin"); origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
						func /*preflightHandler*/ (w http.ResponseWriter, r *http.Request) {
							headers := []string{"Content-Type", "Accept"}
							w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
							methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
							w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
							glog.Infof("preflight request for %s", r.URL.Path)
							return
						}(w, r)
						return
					}
				}
				h.ServeHTTP(w, r)
			})
		}(mux),
	}

	if err := srv.Serve(lstn); nil != err {
		fmt.Fprintln(os.Stderr, "Server died.", err.Error())
	}
}

func (m *myCollector) run() {
	wg := sync.WaitGroup{}
	ch := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.startGRPC(ch)
	}()

	select {
	case <-ch:
		return
	case <-time.After(time.Millisecond * 500):
		fmt.Println("So gRPC running")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.startGateway()
	}()

	m.dispatchersignal = make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		/*
		   default to write to stdout?
		*/
		m.collectormanager.Start(m.dispatchersignal)
	}()

	/*
	   https://github.com/kubernetes/kubernetes/blob/release-1.1/build/pause/pause.go
	*/
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Block until a signal is received.
		<-c

		// to stop stuff...
		m.dispatchersignal <- false
		goruntime.Goexit()
	}()

	wg.Wait()
}

func (m *myCollector) startGRPC(ch chan<- bool) {
	host := m.grpcHost

	s := grpc.NewServer()
	pb.RegisterCollectorServiceServer(s, m)

	l, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start gRPC on host", l.Addr())
	if err := s.Serve(l); nil != err {
		panic(err)
	}
	ch <- false
}

func (m *myCollector) startGateway() {
	gRPCHost := m.grpcHost
	host := m.httpHost

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	// mux.HandleFunc("/swagger/", serveSwagger2)
	//	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
	//		io.Copy(w, strings.NewReader(healthcheckerpb.Swagger))
	//	})

	dopts := []grpc.DialOption{grpc.WithInsecure()}

	fmt.Println("Start gRPC Gateway into host", gRPCHost)
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterCollectorServiceHandlerFromEndpoint(ctx, gwmux, gRPCHost, dopts); err != nil {
		fmt.Println("Failed to run HTTP server. ", err)
		return
	}

	mux.Handle("/", gwmux)
	// serveSwagger(mux)
	//	fmt.Printf("Start HTTP")
	//	if err := http.ListenAndServe(host, allowCORS(mux)); nil != err {
	//		fmt.Fprintf(os.Stderr, "Server died: %s\n", err)
	//	}

	lstn, err := net.Listen("tcp", host)
	if nil != err {
		panic(err)
	}

	fmt.Printf("http on host: %s\n", lstn.Addr())
	srv := &http.Server{
		Handler: func /*allowCORS*/ (h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if origin := r.Header.Get("Origin"); origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
						func /*preflightHandler*/ (w http.ResponseWriter, r *http.Request) {
							headers := []string{"Content-Type", "Accept"}
							w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
							methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
							w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
							glog.Infof("preflight request for %s", r.URL.Path)
							return
						}(w, r)
						return
					}
				}
				h.ServeHTTP(w, r)
			})
		}(mux),
	}

	if err := srv.Serve(lstn); nil != err {
		fmt.Fprintln(os.Stderr, "Server died.", err.Error())
	}
}
