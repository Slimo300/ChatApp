import React, {useState, useEffect} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './Login';
import AuthMain from "./Main";
import Navigation from "./Navigation";
import RegisterForm from "./Register";

function App() {
  
  const [name, setName] = useState('');

  // useEffect(() => {
  //   async function fetchMyAPI() {
  //     let response = await fetch('http://localhost:8080/api/user', {
  //       headers: {'Content-Type': 'application/json'},
  //       credentials: 'include',
  //   })
  //     response = await response.json()
  //     setName(response.name)
  //   }

  //   fetchMyAPI()
  // }, [])

  useEffect(() => {
    (
        async () => {
            const response = await fetch('http://localhost:8080/api/user', {
                headers: {'Content-Type': 'application/json'},
                credentials: 'include',
            });

            const content = await response.json();

            setName(content.name);
        }
    )();
  });

    return (
        <div>
          <Router>
            <Navigation user={name} setName={setName}></Navigation>
            <main className="container pt-4 mt-4">
              <Routes>
                <Route path="/" element={<AuthMain name={name}/>} exact />
                <Route path="/login" element={<SignInForm setName={setName}/>} exact />
                <Route path="/register" element={<RegisterForm/>} exact />
              </Routes>
            </main>
          </Router>
        </div>
    )
}

export default App;