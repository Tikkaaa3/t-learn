import { useState, useRef, useEffect } from "react";
import Markdown from "react-markdown";
import "./index.css";
import { useTerminal } from "./hooks/useTerminal";

function App() {
  const { history, execute, promptLabel, commandHistory } = useTerminal();

  const [input, setInput] = useState("");
  // Pointer tracks our position in history (null = typing new command)
  const [historyPointer, setHistoryPointer] = useState<number | null>(null);
  // Cursor position in the input string
  const [cursorPosition, setCursorPosition] = useState(0);

  const inputRef = useRef<HTMLInputElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [history]);

  const handleFocus = () => inputRef.current?.focus();

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pastedText = e.clipboardData.getData("text");
    const newInput = input.slice(0, cursorPosition) + pastedText + input.slice(cursorPosition);
    setInput(newInput);
    setCursorPosition(cursorPosition + pastedText.length);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      execute(input);
      setInput("");
      setCursorPosition(0);
      setHistoryPointer(null); // Reset pointer on submit
      return;
    }

    if (e.key === "ArrowLeft") {
      e.preventDefault();
      setCursorPosition((prev) => Math.max(0, prev - 1));
      return;
    }

    if (e.key === "ArrowRight") {
      e.preventDefault();
      setCursorPosition((prev) => Math.min(input.length, prev + 1));
      return;
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
      const newInput = commandHistory[newIndex];
      setInput(newInput);
      setCursorPosition(newInput.length); // Move cursor to end
      return;
    }

    if (e.key === "ArrowDown") {
      e.preventDefault();
      if (historyPointer === null) return; // Already at bottom

      if (historyPointer < commandHistory.length - 1) {
        // Go forward 1
        const newIndex = historyPointer + 1;
        setHistoryPointer(newIndex);
        const newInput = commandHistory[newIndex];
        setInput(newInput);
        setCursorPosition(newInput.length); // Move cursor to end
      } else {
        // We reached the end, clear input
        setHistoryPointer(null);
        setInput("");
        setCursorPosition(0);
      }
      return;
    }

    if (e.key === "Backspace") {
      e.preventDefault();
      if (cursorPosition === 0) return; // Nothing to delete
      
      const newInput = input.slice(0, cursorPosition - 1) + input.slice(cursorPosition);
      setInput(newInput);
      setCursorPosition(cursorPosition - 1);
      return;
    }

    if (e.key === "Delete") {
      e.preventDefault();
      if (cursorPosition === input.length) return; // Nothing to delete
      
      const newInput = input.slice(0, cursorPosition) + input.slice(cursorPosition + 1);
      setInput(newInput);
      // Cursor position stays the same
      return;
    }

    if (e.key === "Home") {
      e.preventDefault();
      setCursorPosition(0);
      return;
    }

    if (e.key === "End") {
      e.preventDefault();
      setCursorPosition(input.length);
      return;
    }

    // Handle regular character input
    if (e.key.length === 1 && !e.ctrlKey && !e.altKey && !e.metaKey) {
      e.preventDefault();
      const newInput = input.slice(0, cursorPosition) + e.key + input.slice(cursorPosition);
      setInput(newInput);
      setCursorPosition(cursorPosition + 1);
      return;
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
        <span>{input.slice(0, cursorPosition)}</span>
        <span className="cursor"></span>
        <span>{input.slice(cursorPosition)}</span>
      </div>

      <input
        ref={inputRef}
        className="hidden-input"
        autoFocus
        value={input}
        onKeyDown={handleKeyDown}
        onPaste={handlePaste}
        autoComplete="off"
        // Disable default onChange to prevent conflicts
        readOnly
      />

      <div ref={bottomRef} />
    </div>
  );
}

export default App;
