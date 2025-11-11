package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"logger/config"
	"logger/models/migrations"
	"logger/remotes/blockchain"
	"logger/remotes/postgres"
	"logger/web/router"
	"logger/web/server"

	"github.com/gorilla/mux"
	"github.com/joaopandolfi/blackwhale/handlers"
	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
	"github.com/opentracing/opentracing-go"

	"github.com/joaopandolfi/blackwhale/utils"
)

var tracerCloser io.Closer

func configInit() {
	config.Load(os.Args[1:])

	if config.Get().SnakeByDefault {
		handlers.ActiveSnakeCase()
	}

	// Init tracing
	tracer, closer := jaeger.Init(config.Get().SystemID)
	tracerCloser = closer
	opentracing.SetGlobalTracer(tracer)

	blockChainConfig := config.Get().BlockChain
	blockchain.InitChain(blockChainConfig.PubKey)
	blockchain.Get().SetAuth(blockChainConfig.PrivKey, blockChainConfig.Passphrase)

	postgres.Init(config.Get())

	migrations.Migrate()
	err := migrations.Terraform()
	if err != nil {
		utils.Error("Terraforming error", err.Error())
	}
}

func gracefullShutdown() {
	fmt.Println("<====================================Shutdown==================================>")
	if tracerCloser != nil {
		tracerCloser.Close()
	}
	// postgres.Close()
}

func welcome() {
	//https://patorjk.com/software/taag/#p=display&f=Doom&t=logger
	fmt.Println(`
     _                             
    | |                            
    | | ___   __ _  __ _  ___ _ __ 
    | |/ _ \ / _  |/ _  |/ _ \ '__|
    | | (_) | (_| | (_| |  __/ |   
    |_|\___/ \__, |\__, |\___|_|   
              __/ | __/ |          
             |___/ |___/           
		    
O======================================(Nothing will be forgotten)====>
	`)
}

func main() {
	welcome()
	//Init
	configInit()

	// Initialize Mux Router
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	srv := server.New(r, config.Get())
	nr := router.New(srv)
	nr.Setup()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go srv.Start()

	<-done
	utils.Info("[SERVER] Gracefully shutdown")
	gracefullShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		utils.CriticalError("Server Shutdown Failed", err.Error())
	}

	cancel()

}
