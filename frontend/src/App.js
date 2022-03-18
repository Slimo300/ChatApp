import React, {useState} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './pages/Login';
import Main from "./pages/Main";
import Navigation from "./components/Navigation";
import RegisterForm from "./pages/Register";
import {ModalCreateGroup, ModalAddFriend} from "./components/Modals";

function App() {
  
  const [name, setName] = useState('');
  const [createGrShow, setCreateGrShow] = useState(false);
  const [addFrShow, setAddFrShow] = useState(false);

  const toggleCreateGroup = () => {
    setCreateGrShow(!createGrShow);
  }
  
  const toggleAddFriend = () => {
    setAddFrShow(!addFrShow);
  }

  return (
      <div >
        <Router>
          <Navigation name={name} setName={setName} toggleCrGroup={toggleCreateGroup} toggleFrAdd={toggleAddFriend}/>
          <main>
            <Routes>
              <Route path="/" element={<Main name={name} toggleCrGroup={toggleCreateGroup}/>}/>
              <Route path="/login" element={<SignInForm setName={setName} name={name}/>}/>
              <Route path="/register" element={<RegisterForm/>}/>
            </Routes>
          </main>
          <ModalCreateGroup show={createGrShow} toggle={toggleCreateGroup}/>
          <ModalAddFriend show={addFrShow} toggle={toggleAddFriend}/>
        </Router>
      </div>
  )
}

export default App;