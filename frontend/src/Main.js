import React, { useEffect, useState } from "react";
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

    const [groups, setGroups] = useState([]);
    useEffect(()=>{
        (
            async () => {
                const response = await fetch('http://localhost:8080/api/group/get', {
                    headers: {'Content-Type': 'application/json'},
                    credentials: 'include'});
                const a = await response.json();
                setGroups(a);
            }
        )();
    }, []);

    return (
        <div className="container-fluid row">
            <div className="list-group h-100 col-4 col-md-4 col-lg-2">
                { groups.map(item => {return (<GroupButton name={item.name} key={item.ID} />)} ) }
            </div>
            <div className="h-100 col-8 col-md-8 col-lg-10">
                <h1>Message Window</h1>
            </div>
        </div>
    )
}

const GroupButton = (props) => {
    return (
        <button type="button" data-internalid={props.id} className="list-group-item list-group-item-action">{props.name}</button>
    )
}



export default Main;