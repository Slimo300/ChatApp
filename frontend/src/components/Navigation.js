import React, { useContext } from 'react';
import {NavLink} from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBell } from '@fortawesome/free-solid-svg-icons'
import { Logout } from '../Requests';
import { StorageContext, actionTypes } from '../ChatStorage';
import Invite from './Invite';

const Navigation = (props) => {

    const [state, dispatch] = useContext(StorageContext);

    const logout = async () => {
        let logoutPromise = Logout();
        logoutPromise.then( response => {
            if (response === undefined) {
                dispatch({type: actionTypes.LOGOUT})
                props.setName("");
            }
            else alert(response.message);
        })
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
                    <button type='button' className="navbar-brand order-1 btn btn-dark text-primary" onClick={props.toggleProfile}>Profile</button>
                </ul>

                <div className="btn-group">
                    <button type="button" className="btn btn-secondary dropdown-toggle bg-dark" data-toggle="dropdown" aria-expanded="false">
                        <FontAwesomeIcon className='text-primary pr-2' icon={faBell} />
                        <span className="badge badge-secondary">{state.notifications!==undefined?state.notifications.length:null}</span>
                    </button>
                    <div className="dropdown-menu dropdown-menu-right">
                        {state.notifications!==undefined?state.notifications.map((item)=> <Invite invite={item} />):null}
                    </div>
                </div>

                <NavLink className="nav-item nav-link" to="/login" onClick={logout}>Logout</NavLink>
            </div>
        );
    }

    return (
        <nav className="navbar navbar-expand-md navbar-dark bg-dark mb-4">
            <NavLink className="navbar-brand" to="/" >ChatApp</NavLink>
            <div className="collapse navbar-collapse" id="navbarCollapse">
                {menu}
            </div>
        </nav>
    )
}

export default Navigation;