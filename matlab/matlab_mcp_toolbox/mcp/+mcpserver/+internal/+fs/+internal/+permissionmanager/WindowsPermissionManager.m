classdef WindowsPermissionManager < mcpserver.internal.fs.internal.permissionmanager.PermissionManager
    %WindowsPermissionManager Permission manager for Windows.
    %   Sets exactly 3 ACEs (current user, SYSTEM, Administrators), all
    %   with FullControl, and protects the DACL from inheritance. SIDs
    %   are used for identity resolution to avoid domain-name ambiguity.
    %   This class is the only one that holds a WindowsACLManager and
    %   should never be constructed on non-Windows platforms.

    % Copyright 2026 The MathWorks, Inc.

    properties (Constant, Access = private)
        % Well-known SIDs matching C++ WinLocalSystemSid and WinBuiltinAdministratorsSid
        SYSTEM_SID = "S-1-5-18"
        ADMINISTRATORS_SID = "S-1-5-32-544"
    end

    properties (GetAccess = private, SetAccess = immutable)
        FSFacade(1, 1) mcpserver.internal.facade.fs.FSFacade = mcpserver.internal.facade.fs.DefaultFSFacade()
        WindowsACLManager(1, 1) mcpserver.internal.fs.internal.permissionmanager.internal.windowsacl.WindowsACLManager = mcpserver.internal.fs.internal.permissionmanager.internal.windowsacl.DefaultWindowsACLManager()
    end

    methods
        function obj = WindowsPermissionManager(options)
            arguments
                options.?mcpserver.internal.fs.internal.permissionmanager.WindowsPermissionManager
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function setPermissionsToUserOnly(obj, path)
            isDir = obj.FSFacade.isfolder(path);

            try
                userSid = obj.WindowsACLManager.getCurrentUserSID();
                sids = [userSid, obj.SYSTEM_SID, obj.ADMINISTRATORS_SID];
                obj.WindowsACLManager.setProtectedACL(path, sids, isDir);
            catch ME
                throw(addCause(mcpserver.internal.error.Errors.FailedToSetPermissions(path), ME));
            end
        end

        function tf = checkPermissionsIsUserOnly(obj, path)
            try
                userSid = obj.WindowsACLManager.getCurrentUserSID();
                allowedSids = obj.WindowsACLManager.getAllowedSIDs(path);
                daclProtected = obj.WindowsACLManager.isDACLProtected(path);
            catch ME
                throw(addCause(mcpserver.internal.error.Errors.FailedToGetFileAttributes(path), ME));
            end

            if ~daclProtected
                tf = false;
                return;
            end

            trustedSids = [userSid, obj.SYSTEM_SID, obj.ADMINISTRATORS_SID];

            % Verify all 3 trusted SIDs have access
            allTrustedHaveAccess = true;
            for i = 1:length(trustedSids)
                if ~any(allowedSids == trustedSids(i))
                    allTrustedHaveAccess = false;
                    break;
                end
            end

            % Verify no unexpected SIDs have access
            hasUnexpectedSid = false;
            for i = 1:length(allowedSids)
                if ~any(allowedSids(i) == trustedSids)
                    hasUnexpectedSid = true;
                    break;
                end
            end

            tf = allTrustedHaveAccess && ~hasUnexpectedSid;
        end
    end

end
