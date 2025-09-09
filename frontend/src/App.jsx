import { useEffect, useState, useRef, useMemo } from "react";
import "./App.css";
import { GetWindows, SetFocus } from "../wailsjs/go/main/App";
import { EventsOn, WindowGetSize, WindowSetSize, WindowShow , Quit, WindowSetBackgroundColour, WindowCenter} from "../wailsjs/runtime/runtime";

function App() {
  const [windows, setWindows] = useState([]);
  const [allWindows, setAllWindows] = useState([]);
  const [filterText, setFilterText] = useState("");
  const hiddenInputRef = useRef(null);

  const updateWindows = async (result) => {
    setAllWindows(result);
    setWindows(result);
    setFilterText(""); // Reset filter when windows update
  };
async function updateWindowSize() {
    const size = await WindowGetSize();

    // Measure the .list container
    const listEl = document.querySelector(".list");
    let listHeight = 0;
    if (listEl) {
        listHeight = listEl.scrollHeight;
    }

    // Add app padding (top + bottom) from CSS (#app { padding: 8px })
    const appPadding = 16; // 8px top + 8px bottom

    // Minimum height fallback
    const minHeight = 50;

    const calculatedHeight = Math.max(minHeight, listHeight + appPadding);

    await WindowSetSize(size.w, calculatedHeight);
}


  // Update window size whenever windows change
  useEffect(() => {
    if (windows.length > 0) {
      console.log(`Windows changed: ${windows.length} windows detected`);
      updateWindowSize();
    }
  }, [windows]);

  function getWin() {
    GetWindows().then((result) => {
      updateWindows(result);
      updateWindowSize();
    });
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
  }, [filteredWindows]);

  useEffect(() => {
    getWin();
    EventsOn("windows:update", (windows) => {
      console.log(`EventsOn windows:update - ${windows.length} windows`);
      updateWindows(windows);
      updateWindowSize();
      WindowSetBackgroundColour(25, 23, 36, 180); // Ensure transparency on update
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

  // Handle escape key to quit application
  useEffect(() => {
    const handler = async (e) => {
      if(e.key === "Escape"){
        Quit();
      }
    };

    document.addEventListener("keydown", handler);
    return () => {
      document.removeEventListener("keydown", handler);
    };
  }, []);

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
          <li key={window.Handle}>
            <div
              className="result-item"
              onClick={() => SetFocus(window.Handle)}
            >
              {window.Title} <span className="split">|</span> <span className="hint">{hints[index].toUpperCase()}</span>
            </div>
            {index < windows.length - 1 && <hr className="sep" />}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
