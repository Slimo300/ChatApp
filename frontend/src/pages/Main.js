import React, { useEffect, useState } from "react";
import {Navigate} from "react-router-dom";
import Chat from "../components/Chat";
import { GroupLabel } from "../components/GroupLabel";
import { ModalAddUser } from "../components/modals/AddUser";
import { ModalCreateGroup } from "../components/modals/CreateGroup";

const Main = (props) => {

    return (
        <div>
            {props.name === ""? <Navigate to="/login" />:<AuthMain {...props}/>}
        </div>
    );
}

const AuthMain = (props) => {

    const [groups, setGroups] = useState([]);
    const [messages, setMessages] = useState([]);
    const [current, setCurrent] = useState({});
    const [ws, setWs] = useState({});

    // Effect starting websocket connection
    useEffect(() => {
        let socket = new WebSocket("ws://localhost:8080/ws")
        socket.onopen = () => {
            console.log("Websocket openned");
        };
        socket.onclose = () => {
            console.log("closed");
        };
        setWs(socket);
    }, []);

    ws.onmessage = (e) => {
        const msgJSON = JSON.parse(e.data);
        setMessages([...messages, msgJSON]);
    }

    // effect getting groups in which user has membership
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

    // getting messages from specific group
    useEffect(()=>{
        (
            async () => {
                if (current.ID !== undefined) {
                    const response = await fetch("http://localhost:8080/api/group/messages?group=" + current.ID.toString(), {
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
                                        <ul className="users" style={{height: '85vh', overflow: 'scroll'}}>
                                            {groups.length!==0?groups.map(item => {return <GroupLabel key={item.ID} setCurrent={setCurrent} group={item}/>}):null}
                                        </ul>
                                    </div>
                                </div>
                                <Chat messages={messages} group={current} setGroups={setGroups} groups={groups} ws={ws} setCurrent={setCurrent}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
          <ModalCreateGroup show={props.showCrGroup} toggle={props.toggleCrGroup} groups={groups} setGroups={setGroups}/>
          <ModalAddUser show={props.showFrAdd} toggle={props.toggleFrAdd}/>
        </div>
    )
}
export default Main;