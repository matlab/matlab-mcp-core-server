classdef DefaultFSAdaptor < mcpserver.internal.fs.FSAdaptor
    %DefaultFSAdaptor Default implementation for filesystem operations with security features
    %   This class provides higher-level operations built on top of the
    %   filesystem facade and permission manager.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        FSFacade(1, 1) mcpserver.internal.facade.fs.FSFacade = mcpserver.internal.facade.fs.DefaultFSFacade()
        PermissionManager(1, 1) mcpserver.internal.fs.internal.permissionmanager.PermissionManager = ...
            mcpserver.internal.fs.internal.permissionmanager.createPermissionManager()
    end

    methods
        function obj = DefaultFSAdaptor(options)
            arguments
                options.?mcpserver.internal.fs.DefaultFSAdaptor
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function ensureSecureFolder(obj, folderPath)
            if obj.FSFacade.isfolder(folderPath)
                if ~obj.PermissionManager.checkPermissionsIsUserOnly(folderPath)
                    throw(mcpserver.internal.error.Errors.InsecurePermissions(folderPath));
                end
            else
                if obj.FSFacade.isfile(folderPath)
                    throw(mcpserver.internal.error.Errors.FileExistsAtFolderPath(folderPath));
                end
                [status, msg] = obj.FSFacade.mkdir(folderPath);
                if ~status
                    throw(mcpserver.internal.error.Errors.FailedToCreateDirectory(folderPath, msg));
                end
                obj.PermissionManager.setPermissionsToUserOnly(folderPath);
            end
        end

        function ensureSecureFile(obj, filePath)
            if obj.FSFacade.isfile(filePath)
                if ~obj.PermissionManager.checkPermissionsIsUserOnly(filePath)
                    throw(mcpserver.internal.error.Errors.InsecurePermissions(filePath));
                end
            else
                if obj.FSFacade.isfolder(filePath)
                    throw(mcpserver.internal.error.Errors.FolderExistsAtFilePath(filePath));
                end
                obj.FSFacade.writelines("", filePath);
                obj.PermissionManager.setPermissionsToUserOnly(filePath);
            end
        end
    end

end
