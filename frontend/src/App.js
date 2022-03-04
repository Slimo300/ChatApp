import React, {useState, useEffect} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './Login';
import Main from "./Main";
import Navigation from "./Navigation";
import RegisterForm from "./Register";

function App() {
  
  const [name, setName] = useState('');

  useEffect(() => {
    (
        async () => {
            const response = await fetch('http://localhost:8080/api/user', {
                headers: {'Content-Type': 'application/json'},
                credentials: 'include',
            });

            const content = await response.json();
            console.log(content);
            setName(content.username);
        }
    )();
  });

    return (
        <div>
          <Router>
            <Navigation name={name} setName={setName}></Navigation>
            <main className="container pt-4 mt-4">
              <Routes>
                <Route path="/" element={<Main name={name}/>}/>
                <Route path="/login" element={<SignInForm setName={setName}/>}/>
                <Route path="/register" element={<RegisterForm/>}/>
              </Routes>
            </main>
          </Router>
        </div>
    )
}

export default App;