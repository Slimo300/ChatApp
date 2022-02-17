import React from "react";
import { unstable_renderSubtreeIntoContainer } from "react-dom";
import { BrowserRouter as Router, Route } from "react-router-dom";
import SignInForm from './Login';
import RegisterForm from "./Register";

class App extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      user: {
        loggedin: false,
        name: "",
      }
    };
  }

  renderLogin() {
    return (
      <SignInForm />
    )
  }

  renderRegister() {
    return (
      <RegisterForm />
    )
  }
  
  renderSite() {
    return (
      <div>
        <div className="mt-5 row d-flex flex-row">
          <h1>LoggedIn</h1>
        </div>
      </div>
    )
  }

  render() {
    return (
      <main className="d-flex flex-column">
        <div className="container h-100 d-flex flex-column justify-content-center align-items-center">
        {
          !this.state.user.loggedin ?
          this.renderRegister():
          this.renderSite()
        }
        </div>
      </main>
    )
  }
}

export default App;