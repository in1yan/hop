import {useState} from 'react';
import './App.css';
import {GetWindows, SetFocus} from "../wailsjs/go/main/App";
import {Hide} from "../wailsjs/runtime/runtime";
function App() {
    const [resultText, setResultText] = useState([]);
    const updateResultText = (result) => setResultText(result);

    function getWin() {
        GetWindows().then(updateResultText);
        console.log(resultText);
    }

    return (
        <div id="App">
            <div id="result" className="result">
                <ul>
                    {resultText.map((window, index) => (
                        <div onClick={() => SetFocus(window.Handle)}>
                            <li key={window.Handle} className="result-item">
                                {window.Title}
                            </li>
                        </div>
                    ))}
                </ul>
                <button className="btn" onClick={getWin}>
                    Get Windows
                </button>
                <button className="btn" onClick={Hide}>
                Hide
                </button>
            </div>
        </div>
    );
}

export default App
