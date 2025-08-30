import { useEffect, useState, useRef, useMemo } from "react";
import "./App.css";
import { GetWindows, SetFocus } from "../wailsjs/go/main/App";
import { EventsOn, WindowGetSize, WindowSetSize } from "../wailsjs/runtime/runtime";

function App() {
  const [windows, setWindows] = useState([]);
  const inputBuffer = useRef(""); // 👈 persistent buffer

  const updateWindows = async (result) => {
    setWindows(result);
    const baseHeight = 10;
    const itemHeight = 25;
    const newHeight = baseHeight + itemHeight * result.length;
    const size = await WindowGetSize();
    WindowSetSize(size.w, newHeight);
  };

  function getWin() {
    GetWindows().then(updateWindows);
  }

  function generateHints(n) {
    const chars = "HIJKL";
    let hints = [];
    let queue = [""];
    while (hints.length < n) {
      let pre = queue.shift();
      for (let ch of chars) {
        let hint = pre + ch;
        hints.push(hint);
        if (hints.length === n) break;
        queue.push(hint);
      }
    }
    return hints;
  }

  // memoize hints so it doesn’t regenerate on every render unnecessarily
  const hints = useMemo(() => generateHints(windows.length), [windows.length]);

  useEffect(() => {
    getWin();
    EventsOn("windows:update", (windows) => {
      updateWindows(windows);
    });
  }, []);

  useEffect(() => {
    const handler = async (e) => {
      const key = e.key.toUpperCase();

      if (!"HIJKL".includes(key)) return;

      inputBuffer.current += key;

      const matched = hints.findIndex((h) => h === inputBuffer.current);
      if (matched !== -1) {
        let winid = windows[matched].Handle;
        await SetFocus(winid);
        inputBuffer.current = ""; // reset buffer
      }
    };

    document.addEventListener("keydown", handler);

    return () => {
      document.removeEventListener("keydown", handler);
    };
  }, [hints, windows]);

  return (
    <div id="app">
      <ul className="list">
        {windows.map((window, index) => (
          <li
            key={window.Handle}
            className="result-item"
            onClick={() => SetFocus(window.Handle)}
          >
            {window.Title + "  " + hints[index]}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
