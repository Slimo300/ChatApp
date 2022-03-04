import React, { useState } from "react";
import { Navigate } from "react-router-dom";

const SignInForm = (props) => {

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [redirect, setRedirect] = useState(false);
  const [message, setMessage] = useState('');

  const submit = async (e) => {
    e.preventDefault();
    if (email === "") {
      setMessage("Email can't be blank")
    }

    const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify({
            email,
            password,
        })
    });

    const content = await response.json();
    console.log(content)

    setRedirect(true);
    props.setName(content.name);
  }

  if (redirect) {
    return <Navigate to="/"/>;
  }

  return (
    <div className="mt-5 d-flex justify-content-center">
      <div className="mt-5 row">
        {message}
        <form onSubmit={submit}>
          <div className="display-3 mb-4 text-center text-primary"> Log In</div>
          <div className="mb-3 text-center">
            <label htmlFor="email" className="form-label">Email address</label>
            <input name="email" type="email" className="form-control" id="email" onChange={setEmail}/>
          </div>
          <div className="mb-3 text-center">
            <label htmlFor="password" className="form-label">Password</label>
            <input name="password" type="password" className="form-control" id="password" onChange={setPassword}/>
          </div>
          <div className="text-center">
            <button type="submit" className="btn btn-primary text-center">Submit</button>
          </div>
          <div className="display-5 mt-4 text-center text-primary"><a href="/register">or Register</a></div>
        </form>
      </div>
    </div>
   )
}

export default SignInForm;