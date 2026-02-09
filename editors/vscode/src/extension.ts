import { workspace, ExtensionContext, window } from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  const command = workspace
    .getConfiguration("langz")
    .get<string>("serverPath", "langz");

  const serverOptions: ServerOptions = {
    command: command!,
    args: ["lsp"],
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "langz" }],
    outputChannel: window.createOutputChannel("Langz Language Server"),
  };

  client = new LanguageClient(
    "langz",
    "Langz Language Server",
    serverOptions,
    clientOptions
  );

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
