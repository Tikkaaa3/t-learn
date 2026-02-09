import { useState, useRef, useEffect } from "react";
import "./index.css";
import { useTerminal } from "./hooks/useTerminal";

function App() {
  const { history, execute } = useTerminal();
  const [input, setInput] = useState("");

  const inputRef = useRef<HTMLInputElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Auto-scroll
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [history]);

  // Keep focus
  const handleFocus = () => inputRef.current?.focus();

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      execute(input);
      setInput("");
    }
  };

  return (
    <div
      className="terminal-container"
      onClick={handleFocus}
      style={{ padding: "20px", height: "100%", overflowY: "auto" }}
    >
      <div className="history">
        {history.map((line) => (
          <div
            key={line.id}
            className={`line ${line.type}`}
            style={{ marginBottom: "8px" }}
          >
            {line.type === "command" && (
              <span style={{ color: "#fff", marginRight: "10px" }}>$</span>
            )}
            <span style={{ whiteSpace: "pre-wrap" }}>{line.content}</span>
          </div>
        ))}
      </div>

      <div
        className="input-line"
        style={{ display: "flex", alignItems: "center" }}
      >
        <span style={{ color: "#fff", marginRight: "10px" }}>$</span>
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
