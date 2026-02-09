import { useState, useCallback } from "react";
import type { HistoryLine, LineType } from "../types";
import { commands } from "../commands/registry";

export const useTerminal = () => {
  const [history, setHistory] = useState<HistoryLine[]>([
    {
      id: "init",
      type: "info",
      content: "Welcome to t-learn v1.0.0. Type 'help' to start.",
    },
  ]);

  const pushToHistory = (content: string, type: LineType = "info") => {
    setHistory((prev) => [
      ...prev,
      { id: Date.now().toString() + Math.random(), type, content },
    ]);
  };

  const execute = useCallback(async (commandString: string) => {
    if (!commandString.trim()) return;

    // Echo the user's command to the screen
    pushToHistory(commandString, "command");

    // Parse: Split "login admin 123" -> cmd="login", args=["admin", "123"]
    const parts = commandString.trim().split(/\s+/);
    const cmdName = parts[0].toLowerCase();
    const args = parts.slice(1);

    // Special Case: Clear
    if (cmdName === "clear") {
      setHistory([]);
      return;
    }

    // Find & Execute Command
    const commandDef = commands[cmdName];

    if (commandDef) {
      try {
        const response = await commandDef.execute(args);
        pushToHistory(response.output, response.type);
      } catch (err) {
        pushToHistory(`Error executing '${cmdName}': ${err}`, "error");
      }
    } else {
      pushToHistory(
        `Command not found: ${cmdName}. Type 'help' for list.`,
        "error",
      );
    }
  }, []);

  return {
    history,
    execute,
    clear: () => setHistory([]),
  };
};
