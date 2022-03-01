import React, {useState} from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './Login';
import AuthMain from "./Main";
import Navigation from "./Navigation";
import RegisterForm from "./Register";

class App extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      user: ""
    }
    this.setName = this.setName.bind(this)
  }
  
  setName(name) {
    this.setState(
    {
      user: name
    });
  }

  render() {
    return (
        <div>
          <Router>
            <Navigation user={this.state.user}></Navigation>
            <div className="container pt-4 mt-4">
              <Routes>
                <Route path="/" element={<AuthMain name={this.state.user}/>} exact />
                <Route path="/login" element={<SignInForm setName={this.setName}/>} exact />
                <Route path="/register" element={<RegisterForm/>} exact />
              </Routes>
            </div>
          </Router>
        </div>
    )
  }
}

export default App;