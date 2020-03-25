package platform

import "fmt"

// New returns a new application from the given configuration.
func New(config Config) Platform {
	p := Platform{}
	p.Config = config

	p.setupLogger()
	p.Logger.Debug("[INIT] Logger initialized.")

	// TODO: Decide where the score generation should go

	p.setupContext()
	p.Logger.Debug("[INIT] Context initialized.")

	p.setupServer()
	p.Logger.Debug(fmt.Sprintf("[INIT] Server initialized using HTTP port [%s] and SSL port [%s].", p.Server.HTTPPort, p.Server.SSLport))
	p.Logger.Debug(fmt.Sprintf("[INIT] Database [%s] initialized", p.Server.DbConfig.Name))

	p.setupRouter()
	p.Logger.Debug("[INIT] Router initialized.")

	p.setupValidator() // TODO: Decide where should the custom validators should go
	p.Logger.Debug("[INIT] Validators initialized.")

	p.setupFormDecoder()
	p.Logger.Debug("[INIT] Form decoder initialized.")

	p.setupPermissions()
	p.Logger.Debug("[INIT] Permissions initialized.")

	p.setupUserService()
	p.Logger.Debug("[INIT] User service initialized")

	p.setupDatabase()
	p.Logger.Debug("[INIT] Database configured: Migration, default data and custom indexes.")

	p.setupCloudProvider()
	p.Logger.Debug("[INIT] Cloud provider initialized: AWS.")

	p.setupOrchestrator()
	p.Logger.Debug("[INIT] Orchestrator initialized: Kubernetes.")

	// TODO: Move http handler instance logic to applications

	// TODO: Initialize applications? We might need to register them instead.

	p.setupNodeManager()
	p.Logger.Debug("[INIT] NodeManager initialized. Using: AWS and Kubernetes.")

	p.setupScheduler()
	p.Logger.Debug("[INIT] Scheduler initialized.")

	p.setupQueues()
	p.Logger.Debug("[INIT] Launch and termination queues have been initialized.")

	p.setupWorkers()
	p.Logger.Debug("[INIT] Launch and termination workers have been initialized.")
	return p
}