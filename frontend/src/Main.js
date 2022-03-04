import React from "react";
import {Navigate} from "react-router-dom";

const Main = (props) => {

    console.log(props.name)
    return (
        <div>
            {props.name === undefined? <Navigate to="/login" />:<AuthMain name={props.name}/>}
        </div>
    );
}

const AuthMain = (props) => {

    return (
        <h1>Hello {props.name}!</h1>
    )
}

export default Main;