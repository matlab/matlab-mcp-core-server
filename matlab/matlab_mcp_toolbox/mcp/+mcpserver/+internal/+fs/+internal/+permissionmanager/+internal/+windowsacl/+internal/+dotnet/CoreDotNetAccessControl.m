classdef CoreDotNetAccessControl < handle & mcpserver.internal.fs.internal.permissionmanager.internal.windowsacl.internal.dotnet.DotNetAccessControl
    %CoreDotNetAccessControl .NET Core implementation
    %   Uses FileSystemAclExtensions static methods which are required
    %   on .NET Core where SetAccessControl is no longer an instance method.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        DotNetFacade(1, 1) mcpserver.internal.facade.dotnet.DotNetFacade = ...
            mcpserver.internal.facade.dotnet.DefaultDotNetFacade()
    end

    properties (Access = private)
        AssemblyLoaded(1, 1) logical = false
    end

    properties (Constant, Access = private)
        AssemblyName = "System.IO.FileSystem.AccessControl"
    end

    methods
        function obj = CoreDotNetAccessControl(options)
            arguments
                options.DotNetFacade(1, 1) mcpserver.internal.facade.dotnet.DotNetFacade = ...
                    mcpserver.internal.facade.dotnet.DefaultDotNetFacade()
            end

            obj.DotNetFacade = options.DotNetFacade;
        end

        function setAccessControl(obj, target, security)
            obj.ensureAssemblyLoaded();
            obj.DotNetFacade.setFileSystemAccessControl(target, security);
        end
    end

    methods (Access = private)
        function ensureAssemblyLoaded(obj)
            if obj.AssemblyLoaded
                return;
            end

            try
                obj.DotNetFacade.addAssembly(obj.AssemblyName);
            catch cause
                err = mcpserver.internal.error.Errors.FailedToLoadDotNetAssembly(obj.AssemblyName);
                err = err.addCause(cause);
                throw(err);
            end
            obj.AssemblyLoaded = true;
        end
    end

end
