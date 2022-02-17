import React from "react";
import { BrowserRouter as Router, Route, Routes} from "react-router-dom";
import SignInForm from './Login';
import AuthMain from "./Main";
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

  render() {
    return (
        <div>
          <Router>
            <div className="container pt-4 mt-4">
              <Routes>
                <Route path="/" element={<AuthMain/>} exact />
                <Route path="/login" element={<SignInForm/>} exact />
                <Route path="/register" element={<RegisterForm/>} exact />
              </Routes>
            </div>
          </Router>
        </div>
    )
  }
}

export default App;