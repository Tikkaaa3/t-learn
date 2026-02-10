import { useState, useRef, useEffect } from "react";
import Markdown from "react-markdown"; // <--- 1. Import Markdown
import "./index.css";
import { useTerminal } from "./hooks/useTerminal";

function App() {
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
            {/* CONDITIONAL RENDERING */}
            {line.type === "command" ? (
              // 1. User Command: Render as plain text with '$'
              <>
                <span style={{ color: "#fff", marginRight: "10px" }}>$</span>
                <span style={{ whiteSpace: "pre-wrap" }}>{line.content}</span>
              </>
            ) : (
              // 2. System Output: Render as Markdown
              <div className="markdown-output">
                <Markdown>{line.content}</Markdown>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Input Line */}
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
