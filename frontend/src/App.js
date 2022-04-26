import React, {useState} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './pages/Login';
import Main from "./pages/Main";
import Navigation from "./components/Navigation";
import RegisterForm from "./pages/Register";
import ChatStorage from "./ChatStorage";
function App() {
  
  const [name, setName] = useState('');

  const [createGrShow, setCreateGrShow] = useState(false);
  const toggleCreateGroup = () => {
    setCreateGrShow(!createGrShow);
  }

  return (
      <div >
        <ChatStorage>
        <Router>
          <Navigation name={name} setName={setName} toggleCrGroup={toggleCreateGroup}/>
          <main>
            <Routes>
              <Route path="/" element={<Main name={name} toggleCrGroup={toggleCreateGroup} showCrGroup={createGrShow}/>}/>
              <Route path="/login" element={<SignInForm setName={setName} name={name}/>}/>
              <Route path="/register" element={<RegisterForm/>}/>
            </Routes>
          </main>
        </Router>
        </ChatStorage>
      </div>
  )
}

export default App;