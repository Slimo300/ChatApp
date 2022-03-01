import React from 'react';
import {NavLink} from 'react-router-dom'

export default class Navigation extends React.Component {
    constructor(props){
        super(props);
        this.state = {
          user: this.props.user
        };
      }
    

    render() {
        return (
            <nav className="navbar navbar-expand-md navbar-dark bg-dark mb-4">
                <a className="navbar-brand" href="#">ChatApp</a>
                <div className="collapse navbar-collapse" id="navbarCollapse">
                    <ul className="navbar-nav mr-auto">
                    
                    </ul>
                    <NavLink className="nav-item nav-link" to="/login">Login</NavLink>
                    <NavLink className="nav-item nav-link" to="/register">Register</NavLink>
                </div>
            </nav>
        )
    }
}

class NavUnathenticated extends React.Component {
    render() {
        return (
            <div className="collapse navbar-collapse" id="navbarCollapse">
                <ul className="navbar-nav mr-auto"></ul>
                <NavLink className="nav-item nav-link" to="/login">Login</NavLink>
                <NavLink className="nav-item nav-link" to="/register">Register</NavLink>
            </div>
        )
    }
}

class NavAuthenticated extends React.Component {
    render() {
        return (
            <div className="collapse navbar-collapse" id="navbarCollapse">
                <ul className="navbar-nav mr-auto">
                    <NavLink className="nav-item nav-link" to="/logout">Create Room</NavLink>
                    <NavLink className="nav-item nav-link" to="/logout">Add Friend</NavLink>
                </ul>

                <NavLink className="nav-item nav-link" to="/logout">Account</NavLink>
                <NavLink className="nav-item nav-link" to="/logout">Logout</NavLink>
            </div>
        )
    }
}