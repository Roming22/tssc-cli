# `tsf`: Model Context Protocol Server (MCP)

## Usage

Configure your agentic LLM client to run the following command, and use STDIO for communication with the server.

```sh
tsf mcp-server --image="ghcr.io/redhat-appstudio/tsf:latest"
```

Note: if you use Cursor, you can use the [Cursor MCP configuration](../.cursor/mcp.json)

Then use the [recommended prompt](#prompt) to get started, and follow the instructions.

## Instructions

When the MCP server is initialized, it will provide instructions for the LLM client to follow. The instructions are documented [here](../pkg/mcpserver/instructions.md).

The instructions provide a clear overview of what the MCP server is designed to do, how it works, and the expected sequence of tools to follow.

## Prompt

This is the recommended prompt to instruct the AI agent to start a deployment.

```text
You are a senior Red Hat OpenShift Platform Engineer. You leverage a specialized OpenShift-based platform focused on Trusted Software Factory (TSF), primarily consisting of:

- Advanced Cluster Security (ACS)
- OpenShift GitOps
- Trusted Artifact Signer (TAS)
- OpenShift Pipelines
- Trusted Profile Analyzer (TPA)
- Red Hat Developer Hub (DH)
- Quay and other container registries

Your primary function is to provide actionable commands: YAML manifests and `oc`/`kubectl` CLI commands, all tailored for OpenShift 4.17+. You prioritize security-hardened configurations and adherence to CNCF best practices.

Your interaction schema is built around Model Context Protocol (MCP) tools, enabling you to interact with various OpenShift work as a platform engineer. Your responses are structured to include a quick command block for immediate use.

Now let's get started on deploying TSF.
```

## Tools

The following MCP tools are exposed by the MCP server.

### Configuration

This section covers for the features in `tsf config` subcommand.

#### `tsf_config_get`

- *Description*: Get the existing TSF configuration in the cluster, or return the default if none exists yet. Use the default configuration as the reference to create a new TSF configuration for the cluster.
- *Arguments*: None.

#### `tsf_config_create`

- *Description*: Create a new TSF configuration in the cluster, in case none exists yet. Use the defaults as the reference to create a new TSF cluster configuration.
- *Arguments*:
    - **namespace** (string):
        - **Description**: The main namespace for TSF (`.tssc.namespace`), where Red Hat Developer Hub (DH) and other fundamental services will be deployed.
        - **Default**: "tssc".
    - **setting** (object):
        - **Description**: The global settings object for TSF (`.tssc.settings{}`). When empty the default settings will be used.

### Integration

Integrations with external services are managed via the `tssc integration <integration-name>` command. Since these integrations often require sensitive information, credentials, the MCP server will not handle them directly. Instead, it will provide the user with the exact `tssc integration` commands needed to configure the integration.

#### `tssc_integration_list`

- *Description*: List the TSF integrations available for the user. Certain integrations are required for certain features, make sure to configure the integrations accordingly.
- *Arguments*: None.

### Deploy

To handle the long-running deployment process without blocking the MCP server, the MCP server will delegate tasks to the cluster.

The `tssc` design maintains a single deployment per cluster by centralizing the installer configuration in a unique `ConfigMap`, ensuring a single RHADS installation per cluster.

Following the same principle, the MCP server generates a Kubernetes Job to run the `tssc` container image (specifically the `tssc deploy` subcommand) which proceeds with installing the predefined sequence of Helm charts that makes up RHADS.

#### `tssc_deploy_status`

- *Description*: Reports the status of the TSF deploy Job running in the cluster.
- *Arguments*: None.

#### `tssc_deploy`

- *Description*: Deploys TSF components to the cluster, uses the cluster configuration to deploy the TSF components sequentially.
- *Arguments*: None.
