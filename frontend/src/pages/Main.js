import React, { useEffect, useState } from "react";
import {Navigate} from "react-router-dom";
import Chat from "../components/Chat";
import {ws} from "../services/ws";
import { ModalAddFriend, ModalCreateGroup } from "../components/GroupModals";

const Main = (props) => {

    return (
        <div>
            {props.name === ""? <Navigate to="/login" />:<AuthMain {...props}/>}
        </div>
    );
}

const AuthMain = (props) => {

    const [groups, setGroups] = useState([]);
    const [current, setCurrent] = useState(0);
    const [messages, setMessages] = useState([]);
    const [groupname, setGroupName] = useState("");

    ws.onmessage = (e) => {
        console.log("got message");
        const msgJSON = JSON.parse(e.data);
        console.log(msgJSON);
        setMessages([...messages, msgJSON]);
    }

    useEffect(()=>{
        (
            async () => {
                const response = await fetch('http://localhost:8080/api/group/get', {
                    headers: {'Content-Type': 'application/json'},
                    credentials: 'include'});
                const responseJSON = await response.json();
                console.log(responseJSON.message);
                if (responseJSON.message === undefined) {
                    setGroups(responseJSON);
                }
            }
        )();
    }, []);

    useEffect(()=>{
        (
            async () => {
                if (current !== 0) {
                    const response = await fetch("http://localhost:8080/api/group/messages?group=" + current.toString(), {
                        headers: {"Content-Type": "application/json"},
                        credentials: "include",
                    });
                    const responseJSON = await response.json();
                    if (responseJSON.message === "no messages") {
                        setMessages([]);
                    } else {
                        setMessages(responseJSON);
                    }
                }
            }
        )();
    }, [current])

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
                                            {groups.length!==0?groups.map(item => {return <GroupLabel name={item.name} key={item.ID} id={item.ID} setCurrent={setCurrent} setGroupName={setGroupName}/>}):null}
                                        </ul>
                                    </div>
                                </div>
                                <Chat messages={messages} group={current} groupname={groupname} setGroups={setGroups} groups={groups}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
          <ModalCreateGroup show={props.showCrGroup} toggle={props.toggleCrGroup} groups={groups} setGroups={setGroups}/>
          <ModalAddFriend show={props.showFrAdd} toggle={props.toggleFrAdd}/>
        </div>
    )
}

const GroupLabel = (props) => {
    const change = () => {
        props.setCurrent(props.id);
        props.setGroupName(props.name);
    };
    return (
        <li className="person" onClick={change}>
            <div className="user">
                <img src="https://www.bootdey.com/img/Content/avatar/avatar3.png" alt="Retail Admin"/>
                <span className={"status" + props.status}></span>
            </div>
            <p className="name-time">
                <span className="name">{props.name}</span>
            </p>
        </li>
    );
}

export default Main;