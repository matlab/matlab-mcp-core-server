classdef UnixPermissionManager < mcpserver.internal.fs.internal.permissionmanager.PermissionManager
    %UnixPermissionManager Permission manager for Unix and macOS.
    %   Uses chmod 700 for directories and 600 for files. Holds no
    %   reference to Windows-only classes, so the load chain on
    %   Linux/macOS never touches +windowsacl or +dotnet.

    % Copyright 2026 The MathWorks, Inc.

    properties (Constant, Access = private)
        CHMOD_FOLDER_PERMISSIONS = "700"
        CHMOD_FILE_PERMISSIONS = "600"
    end

    properties (GetAccess = private, SetAccess = immutable)
        FSFacade(1, 1) mcpserver.internal.facade.fs.FSFacade = mcpserver.internal.facade.fs.DefaultFSFacade()
        OSFacade(1, 1) mcpserver.internal.facade.os.OSFacade = mcpserver.internal.facade.os.DefaultOSFacade()
    end

    methods
        function obj = UnixPermissionManager(options)
            arguments
                options.?mcpserver.internal.fs.internal.permissionmanager.UnixPermissionManager
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function setPermissionsToUserOnly(obj, path)
            [status, attribs, ~] = obj.FSFacade.fileattrib(path);
            if ~status
                throw(mcpserver.internal.error.Errors.FailedToGetFileAttributes(path));
            end

            if attribs.directory
                mode = obj.CHMOD_FOLDER_PERMISSIONS;  % rwx------
            else
                mode = obj.CHMOD_FILE_PERMISSIONS;  % rw-------
            end

            escapedPath = strrep(path, "'", "'\''");
            cmd = sprintf("chmod %s '%s'", mode, escapedPath);
            [status, ~] = obj.OSFacade.system(cmd);

            if status ~= 0
                throw(mcpserver.internal.error.Errors.FailedToSetPermissions(path));
            end
        end

        function tf = checkPermissionsIsUserOnly(obj, path)
            [status, attribs, ~] = obj.FSFacade.fileattrib(path);
            if ~status
                throw(mcpserver.internal.error.Errors.FailedToGetFileAttributes(path));
            end

            groupHasAccess = attribs.GroupRead || attribs.GroupWrite || attribs.GroupExecute;
            otherHasAccess = attribs.OtherRead || attribs.OtherWrite || attribs.OtherExecute;
            userHasRead = attribs.UserRead;
            userHasWrite = attribs.UserWrite;
            userHasExecute = attribs.UserExecute;

            % Directories: rwx------ (700)
            % Files: rw------- (600)
            if attribs.directory
                tf = userHasRead && userHasWrite && userHasExecute && ~groupHasAccess && ~otherHasAccess;
            else
                tf = userHasRead && userHasWrite && ~userHasExecute && ~groupHasAccess && ~otherHasAccess;
            end
        end
    end

end
