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
    const [current, setCurrent] = useState(0);
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
        <div className="container" >
            <div className="content-wrapper">
                <div className="row gutters">
                    <div className="col-xl-12 col-lg-12 col-md-12 col-sm-12 col-12">
                        <div className="card m-0">
                            <div className="row no-gutters">
                                <div className="col-xl-4 col-lg-4 col-md-4 col-sm-3 col-3">
                                    <div className="users-container">
                                        <ul className="users">
                                            {groups.map(item => {return <GroupLabel name={item.name} key={item.ID}/>})}
                                        </ul>
                                    </div>
                                </div>
                                <Chat />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

const Chat = (props) => {
    let load;

    if (props.group === undefined) {
        load = <h1>Select a group to chat!</h1>;
    } else {
    load = (
        <div class="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
            <div class="selected-user">
                <span>To: <span class="name">{props.name}</span></span>
            </div>
            <div class="chat-container">
                <ul class="chat-box chatContainerScroll">
                    
                </ul>
                <div class="form-group mt-3 mb-0">
                    <textarea class="form-control" rows="3" placeholder="Type your message here..."></textarea>
                </div>
            </div>
        </div>
    );
    }
    return load;
}

const GroupLabel = (props) => {
    return (
        <li className="person" >
            <div className="user">
                <img src="https://www.bootdey.com/img/Content/avatar/avatar3.png" alt="Retail Admin"/>
                <span className={"status busy" + props.status}></span>
            </div>
            <p className="name-time">
                <span className="name">{props.name}</span>
            </p>
        </li>
    )
}

const Message = (props) => {
    const right = (
        <li className="chat-right">
            <div className="chat-hour">{props.time} <span className="fa fa-check-circle"></span></div>
            <div className="chat-text">{props.message}</div>
            <div className="chat-avatar">
                <img src="https://www.bootdey.com/img/Content/avatar/avatar3.png" alt="Retail Admin"/>
                <div className="chat-name">{props.name}</div>
            </div>
        </li>
    );

    const left = (
        <li className="chat-left">
            <div className="chat-avatar">
                <img src="https://www.bootdey.com/img/Content/avatar/avatar3.png" alt="Retail Admin"/>
                <div className="chat-name">{props.name}</div>
            </div>
            <div className="chat-text">{props.message}</div>
            <div className="chat-hour">{props.time} <span class="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {props.owned?left:right}
        </div>
    )
}

export default Main;