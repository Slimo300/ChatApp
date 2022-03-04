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
        <div className="container-fluid row">
            <div className="list-group h-100 col-4 col-md-4 col-lg-2">
                <button type="button" className="list-group-item list-group-item-action active" aria-current="true">
                    {props.name}
                </button>
                <button type="button" className="list-group-item list-group-item-action">A second item</button>
                <button type="button" className="list-group-item list-group-item-action">A third button item</button>
                <button type="button" className="list-group-item list-group-item-action">A fourth button item</button>
                <button type="button" className="list-group-item list-group-item-action">A disabled button item</button>
            </div>
            <div className="h-100 col-8 col-md-8 col-lg-10">
                <h1>Message Window</h1>
            </div>
        </div>
    )
}



export default Main;