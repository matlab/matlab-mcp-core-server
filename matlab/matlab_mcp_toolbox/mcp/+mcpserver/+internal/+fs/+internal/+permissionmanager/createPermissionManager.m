function mgr = createPermissionManager(options)
    %createPermissionManager Returns the platform-appropriate PermissionManager.
    %   On Windows, returns a WindowsPermissionManager (which holds a
    %   WindowsACLManager and the .NET stack). On Unix and macOS, returns
    %   a UnixPermissionManager (which has no reference to any Windows
    %   class). This is the single place that decides which platform
    %   implementation to load.

    % Copyright 2026 The MathWorks, Inc.

    arguments
        options.OSFacade(1, 1) mcpserver.internal.facade.os.OSFacade = ...
            mcpserver.internal.facade.os.DefaultOSFacade()
    end

    if options.OSFacade.ispc()
        mgr = mcpserver.internal.fs.internal.permissionmanager.WindowsPermissionManager();
    else
        mgr = mcpserver.internal.fs.internal.permissionmanager.UnixPermissionManager();
    end
end
