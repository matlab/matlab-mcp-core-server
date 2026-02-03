# MATLAB MCP Core Server

Run MATLAB® using AI applications with the official MATLAB MCP Server from MathWorks®. The MATLAB MCP Core Server allows your AI applications to:

- Start and quit MATLAB.
- Write and run MATLAB code.
- Assess your MATLAB code for style and correctness.

## Table of Contents

- [Setup](#setup)
  - [Claude Code](#claude-code)
  - [Claude Desktop](#claude-desktop)
  - [GitHub Copilot in Visual Studio Code](#github-copilot-in-visual-studio-code)
- [Arguments](#arguments)
- [Tools](#tools)
- [Resources](#resources)
- [Data Collection](#data-collection)

## Setup

1. Install [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) 2020b or later and add it to the system PATH.
1. For Windows or Linux, download the [Latest Release](https://github.com/matlab/matlab-mcp-core-server/releases/latest) from GitHub®. Alternatively, you can install [Go](https://go.dev/doc/install) and build the binary from source using

    ```sh
    go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server@latest
    ```

    For macOS, first download the latest release by running the following command in your terminal:
    - For Apple silicon processors, run:
        ```sh
        curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
        ```
    - For Intel processors, run:
        ```sh
        curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
        ```
    Then grant executable permissions to the downloaded binary so you can run the MATLAB MCP Core Server:

    ```sh
    chmod +x ~/Downloads/matlab-mcp-core-server
    ```

1. Add the MATLAB MCP Core Server to your AI application. You can find instructions for adding MCP servers in the documentation of your AI application. For example instructions on using Claude Code®, Claude Desktop®, and GitHub Copilot in Visual Studio® Code, see below. Note that you can customize the server by specifying optional [Arguments](#arguments).

### Claude Code

In your terminal, run the following, remembering to insert the full path to the server binary you acquired in the setup:

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary
```

You can customize the server by specifying optional [Arguments](#arguments). Note the `--` separator between Claude Code's options and the server arguments:

```sh
claude mcp add --transport stdio matlab -- /fullpath/to/matlab-mcp-core-server-binary --initial-working-folder=/home/username/myproject
```

For details on adding MCP servers in Claude Code, see [Add a local stdio server (Claude Code)](https://docs.claude.com/en/docs/claude-code/mcp#option-3%3A-add-a-local-stdio-server). To remove the server later, run:

```sh
claude mcp remove matlab
```

### Claude Desktop

Follow the instructions on the page [Connect to local MCP servers (MCP)](https://modelcontextprotocol.io/docs/develop/connect-local-servers) to install Node.js and the Filesystem Server. These are required to allow Claude to create files on your filesystem that MATLAB can access. In your Claude Desktop configuration file, you need to add the configuration for the MATLAB MCP Core Server as well as the Filesystem Server. You can use the combined JSON below. In the Filesystem `args`, remember to specify which paths the server can access. In the MATLAB `args`, remember to insert the full path to the server binary you acquired, as well as any other [Arguments](#arguments). (Note that on Windows, your paths require extra backslashes as escape characters).

```json
{
    "mcpServers": {
        "filesystem": {
            "command": "npx",
            "args": [
                "-y",
                "@modelcontextprotocol/server-filesystem",
                "C:\\Users\\username"
            ]
        },
        "matlab": {
            "command": "C:\\fullpath\\to\\matlab-mcp-core-server-binary",
            "args": ["--initial-working-folder=C:\\Users\\username\\Documents"]
        }
    }
}
```

After saving the configuration file, quit Claude Desktop by clicking **File > Exit**, then restart Claude Desktop.

### GitHub Copilot in Visual Studio Code

VS Code provides different methods to [Add an MCP Server (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers?wt.md_id=AZ-MVP-5004796#_add-an-mcp-server). MathWorks recommends you follow the steps in the section **"Add an MCP server to a workspace `mcp.json` file"**. In your `mcp.json` configuration file, add the following, remembering to insert the full path to the server binary you acquired in the setup, as well as any [Arguments](#arguments):

```json
{
    "servers": {
        "matlab": {
            "type": "stdio",
            "command": "/fullpath/to/matlab-mcp-core-server-binary",
            "args": []
        }
    }
}
```

## Arguments

Customize the behavior of the server by providing arguments in the `args` array when configuring your AI application.

| Argument | Description | Example |
| ------------- | ------------- | ------------- |
| matlab-root | Full path specifying which MATLAB to start. Do not include `/bin` in the path. By default, the server tries to find the first MATLAB on the system PATH. | `"--matlab-root=/home/usr/MATLAB/R2025a"` |
| initialize-matlab-on-startup | To initialize MATLAB as soon as you start the server, set this argument to `true`. By default, MATLAB only starts when the first tool is called. | `"--initialize-matlab-on-startup=true"` |
| initial-working-folder | Specify the folder where MATLAB starts. If you do not provide the argument, MATLAB starts in these locations: <br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | `"--initial-working-folder=C:\\Users\\name\\MyProject"` |
| matlab-display-mode | Specify whether to show the MATLAB desktop. Use `desktop` mode (default) to show the MATLAB desktop. Use `nodesktop` mode to use MATLAB only from your AI application, without the MATLAB desktop. Note that in `nodesktop` mode, commands requiring a graphical interface (such as `edit`, `open`, `open_system`, `uifigure`, and `appdesigner`) will still open MATLAB windows on your desktop.| `"--matlab-display-mode=nodesktop"` |
| disable-telemetry | To disable anonymized data collection, set this argument to `true`. For details, see [Data Collection](#data-collection). | `"--disable-telemetry=true"` |


## Tools

1. `detect_matlab_toolboxes`
    - Lists installed MATLAB toolboxes with version information.

1. `check_matlab_code`
    - Performs static code analysis on a MATLAB script. Returns warnings about coding style, potential errors, deprecated functions, performance issues, and best practice violations. This is a non-destructive, read-only operation that helps identify code quality issues without executing the script.
    - Inputs:
        - `script_path` (string): Absolute path to the MATLAB script file to analyze. Must be a valid `.m` file. The file is not modified during analysis. Example: `C:\Users\username\matlab\myFunction.m` or `/home/user/scripts/analysis.m`.

1. `evaluate_matlab_code`
    - Evaluates a string of MATLAB code and returns the output.
    - Inputs:
        - `code` (string): MATLAB code to evaluate.
        - `project_path` (string): Absolute path to your project directory. MATLAB sets this directory as the current working folder. Example: `C:\Users\username\matlab-project` or `/home/user/research`.

1. `run_matlab_file`
    - Executes a MATLAB script and returns the output. The script must be a valid `.m file`.
    - Inputs:
        - `script_path` (string): Absolute path to the MATLAB script file to execute. Must be a valid `.m` file. Example: `C:\Users\username\projects\analysis.m` or `/home/user/matlab/simulation.m`.

1. `run_matlab_test_file`
    - Executes a MATLAB test script and returns comprehensive test results. Designed specifically for MATLAB unit test files that follow MATLAB testing framework conventions.
    - Inputs:
        - `script_path` (string): Absolute path to the MATLAB test script file. Must be a valid `.m` file containing MATLAB unit tests. Example: `C:\Users\username\tests\testMyFunction.m` or `/home/user/matlab/tests/test_analysis.m`.

## Resources

The MCP server provides [Resources (MCP)](https://modelcontextprotocol.io/specification/2025-03-26/server/resources) to help your AI application write MATLAB code. To see instructions for using this resource, refer to the documentation of your AI application that explains how to use resources.

1. `matlab_coding_guidelines`
    - Provides comprehensive MATLAB coding standards for improving code readability, maintainability, and collaboration. The guidelines encompass naming conventions, formatting, commenting, performance optimization, and error handling.
    - URI: `guidelines://coding`
    - MIME Type: `text/markdown`
    - Source: [MATLAB Coding Standards (GitHub)](https://github.com/matlab/rules/blob/main/matlab-coding-standards.md)

1. `plain_text_live_code_guidelines`
    - Provides rules and guidelines for generating live scripts using the plain text Live Code `.m` file format, suitable for version control and AI-assisted development. Note that to run plain text live scripts you need MATLAB R2025a or newer. For details, see [Live Code File Format (MathWorks)](https://www.mathworks.com/help/matlab/matlab_prog/plain-text-file-format-for-live-scripts.html).
    - URI: `guidelines://plain-text-live-code`
    - MIME Type: `text/markdown`
    - Source: [Plain Text Live Code Generation (GitHub)](https://github.com/matlab/rules/blob/main/live-script-generation.md)

## Data Collection

The MATLAB MCP Core Server may collect fully anonymized information about your usage of the server and send it to MathWorks. This data collection helps MathWorks improve products and is on by default. To opt out of data collection, set the argument `--disable-telemetry` to `true`.

#

When using the MATLAB MCP Core Server, you should thoroughly review and validate all tool calls before you run them. Always keep a human in the loop for important actions and only proceed once you are confident the call will do exactly what you expect. For more information, see [User Interaction Model (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/server/tools#user-interaction-model) and [Security Considerations (MCP)](https://modelcontextprotocol.io/specification/2025-06-18/server/tools#security-considerations).

The MATLAB MCP Core server may only be used with MATLAB installations that are used as a Personal Automation Server. Use with a central Automation Server is not allowed. Please contact MathWorks if Automation Server use is required. For more information see the [Program Offering Guide (MathWorks)](https://www.mathworks.com/help//pdf_doc/offering/offering.pdf).

---

Copyright 2025-2026 The MathWorks, Inc.

---
