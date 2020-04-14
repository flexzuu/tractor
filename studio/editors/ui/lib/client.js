
const AgentEndpoint = "ws://localhost:3001/";
const RetryInterval = 500;

function scheduleRetry(fn) {
    setTimeout(fn, RetryInterval);
}

class Session {
    constructor(workspacePath, onconnect) {
        this.workspacePath = workspacePath;
        this.onconnect = onconnect;
        this.api = new qrpc.API();
        this.discover();
    }

    async discover() {
        try {
            var conn = await qmux.DialWebsocket(AgentEndpoint);
        } catch (e) {
            scheduleRetry(() => this.discover());
            return;
        }
        var session = new qmux.Session(conn);
        var client = new qrpc.Client(session);
        var resp = await client.call("list");
        var workspace = resp.reply.find((el) => el.Path == this.workspacePath);
        if (workspace) {
            this.workspaceEndpoint = workspace.Endpoint;
            this.connect(workspace.Endpoint);
        } else {
            console.error("workspace path not a known workspace");
        }
        conn.close();
    }

    async connect(endpoint) {
        try {
            var conn = await qmux.DialWebsocket(endpoint);
        } catch (e) {
            scheduleRetry(() => this.connect(endpoint));
            return;
        }
        this.session = new qmux.Session(conn);
        this.client = new qrpc.Client(this.session, this.api);
        this.client.serveAPI();
        this.onconnect(this.client);
    }

    reconnect() {
        if (this.workspaceEndpoint) {
            scheduleRetry(() => this.connect(this.workspaceEndpoint));
        } else {
            scheduleRetry(() => this.discover());
        }
    }
}

export { Session };