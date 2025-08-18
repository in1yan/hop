import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {GetWindows} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState([]);
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    function getWin() {
        GetWindows().then(updateResultText);
        console.log(resultText);
    }

    return (
        <div id="App">
            <div id="result" className="result">
                <ul>
                {resultText.map((window, index) =>(
                        <li key={window.Handle} className="result-item">{window.Title}</li>
                ))}
                </ul>
            </div>
                {/* <button className="btn" onClick={getWin}>Get Windows</button> */}
        </div>
    )
}

export default App
