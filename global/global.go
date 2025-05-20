package global

import (
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


var (
	Logger        	*zap.Logger
	FilteredData 	= make(chan []byte, 30)
	ParsedData 		= make(chan []byte, 30)
	Clients 		= make(map[*websocket.Conn]bool)
	MuClient       	sync.Mutex
	ParsedDataOneGeo  = make(chan []byte, 30)
	ClientsOneGeo  = make(map[*websocket.Conn]bool)
	MuClientOneGeo    sync.Mutex
	



	NUMBER_OF_RECEIVED_MESSAGES = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "NUMBER_OF_RECEIVED_MESSAGES",
		Help: "The total number of received messages",
	})
	NUMBER_OF_PARSED_MESSAGES = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "NUMBER_OF_PARSED_MESSAGES",
		Help: "The total number of parsed messages",
	})
	NUMBER_OF_FAILED_MESSAGES = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "NUMBER_OF_FAILED_MESSAGES",
		Help: "The total number of failed messages",
	})
	NUMBER_OF_SENT_MESSAGES = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "NUMBER_OF_SENT_MESSAGES",
		Help: "The total number of sent messages",
	})
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},}


	func init() {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	
		// Create a write syncer
		fileWriter := zapcore.AddSync(logFile)
	
		// Create a core that writes logs to the file with JSON encoding
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), // Use default production encoder
			fileWriter,         // Write to file
			zapcore.DebugLevel, // Log all levels
		)
	
		// Create logger with caller and stacktrace for errors
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	
		// Set as global logger
		zap.ReplaceGlobals(Logger)
		prometheus.MustRegister(NUMBER_OF_RECEIVED_MESSAGES)
		prometheus.MustRegister(NUMBER_OF_PARSED_MESSAGES)
		prometheus.MustRegister(NUMBER_OF_FAILED_MESSAGES)
		prometheus.MustRegister(NUMBER_OF_SENT_MESSAGES)
	}