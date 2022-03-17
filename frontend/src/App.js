import React, {useState} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './Login';
import Main from "./Main";
import Navigation from "./Navigation";
import RegisterForm from "./Register";

function App() {
  
  const [name, setName] = useState('');
  

  const ws = new WebSocket("ws://localhost:8080/ws")

  ws.onopen = () => {
      console.log("Websocket openned");
  };

  ws.onclose = () => {
      console.log("closed");
  }


  return (
      <div >
        <Router>
          <Navigation name={name} setName={setName}></Navigation>
          <main>
            <Routes>
              <Route path="/" element={<Main name={name} ws={ws}/>}/>
              <Route path="/login" element={<SignInForm setName={setName} name={name}/>}/>
              <Route path="/register" element={<RegisterForm/>}/>
            </Routes>
          </main>
        </Router>
      </div>
  )
}

export default App;