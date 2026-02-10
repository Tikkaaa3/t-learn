import { useState, useRef, useEffect } from "react";
import "./index.css";
import { useTerminal } from "./hooks/useTerminal";

function App() {
  // 1. Destructure promptLabel from the hook
  const { history, execute, promptLabel } = useTerminal();
  const [input, setInput] = useState("");

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
            {line.type === "command" && (
              // Optional: You can make this dynamic too, but keeping it '$'
              // for history is often cleaner. Let's stick to '$' for history for now.
              <span style={{ color: "#fff", marginRight: "10px" }}>$</span>
            )}
            <span style={{ whiteSpace: "pre-wrap" }}>{line.content}</span>
          </div>
        ))}
      </div>

      {/* The Active Input Line */}
      <div
        className="input-line"
        style={{ display: "flex", alignItems: "center" }}
      >
        {/* 2. USE THE DYNAMIC PROMPT HERE */}
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
      />

      <div ref={bottomRef} />
    </div>
  );
}

export default App;
