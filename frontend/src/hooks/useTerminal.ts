import { useState, useCallback } from "react";
import type { HistoryLine, LineType } from "../types";
import { commands, getPrompt } from "../commands/registry"; // <--- 1. Import getPrompt

export const useTerminal = () => {
  const [history, setHistory] = useState<HistoryLine[]>([
    {
      id: "init",
      type: "info",
      content: "Welcome to t-learn v1.0.0. Type 'help' to start.",
    },
  ]);

  // 2. Add state for the prompt label (Initialize with current global state)
  const [promptLabel, setPromptLabel] = useState(getPrompt());

  const pushToHistory = (content: string, type: LineType = "info") => {
    setHistory((prev) => [
      ...prev,
      { id: Date.now().toString() + Math.random(), type, content },
    ]);
  };

  const execute = useCallback(async (commandString: string) => {
    if (!commandString.trim()) return;

    // Echo the command (We use the *current* promptLabel here effectively in the UI)
    pushToHistory(commandString, "command");

    // NEW: Regex to match spaces BUT ignore spaces inside quotes
    const parts = commandString.match(/(?:[^\s"]+|"[^"]*")+/g);

    if (!parts) return;
    const cmdName = parts[0].toLowerCase();
    // Remove quotes from args (e.g., "Hello World" -> Hello World)
    const args = parts.slice(1).map((arg) => {
      if (arg.startsWith('"') && arg.endsWith('"')) {
        return arg.slice(1, -1);
      }
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

    // 3. FORCE UPDATE PROMPT
    // After the command finishes (e.g., 'lessons' updates the path),
    // we fetch the new string from the registry.
    setPromptLabel(getPrompt());
  }, []);

  return {
    history,
    execute,
    promptLabel, // <--- 4. Export it so App.tsx can use it
    clear: () => setHistory([]),
  };
};
