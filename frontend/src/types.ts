export type LineType = "command" | "output" | "error" | "info" | "success";

export interface HistoryLine {
  id: string;
  type: LineType;
  content: string;
}

export interface CommandResponse {
  output: string;
  type: LineType;
}

// Every command function (login, help, courses) must follow this signature
export type CommandFunction = (args: string[]) => Promise<CommandResponse>;

export interface CommandDefinition {
  description: string;
  execute: CommandFunction;
}
