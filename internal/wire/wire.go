// Copyright 2025-2026 The MathWorks, Inc.

//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/lifecyclesignaler"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/orchestrator"
	files "github.com/matlab/matlab-mcp-core-server/internal/adaptors/filesystem/files"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/matlabrootselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/matlabstartingdirselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/logger"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession"
	localmatlabsessiondirectory "github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory/matlabfiles"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processdetails"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processlauncher"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabroot"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabversion"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/baseresource"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/configurator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/sdk"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	evalmatlabcodemultisessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	listavailablematlabstool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	startmatlabsessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	stopmatlabsessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	checkmatlabcodesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	detectmatlabtoolboxessinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabcodesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	runmatlabfilesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	runmatlabtestfilesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/messagecatalog"
	watchdogclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/iofacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/detectmatlabtoolboxes"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/utils/pathvalidator"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpserverfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/ossignaler"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/oswrapper"
	watchdogprocess "github.com/matlab/matlab-mcp-core-server/internal/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/processhandler"
	transportclient "github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/client"
	transportserver "github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

type Application struct {
	ModeSelector      *modeselector.ModeSelector
	MessageCatalog    *messagecatalog.MessageCatalog
	HTTPClientFactory *httpclientfactory.HTTPClientFactory
	HTTPServerFactory *httpserverfactory.HTTPServerFactory
}

func Initialize() *Application {
	wire.Build(
		// Application
		wire.Struct(new(Application), "*"),

		// Mode Selector
		modeselector.New,
		wire.Bind(new(modeselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(modeselector.Parser), new(*parser.Parser)),
		wire.Bind(new(modeselector.WatchdogProcess), new(*watchdogprocess.Watchdog)),
		wire.Bind(new(modeselector.Orchestrator), new(*orchestrator.Orchestrator)),
		wire.Bind(new(modeselector.OSLayer), new(*osfacade.OsFacade)),

		// Watchdog Process
		watchdogprocess.New,
		wire.Bind(new(watchdogprocess.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(watchdogprocess.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(watchdogprocess.ProcessHandler), new(*processhandler.ProcessHandler)),
		wire.Bind(new(watchdogprocess.OSSignaler), new(*ossignaler.OSSignaler)),
		wire.Bind(new(watchdogprocess.ServerHandlerFactory), new(*handler.Factory)),
		wire.Bind(new(watchdogprocess.ServerFactory), new(*transportserver.Factory)),
		wire.Bind(new(watchdogprocess.SocketFactory), new(*socket.Factory)),

		// Watchdog Transport Server Factory
		transportserver.NewFactory,
		wire.Bind(new(transportserver.HTTPServerFactory), new(*httpserverfactory.HTTPServerFactory)),
		wire.Bind(new(transportserver.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(transportserver.HandlerFactory), new(*handler.Factory)),

		// HTTP Server Factory
		httpserverfactory.New,
		wire.Bind(new(httpserverfactory.OSLayer), new(*osfacade.OsFacade)),

		// Orchestrator
		orchestrator.New,
		wire.Bind(new(orchestrator.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(orchestrator.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(orchestrator.Server), new(*server.Server)),
		wire.Bind(new(orchestrator.WatchdogClient), new(*watchdogclient.Watchdog)),
		wire.Bind(new(orchestrator.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(orchestrator.OSSignaler), new(*ossignaler.OSSignaler)),
		wire.Bind(new(orchestrator.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),
		wire.Bind(new(orchestrator.DirectoryFactory), new(*directory.Factory)),

		// MCP Server
		server.New,
		wire.Bind(new(server.MCPSDKServerFactory), new(*sdk.Factory)),
		wire.Bind(new(server.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(server.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(server.MCPServerConfigurator), new(*configurator.Configurator)),

		// MCP Server (SDK)
		sdk.NewFactory,
		wire.Bind(new(sdk.ConfigFactory), new(*config.Factory)),

		// MCP Server Configurator
		configurator.New,
		wire.Bind(new(configurator.ConfigFactory), new(*config.Factory)),

		// Tools
		wire.Bind(new(basetool.LoggerFactory), new(*logger.Factory)),

		listavailablematlabstool.New,
		wire.Bind(new(listavailablematlabstool.Usecase), new(*listavailablematlabs.Usecase)),

		listavailablematlabs.New,

		startmatlabsessiontool.New,
		wire.Bind(new(startmatlabsessiontool.Usecase), new(*startmatlabsession.Usecase)),

		startmatlabsession.New,

		stopmatlabsessiontool.New,
		wire.Bind(new(stopmatlabsessiontool.Usecase), new(*stopmatlabsession.Usecase)),

		stopmatlabsession.New,

		evalmatlabcodemultisessiontool.New,
		wire.Bind(new(evalmatlabcodemultisessiontool.Usecase), new(*evalmatlabcode.Usecase)),

		evalmatlabcodesinglesessiontool.New,
		wire.Bind(new(evalmatlabcodesinglesessiontool.Usecase), new(*evalmatlabcode.Usecase)),

		evalmatlabcode.New,
		wire.Bind(new(evalmatlabcode.PathValidator), new(*pathvalidator.PathValidator)),

		checkmatlabcodesinglesessiontool.New,
		wire.Bind(new(checkmatlabcodesinglesessiontool.Usecase), new(*checkmatlabcode.Usecase)),

		checkmatlabcode.New,
		wire.Bind(new(checkmatlabcode.PathValidator), new(*pathvalidator.PathValidator)),

		detectmatlabtoolboxessinglesessiontool.New,
		wire.Bind(new(detectmatlabtoolboxessinglesessiontool.Usecase), new(*detectmatlabtoolboxes.Usecase)),

		detectmatlabtoolboxes.New,

		runmatlabfilesinglesessiontool.New,
		wire.Bind(new(runmatlabfilesinglesessiontool.Usecase), new(*runmatlabfile.Usecase)),

		runmatlabfile.New,
		wire.Bind(new(runmatlabfile.PathValidator), new(*pathvalidator.PathValidator)),

		runmatlabtestfilesinglesessiontool.New,
		wire.Bind(new(runmatlabtestfilesinglesessiontool.Usecase), new(*runmatlabtestfile.Usecase)),

		runmatlabtestfile.New,
		wire.Bind(new(runmatlabtestfile.PathValidator), new(*pathvalidator.PathValidator)),

		// Resources
		wire.Bind(new(baseresource.LoggerFactory), new(*logger.Factory)),

		codingguidelines.New,
		plaintextlivecodegeneration.New,

		// Watchdog Client
		watchdogclient.New,
		wire.Bind(new(watchdogclient.WatchdogProcess), new(*process.Factory)),
		wire.Bind(new(watchdogclient.ClientFactory), new(*transportclient.Factory)),
		wire.Bind(new(watchdogclient.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(watchdogclient.SocketFactory), new(*socket.Factory)),

		// Watchdog Process Handler for Watchdog Client
		process.New,
		wire.Bind(new(process.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(process.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(process.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(process.ConfigFactory), new(*config.Factory)),

		// Watchdog Transport Client Factory
		transportclient.NewFactory,
		wire.Bind(new(transportclient.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(transportclient.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(transportclient.HTTPClientFactory), new(*httpclientfactory.HTTPClientFactory)),

		// Global MATLAB
		globalmatlab.New,
		wire.Bind(new(globalmatlab.MATLABManager), new(*matlabmanager.MATLABManager)),
		wire.Bind(new(globalmatlab.MATLABRootSelector), new(*matlabrootselector.MATLABRootSelector)),
		wire.Bind(new(globalmatlab.MATLABStartingDirSelector), new(*matlabstartingdirselector.MATLABStartingDirSelector)),

		// MATLAB Root Selector
		matlabrootselector.New,
		wire.Bind(new(matlabrootselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(matlabrootselector.MATLABManager), new(*matlabmanager.MATLABManager)),

		// MATLAB Starting Dir Selector
		matlabstartingdirselector.New,
		wire.Bind(new(matlabstartingdirselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(matlabstartingdirselector.OSLayer), new(*osfacade.OsFacade)),

		// Entities
		wire.Bind(new(entities.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),
		wire.Bind(new(entities.MATLABManager), new(*matlabmanager.MATLABManager)),

		// MATLAB Manager
		matlabmanager.New,
		wire.Bind(new(matlabmanager.MATLABServices), new(*matlabservices.MATLABServices)),
		wire.Bind(new(matlabmanager.MATLABSessionStore), new(*matlabsessionstore.Store)),
		wire.Bind(new(matlabmanager.MATLABSessionClientFactory), new(*matlabsessionclient.Factory)),

		// MATLAB Services
		matlabservices.New,
		wire.Bind(new(matlabservices.MATLABLocator), new(*matlablocator.MATLABLocator)),
		wire.Bind(new(matlabservices.LocalMATLABSessionLauncher), new(*localmatlabsession.Starter)),

		// MATLAB Locator
		matlablocator.New,
		wire.Bind(new(matlablocator.MATLABRootGetter), new(*matlabroot.Getter)),
		wire.Bind(new(matlablocator.MATLABVersionGetter), new(*matlabversion.Getter)),

		// MATLAB Root Getter
		matlabroot.New,
		wire.Bind(new(matlabroot.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(matlabroot.FileLayer), new(*filefacade.FileFacade)),

		// MATLAB Version Getter
		matlabversion.New,
		wire.Bind(new(matlabversion.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(matlabversion.IOLayer), new(*iofacade.IoFacade)),

		// Local MATLAB Session
		localmatlabsession.NewStarter,
		wire.Bind(new(localmatlabsession.SessionDirectoryFactory), new(*localmatlabsessiondirectory.Factory)),
		wire.Bind(new(localmatlabsession.ProcessDetails), new(*processdetails.ProcessDetails)),
		wire.Bind(new(localmatlabsession.MATLABProcessLauncher), new(*processlauncher.MATLABProcessLauncher)),
		wire.Bind(new(localmatlabsession.Watchdog), new(*watchdogclient.Watchdog)),

		// Local MATLAB Session Directory
		localmatlabsessiondirectory.NewFactory,
		wire.Bind(new(localmatlabsessiondirectory.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(localmatlabsessiondirectory.ApplicationDirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(localmatlabsessiondirectory.MATLABFiles), new(matlabfiles.MATLABFiles)),

		// MATLAB Files Provider
		matlabfiles.New,

		// Local MATLAB Session Process Details
		processdetails.New,
		wire.Bind(new(processdetails.OSLayer), new(*osfacade.OsFacade)),

		// Local MATLAB Process Launcher
		processlauncher.New,

		// MATLAB Session Store
		matlabsessionstore.New,
		wire.Bind(new(matlabsessionstore.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(matlabsessionstore.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),

		// MATLAB Session Client Factory
		matlabsessionclient.NewFactory,
		wire.Bind(new(matlabsessionclient.HttpClientFactory), new(*httpclientfactory.HTTPClientFactory)),

		// Shared Dependencies

		// Path Validator
		pathvalidator.New,
		wire.Bind(new(pathvalidator.OSLayer), new(*osfacade.OsFacade)),

		// Process Handler
		processhandler.New,
		wire.Bind(new(processhandler.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(processhandler.OSWrapper), new(*oswrapper.OSWrapper)),

		// OS Wrapper
		oswrapper.New,
		wire.Bind(new(oswrapper.OSLayer), new(*osfacade.OsFacade)),

		// HTTP Server Handler Factory
		handler.NewFactory,
		wire.Bind(new(handler.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(handler.ProcessHandler), new(*processhandler.ProcessHandler)),

		// Socket Factory
		socket.NewFactory,
		wire.Bind(new(socket.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(socket.OSLayer), new(*osfacade.OsFacade)),

		// Logger Factory
		logger.NewFactory,
		wire.Bind(new(logger.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(logger.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(logger.FilenameFactory), new(*files.Factory)),
		wire.Bind(new(logger.OSLayer), new(*osfacade.OsFacade)),

		// Directory Factory
		directory.NewFactory,
		wire.Bind(new(directory.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(directory.FilenameFactory), new(*files.Factory)),
		wire.Bind(new(directory.OSLayer), new(*osfacade.OsFacade)),

		// Lifecycle Signaler
		lifecyclesignaler.New,

		// Config Factory
		config.NewFactory,
		wire.Bind(new(config.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(config.Parser), new(*parser.Parser)),

		// Parser
		parser.New,
		wire.Bind(new(parser.MessageCatalog), new(*messagecatalog.MessageCatalog)),

		// Message Catalog
		messagecatalog.New,

		// Files Factory
		files.NewFactory,
		wire.Bind(new(files.OSLayer), new(*osfacade.OsFacade)),

		// HTTP Client Factory
		httpclientfactory.New,

		// OS Signaler
		ossignaler.New,

		// Facades
		osfacade.New,
		iofacade.New,
		filefacade.New,
	)

	return nil
}
