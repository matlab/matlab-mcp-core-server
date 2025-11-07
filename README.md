# MATLAB MCP Core Server

Run MATLAB® using AI applications. The MATLAB MCP Core Server allows your AI applications to:

- Start and quit MATLAB.
- Write and run MATLAB code.
- Assess your MATLAB code for style and correctness.
  
## Table of Contents
  - [Setup](#setup)
    - [Claude Code](#claude-code) 
    - [Claude Desktop](#claude-desktop)
    - [Cursor](#cursor)
    - [GitHub Copilot in Visual Studio Code](#github-copilot-in-visual-studio-code)
  - [Arguments](#arguments)
  - [Tools](#tools)
  - [Data Collection](#data-collection)
  - [Disclaimer](#disclaimer)

## Setup

1. Install [MATLAB (MathWorks)](https://www.mathworks.com/help/install/ug/install-products-with-internet-connection.html) 2020b or later and add it to the system PATH.
2. Download the [Latest Release](https://github.com/matlab/matlab-mcp-core-server/releases/latest) from GitHub®. Alternatively, you can install [Go](https://go.dev/doc/install) and build the binary from source using `go install github.com/matlab/matlab-mcp-core-server/cmd/matlab-mcp-core-server`. 
3. Add the MATLAB MCP Core Server to your AI application. You can find instructions for adding MCP servers in the documentation of your AI application. For example instructions on using Claude Code®, Claude Desktop®, Cursor, and GitHub Copilot in Visual Studio® Code, see below. Note that you can customize the server by specifying optional [arguments](#arguments).

#### Download on macOS

From a macOS terminal, use `curl` to download the latest version of the MATLAB MCP Core Server for your macOS architecture.

* macOS (Apple silicon processor)

   ```sh
   curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maca64
   ```

* macOS (Intel processor)

   ```sh
   curl -L -o ~/Downloads/matlab-mcp-core-server https://github.com/matlab/matlab-mcp-core-server/releases/latest/download/matlab-mcp-core-server-maci64
   ```

Give the downloaded file executable permissions so that you can run the MATLAB MCP Core Server:

```sh
chmod +x ~/Downloads/matlab-mcp-core-server
```

### Claude Code

In your terminal, run the following, remembering to insert the full path to the server binary you acquired in the setup:
```sh
claude mcp add --transport stdio matlab /fullpath/to/matlab-mcp-core-server-binary [arguments...]
```
You can customize the server by specifying [arguments](#arguments):
```sh
claude mcp add --transport stdio matlab /fullpath/to/matlab-mcp-core-server-binary --initial-working-folder=/home/username/myproject
```

For details on adding MCP servers in Claude Code, see [Add a local stdio server (Claude Code)](https://docs.claude.com/en/docs/claude-code/mcp#option-3%3A-add-a-local-stdio-server). To remove the server later, run:
```sh
claude mcp remove matlab
```

### Claude Desktop

Follow the instructions on the page [Connect to local MCP servers (MCP)](https://modelcontextprotocol.io/docs/develop/connect-local-servers) to install Node.js and the Filesystem Server. In your Claude Desktop configuration file, you need to add the configuration for the MATLAB MCP Core Server as well as the Filesystem Server. You can use the combined JSON below. In the Filesystem `args`, remember to specify which paths the server can access. In the MATLAB `args`, remember to insert the full path to the server binary you acquired, as well as any other arguments:

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
         "command": "fullpath\\to\\matlab-mcp-core-server-binary",
         "args": [
            "--initial-working-folder=C:\\Users\\username\\Documents"
         ]
      }
   }
}
```
After saving the configuration file, quit and restart Claude Desktop. 

### Cursor

Configure the MATLAB MCP Core Server in Cursor by creating an `mcp.json` file. You can create a project-specific configuration in `.cursor/mcp.json` within your project directory, or a global configuration in `~/.cursor/mcp.json` in your home directory. In the configuration file, add the following, remembering to insert the full path to the server binary you acquired in the setup, as well as any other arguments:

```json
{
   "mcpServers": {
      "matlab": {
         "command": "fullpath/to/matlab-mcp-core-server-binary",
         "args": [
            "--initial-working-folder=/home/username/Documents"
         ]
      }
   }
}
```

After saving the configuration file, restart Cursor to apply the changes.

For details on adding MCP servers in Cursor, see [Model Context Protocol (Cursor)](https://cursor.com/docs/context/mcp#model-context-protocol).

### GitHub Copilot in Visual Studio Code

VS Code provides different methods to [Add an MCP Server (VS Code)](https://code.visualstudio.com/docs/copilot/customization/mcp-servers?wt.md_id=AZ-MVP-5004796#_add-an-mcp-server). MathWorks recommends you follow the steps in the section **"Add an MCP server to a workspace `mcp.json` file"**. In your `mcp.json` configuration file, add the following, remembering to insert the full path to the server binary you acquired in the setup, as well as any arguments:
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
| initial-working-folder | Specify the folder where MATLAB starts and where the server generates any MATLAB scripts. If you do not provide the argument, MATLAB starts in these locations: <br><br> <ul><li>Linux: `/home/username` </li><li> Windows: `C:\Users\username\Documents`</li><li>Mac: `/Users/username/Documents`</li></ul> | `"--initial-working-folder=C:\\Users\\name\\MyProject"` |  
| disable-telemetry | To disable anonymized data collection, set this argument to `true`. For details, see [Data Collection](#data-collection). | `"--disable-telemetry=true"`  |

## Tools

1. `detect_matlab_toolboxes`
   - Lists installed MATLAB toolboxes with version information.
 
2. `check_matlab_code`
   - Performs static code analysis on a MATLAB script. Returns warnings about coding style, potential errors, deprecated functions, performance issues, and best practice violations. This is a non-destructive, read-only operation that helps identify code quality issues without executing the script.
   - Inputs:
     - `script_path` (string): Absolute path to the MATLAB script file to analyze. Must be a `.m` file within an allowed directory. The file is not modified during analysis. Example: `C:\Users\username\matlab\myFunction.m` or `/home/user/scripts/analysis.m`.
 
3. `evaluate_matlab_code`
   - Evaluates a string of MATLAB code and returns the output.
   - Inputs:
     - `code` (string): MATLAB code to evaluate.
     - `project_path` (string): Absolute path to an allowed project directory. MATLAB sets this directory as the current working folder. Example: `C:\Users\username\matlab-project` or `/home/user/research`.
 
4. `run_matlab_file`
   - Executes a MATLAB script and returns the output. The script must be a valid `.m file`.
   - Inputs:
     - `script_path` (string): Absolute path to the MATLAB script file to execute. Must be a valid `.m` file within an allowed directory. Example: `C:\Users\username\projects\analysis.m` or `/home/user/matlab/simulation.m`.
 
5. `run_matlab_test_file`
   - Executes a MATLAB test script and returns comprehensive test results. Designed specifically for MATLAB unit test files that follow MATLAB testing framework conventions.
   - Inputs:
     - `script_path` (string): Absolute path to the MATLAB test script file. Must be a valid `.m` file containing MATLAB unit tests, within an allowed directory. Example: `C:\Users\username\tests\testMyFunction.m` or `/home/user/matlab/tests/test_analysis.m`.

## Data Collection

The MATLAB MCP Core Server may collect fully anonymized information about your usage of the server and send it to MathWorks. This data collection helps MathWorks improve products and is on by default. To opt out of data collection, set the argument `--disable-telemetry` to `true`.

## Disclaimer

The MATLAB MCP Core Server is provided "as is" without warranties of any kind, expressed or implied. By using this server, you acknowledge and accept that you are solely responsible for any actions taken and any consequences arising from its use. It is your responsibility to thoroughly review and validate all tool calls before execution. The developers and providers of this server disclaim any liability for loss, damage, or injury resulting from its use.

---

Copyright 2025 The MathWorks, Inc.

----
