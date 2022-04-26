import React from 'react';
import {NavLink} from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBell } from '@fortawesome/free-solid-svg-icons'

const Navigation = (props) => {
    const logout = async () => {
        await fetch("http://localhost:8080/api/signout", {
            method: "POST",
            credentials: "include",
            headers: {"Content-Type": "application/json"}
        });
        props.setName("");
    };

    let menu;

    if (props.name === "") {
        menu = (
            <div className="collapse navbar-collapse" id="navbarCollapse">
                <ul className="navbar-nav mr-auto"></ul>
                <NavLink className="nav-item nav-link" to="/login">Login</NavLink>
                <NavLink className="nav-item nav-link" to="/register">Register</NavLink>
            </div>
        );
    } else {
        menu = (
            <div className="collapse navbar-collapse" id="navbarCollapse">
                <ul className="navbar-nav mr-auto">
                    <button type='button' className="navbar-brand order-1 btn btn-dark text-primary" onClick={props.toggleCrGroup}>Create Room</button>
                </ul>

                <button className='btn btn-dark'><FontAwesomeIcon className='text-primary' icon={faBell} /></button>
                <NavLink className="nav-item nav-link" to="/login" onClick={logout}>Logout</NavLink>
            </div>
        );
    }

    return (
        <nav className="navbar navbar-expand-md navbar-dark bg-dark mb-4">
            <div className="navbar-brand" >ChatApp</div>
            <div className="collapse navbar-collapse" id="navbarCollapse">
                {menu}
            </div>
        </nav>
    )
}

export default Navigation;