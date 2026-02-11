import { useState, useCallback } from "react";
import type { HistoryLine, LineType } from "../types";
import { commands, getPrompt } from "../commands/registry";

export const useTerminal = () => {
  const [history, setHistory] = useState<HistoryLine[]>([
    {
      id: "init",
      type: "info",
      content: "Welcome to t-learn v1.0.0. Type 'help' to start.",
    },
  ]);

  // New State for Command History
  const [commandHistory, setCommandHistory] = useState<string[]>([]);

  const [promptLabel, setPromptLabel] = useState(getPrompt());

  const pushToHistory = (content: string, type: LineType = "info") => {
    setHistory((prev) => [
      ...prev,
      { id: Date.now().toString() + Math.random(), type, content },
    ]);
  };

  const execute = useCallback(async (commandString: string) => {
    if (!commandString.trim()) return;

    // Save to Command History
    setCommandHistory((prev) => [...prev, commandString]);

    pushToHistory(commandString, "command");

    // Regex parsing and command execution logic stays the same
    const parts = commandString.match(/(?:[^\s"]+|"[^"]*")+/g);
    if (!parts) return;

    const cmdName = parts[0].toLowerCase();
    const args = parts.slice(1).map((arg) => {
      if (arg.startsWith('"') && arg.endsWith('"')) return arg.slice(1, -1);
      return arg;
    });

    if (cmdName === "clear") {
      setHistory([]);
      return;
    }

    const commandDef = commands[cmdName];
    if (commandDef) {
      try {
        const response = await commandDef.execute(args);
        pushToHistory(response.output, response.type);
      } catch (err: any) {
        pushToHistory(`Error executing '${cmdName}': ${err.message}`, "error");
      }
    } else {
      pushToHistory(
        `Command not found: ${cmdName}. Type 'help' for list.`,
        "error",
      );
    }

    setPromptLabel(getPrompt());
  }, []);

  return {
    history,
    execute,
    promptLabel,
    commandHistory,
    clear: () => setHistory([]),
  };
};
