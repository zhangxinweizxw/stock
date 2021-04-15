package app

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"stock/share/lib"
	"stock/share/logging"
	"stock/share/middleware/cors"
	"stock/share/middleware/session"
	_ "stock/share/middleware/session/redis"
	"stock/share/store/mysql"
)

// Application initialization information
type App struct {
	Name                string
	Version             string
	EnvVarName          string
	EnvVarValue         string
	PidName             string
	WSPort              string
	MonitorPort         string
	LogAddr             string
	LogPort             string
	LogOn               bool
	IsDaemon            bool
	DisableGzip         bool
	CpuTimes            int
	SessionOn           bool
	SessionProviderName string
	SessionConfig       string
	Cors                []string
	Destroy             func()
}

func NewApp(name, version string) *App {
	return &App{
		Name:     name,
		Version:  version,
		IsDaemon: true,
	}
}

// Application Initialization
func (this *App) Init() *gin.Engine {
	if this.CpuTimes > 0 {
		runtime.GOMAXPROCS(runtime.NumCPU() * this.CpuTimes)
	}

	if this.IsDaemon {
		lib.RunAsDaemon(this.Name, this.Version, "_GO_DAEMON", "1", this.PidName)
	}

	// Check debug mode
	var debug, nc bool
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-debug" {
				debug = true
			}
			if arg == "-nc" {
				debug = true
				nc = true
			}
		}
	}

	logging.SetLogModel(debug, nc)
	localIPAddr, _ := lib.GetLocalIPAddr()

	logging.Info("---------------------------------------------")
	logging.Info("App Name:%s", this.PidName)
	logging.Info("Server started, version %s", this.Version)
	logging.Info("Start-up time: %s", time.Now().Format("2006-01-02 15:04:05"))
	logging.Info("Local Address: %s", localIPAddr)

	if this.LogOn {
		l := logging.NewQueueLogBackend()
		l.ListenLog(this.PidName, this.LogAddr, this.LogPort, this.LogOn)

		logging.Info("Remote Logging: %s:%s", this.LogAddr, this.LogPort)
	}

	logging.Info("%s - Listening", this.WSPort)
	logging.Info("--------------------------------------------")

	// Production Environment
	if nc {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	if len(this.Cors) > 0 {
		r.Use(cors.Default(this.Cors))
	}

	// GZip
	if !this.DisableGzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// Session
	if this.SessionOn {
		initSession(r, this.SessionProviderName, this.SessionConfig)
	}

	// Graceful close
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		<-signalChan

		if this.Destroy != nil {
			this.Destroy()
		}
		destroy()
	}()

	return r
}

// Init session
func initSession(engine *gin.Engine, providerName, sessionConfig string) {
	var globalSessions *session.Manager
	var err error

	if globalSessions, err = session.NewManager(providerName, sessionConfig); err != nil {
		logging.Fatal(err)
		os.Exit(1)
	}

	// init session middleware
	engine.Use(session.Middleware(globalSessions))

	go globalSessions.GC()
}

// When the application is closed for resource recovery
func destroy() {
	mysql.Close()
	//mongo.Close()

	os.Exit(1)
}
