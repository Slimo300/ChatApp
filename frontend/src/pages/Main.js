import React, { useContext, useEffect, useState } from "react";
import {Navigate} from "react-router-dom";
import { actionTypes, StorageContext } from "../ChatStorage";
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

    const [state, dispatch] = useContext(StorageContext);
    const [counter, setCounter] = useState({}); // object mapping group_id to unread messages
    const [messages, setMessages] = useState([]); // list of messages for current group
    const [current, setCurrent] = useState({}); // current group
    const [ws, setWs] = useState({}); // websocket connection

    // Effect getting user info
    useEffect(() => {
        (
            async () => {
                const response = await fetch('http://localhost:8080/api/user', {
                    method: 'GET',
                    headers: {'Content-Type': 'application/json'},
                    credentials: 'include'});
                if (response.status !== 200) {
                    throw new Error("couldn't get user");
                }
                const responseJSON = await response.json()
                dispatch({type: actionTypes.LOGIN, payload: responseJSON})
            }
        )();
        
    }, [dispatch]);

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
        if (msgJSON.action !== undefined) {
            switch (msgJSON.action) {
                case "DELETE_GROUP":
                    dispatch({type: actionTypes.DELETE_GROUP, payload: msgJSON.group.ID});
                    break;
                case "DELETE_MEMBER":
                    dispatch({type: actionTypes.DELETE_MEMBER, payload: msgJSON.member});
                    break;
                case "ADD_MEMBER":
                    dispatch({type: actionTypes.ADD_MEMBER, payload: msgJSON.member});
                    break;
                default:
                    console.log("Unexpected action from websocket: ", msgJSON.action);
            }
            return;
        }
        if (msgJSON.group === current.ID) {
            setMessages([...messages, msgJSON]);
        } else {
            let newCounter = {
                ...counter
              };
            newCounter[msgJSON.group]++;
            setCounter(newCounter);
        }
    }

    // effect getting groups of which user is a member
    useEffect(()=>{
        (
            async () => {
                const response = await fetch('http://localhost:8080/api/group/get', {
                    headers: {'Content-Type': 'application/json'},
                    credentials: 'include'});
                if (response.status !== 200 && response.status !== 204 ) {
                    throw new Error("Invalid response when requesting user groups");
                }
                const responseJSON = await response.json();
                if (responseJSON.message === undefined) {
                    let newCounter = {};
                    for (let i = 0; i < responseJSON.length; i++) {
                        newCounter[responseJSON[i].ID] = 0;
                    }
                    setCounter(newCounter);
                    dispatch({type: actionTypes.SET_GROUPS, payload: responseJSON});
                }
            }
        )();
    }, []);

    // getting messages from specific group
    useEffect(()=>{
        (
            async () => {
                if (current.ID !== undefined) {
                    const response = await fetch("http://localhost:8080/api/group/messages?group=" + current.ID.toString() + "&num=4", {
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
                                            {state.groups.length!==0?state.groups.map(item => {return <GroupLabel counter={counter} setCounter={setCounter} key={item.ID} setCurrent={setCurrent} group={item}/>}):null}
                                        </ul>
                                    </div>
                                </div>
                                <Chat messages={messages} group={current} ws={ws} setCurrent={setCurrent}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
          <ModalCreateGroup show={props.showCrGroup} toggle={props.toggleCrGroup}/>
          <ModalAddUser show={props.showFrAdd} toggle={props.toggleFrAdd}/>
        </div>
    )
}
export default Main;