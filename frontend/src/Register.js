import React, { useState } from "react";
import { Navigate } from "react-router-dom";

function RegisterForm() {
  const [username, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rpassword, setRPassword] = useState("");
  const [message, setMessage] = useState("");
  const [redirect, setRedirect] = useState(false);

  const submit = async (e) => {
    e.preventDefault();

    if (password !== rpassword) {
      setMessage("Passwords don't match")
    }

    await fetch('http://localhost:8080/api/register', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({
          username, 
          email,
          password
        })
    });

    setRedirect(true);
  };

  if (redirect) {
    return <Navigate to="/login"/>;
  }

  return (
      <div className="mt-5 d-flex justify-content-center">
        <div className="mt-5 row">
          <form onSubmit={submit}>
            <div className="display-1 mb-4 text-center text-primary"> Register</div>
            <div className="mb-3 text-center">{message}</div>
            <div className="mb-3 text-center">
              <label htmlFor="username" className="form-label">Username</label>
              <input name="username" type="text" className="form-control" id="username" onChange={setName}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="email" className="form-label">Email address</label>
              <input name="email" type="email" className="form-control" id="email" onChange={setEmail}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="pass" className="form-label">Password</label>
              <input name="password" type="password" className="form-control" id="password" onChange={setPassword}/>
            </div>
            <div className="mb-3 text-center">
              <label htmlFor="pass" className="form-label">Repeat Password</label>
              <input name="rpassword" type="password" className="form-control" id="rpassword" onChange={setRPassword}/>
            </div>
            <div className="text-center">
              <button type="submit" className="btn btn-primary text-center">Submit</button>
            </div>
            <div className="display-5 mt-4 text-center text-primary"><a href="#">or Log in</a></div>
          </form>
        </div>
      </div>
    );
  }

export default RegisterForm;