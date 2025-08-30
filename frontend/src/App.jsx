import { useEffect, useState, useRef, useMemo } from "react";
import "./App.css";
import { GetWindows, SetFocus } from "../wailsjs/go/main/App";
import { EventsOn, WindowGetSize, WindowSetSize, WindowShow , Quit} from "../wailsjs/runtime/runtime";

function App() {
  const [windows, setWindows] = useState([]);
  const [allWindows, setAllWindows] = useState([]);
  const [filterText, setFilterText] = useState("");
  const hiddenInputRef = useRef(null);

  const updateWindows = async (result) => {
    setAllWindows(result);
    setWindows(result);
    setFilterText(""); // Reset filter when windows update
    const baseHeight = 10;
    const itemHeight = 30;
    const newHeight = baseHeight + itemHeight * result.length;
    const size = await WindowGetSize();
    WindowSetSize(size.w, newHeight);
  };

  function getWin() {
    GetWindows().then(updateWindows);
  }

  function generateHints(n) {
    const chars = "abcdefghijklmnopqrstuvwxyz";
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

  const hints = useMemo(() => generateHints(windows.length), [windows.length]);

  // Filter windows based on hint starting character
  const filteredWindows = useMemo(() => {
    if (!filterText) return allWindows;
    
    const filteredHints = generateHints(allWindows.length);
    return allWindows.filter((_, index) => 
      filteredHints[index] && filteredHints[index].toLowerCase().startsWith(filterText.toLowerCase())
    );
  }, [allWindows, filterText]);

  // Update windows when filtered
  useEffect(() => {
    setWindows(filteredWindows);
    const baseHeight = 10;
    const itemHeight = 30;
    const newHeight = baseHeight + itemHeight * filteredWindows.length;
    WindowGetSize().then(size => {
      WindowSetSize(size.w, newHeight);
    });
  }, [filteredWindows]);

  useEffect(() => {
    getWin();
    EventsOn("windows:update", (windows) => {
      updateWindows(windows);
      WindowShow();
    });

    // Listen for focus event from backend
    EventsOn("focus:input", () => {
      if (hiddenInputRef.current) {
        hiddenInputRef.current.focus();
      }
    });
  }, []);

  useEffect(() => {
    const handleHiddenInputChange = (e) => {
      const value = e.target.value;
      setFilterText(value);
      
      // Check if we have an exact match for a hint
      const currentHints = generateHints(allWindows.length);
      const matchedIndex = currentHints.findIndex(hint => 
        hint.toLowerCase() === value.toLowerCase()
      );
      
      if (matchedIndex !== -1 && allWindows[matchedIndex]) {
        SetFocus(allWindows[matchedIndex].Handle);
        setFilterText(""); // Reset after selection
        if (hiddenInputRef.current) {
          hiddenInputRef.current.value = "";
        }
      }
    };

    const hiddenInput = hiddenInputRef.current;
    if (hiddenInput) {
      hiddenInput.addEventListener('input', handleHiddenInputChange);
      return () => {
        hiddenInput.removeEventListener('input', handleHiddenInputChange);
      };
    }
  }, [allWindows]);

  // Keep the old keyboard handler for direct key presses (fallback)
  useEffect(() => {
    const handler = async (e) => {
      if(e.key === "Escape"){
        Quit();
      }
      const key = e.key.toUpperCase();

      if (!"HIJKL".includes(key)) return;
    };

    document.addEventListener("keydown", handler);
    return () => {
      document.removeEventListener("keydown", handler);
    };
  }, [allWindows]);

  return (
    <div id="app">
      {/* Hidden input for capturing typed hints */}
      <input
        ref={hiddenInputRef}
        type="text"
        style={{
          position: 'absolute',
          left: '-9999px',
          top: '-9999px',
          opacity: 0,
          width: '1px',
          height: '1px',
        }}
        autoComplete="off"
        tabIndex={-1}
      />
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
