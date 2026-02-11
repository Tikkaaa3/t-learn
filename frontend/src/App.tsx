import { useState, useRef, useEffect } from "react";
import Markdown from "react-markdown";
import "./index.css";
import { useTerminal } from "./hooks/useTerminal";

function App() {
  const { history, execute, promptLabel, commandHistory } = useTerminal();

  const [input, setInput] = useState("");
  // Pointer tracks our position in history (null = typing new command)
  const [historyPointer, setHistoryPointer] = useState<number | null>(null);

  const inputRef = useRef<HTMLInputElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [history]);

  const handleFocus = () => inputRef.current?.focus();

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      execute(input);
      setInput("");
      setHistoryPointer(null); // Reset pointer on submit
    }

    if (e.key === "ArrowUp") {
      e.preventDefault();
      if (commandHistory.length === 0) return;

      // Calculate new index: If null, start at end. Else go back 1.
      const newIndex =
        historyPointer === null
          ? commandHistory.length - 1
          : Math.max(0, historyPointer - 1);

      setHistoryPointer(newIndex);
      setInput(commandHistory[newIndex]);
    }

    if (e.key === "ArrowDown") {
      e.preventDefault();
      if (historyPointer === null) return; // Already at bottom

      if (historyPointer < commandHistory.length - 1) {
        // Go forward 1
        const newIndex = historyPointer + 1;
        setHistoryPointer(newIndex);
        setInput(commandHistory[newIndex]);
      } else {
        // We reached the end, clear input
        setHistoryPointer(null);
        setInput("");
      }
    }
  };

  return (
    <div className="terminal-container" onClick={handleFocus}>
      <div className="history">
        {history.map((line) => (
          <div
            key={line.id}
            className={`line ${line.type}`}
            style={{ marginBottom: "8px" }}
          >
            {line.type === "command" ? (
              <>
                <span style={{ color: "#fff", marginRight: "10px" }}>$</span>
                <span style={{ whiteSpace: "pre-wrap" }}>{line.content}</span>
              </>
            ) : (
              <div className="markdown-output">
                <Markdown>{line.content}</Markdown>
              </div>
            )}
          </div>
        ))}
      </div>

      <div
        className="input-line"
        style={{ display: "flex", alignItems: "center" }}
      >
        <span style={{ color: "#fff", marginRight: "10px" }}>
          {promptLabel}
        </span>
        <span>{input}</span>
        <span className="cursor"></span>
      </div>

      <input
        ref={inputRef}
        className="hidden-input"
        autoFocus
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        autoComplete="off"
      />

      <div ref={bottomRef} />
    </div>
  );
}

export default App;
