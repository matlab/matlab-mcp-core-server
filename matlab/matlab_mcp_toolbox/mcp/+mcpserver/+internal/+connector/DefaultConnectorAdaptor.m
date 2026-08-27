classdef DefaultConnectorAdaptor < mcpserver.internal.connector.ConnectorAdaptor
    %DefaultConnectorAdaptor Default implementation for MATLAB connector operations
    %   This class provides higher-level operations built on top of the
    %   connector facade.

    % Copyright 2026 The MathWorks, Inc.

    properties (GetAccess = private, SetAccess = immutable)
        ConnectorFacade(1, 1) mcpserver.internal.facade.connector.ConnectorFacade = mcpserver.internal.facade.connector.DefaultConnectorFacade()
        APIKeyManager(1, 1) mcpserver.internal.connector.internal.apikeymanager.APIKeyManager = mcpserver.internal.connector.internal.apikeymanager.DefaultAPIKeyManager()
        OSFacade(1, 1) mcpserver.internal.facade.os.OSFacade = mcpserver.internal.facade.os.DefaultOSFacade()
    end

    methods
        function obj = DefaultConnectorAdaptor(options)
            arguments
                options.?mcpserver.internal.connector.DefaultConnectorAdaptor
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function details = getConnectionDetails(obj)
            baseUrl = obj.ConnectorFacade.getBaseUrl();
            basePath = mcpserver.internal.connector.basePathFromBaseUrl(baseUrl);

            details = struct( ...
                "port", obj.ConnectorFacade.securePort(), ...
                "certificate", obj.ConnectorFacade.getCertificateLocation(), ...
                "apiKey", obj.APIKeyManager.getAPIKey(), ...
                "pid", obj.OSFacade.feature("getpid"), ...
                "basePath", basePath ...
            );
        end
    end

end
