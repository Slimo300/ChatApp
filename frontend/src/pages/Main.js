import React, { useContext, useEffect, useState } from "react";
import {Navigate} from "react-router-dom";
import { actionTypes, StorageContext } from "../ChatStorage";
import Chat from "../components/Chat";
import { GroupLabel } from "../components/GroupLabel";
import { ModalCreateGroup } from "../components/modals/CreateGroup";
import { GetGroups, GetInvites, GetMessages, GetUser, GetWebsocket } from "../Requests";

const Main = (props) => {

    return (
        <div>
            {props.name === ""? <Navigate to="/login" />:<AuthMain {...props}/>}
        </div>
    );
}

const AuthMain = (props) => {

    const [state, dispatch] = useContext(StorageContext);
    const [current, setCurrent] = useState({}); // current group
    const [toggler, setToggler] = useState(false); // toggler for scrollRef
    function toggleToggler(){
        setToggler(!toggler);
    }
    const [ws, setWs] = useState({}); // websocket connection

    // Getting user data, groups and invites and setting websocket connection
    useEffect(() => {
        let userPromise = GetUser();
        userPromise.then( response => { dispatch({type: actionTypes.LOGIN, payload: response}) } );
        let groupsPromise = GetGroups();
        groupsPromise.then( response => { dispatch({type: actionTypes.SET_GROUPS, payload: response}) } );
        let invites = GetInvites();
        invites.then( response => { dispatch({type: actionTypes.SET_NOTIFICATIONS, payload: response}) } );
        let websocketPromise = GetWebsocket();
        websocketPromise.then( response => { setWs(response) } );
    }, [dispatch]);

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
        if (msgJSON.group === current.ID) { // add message to state
            dispatch({type: actionTypes.ADD_MESSAGE, payload: {message: msgJSON, current: true}})
            toggleToggler();
        } else {
            dispatch({type: actionTypes.ADD_MESSAGE, payload: {message: msgJSON, current: false}})
        }
    }

    // getting messages from specific group
    useEffect(()=>{
        (
            async () => {
                if (current.ID !== undefined && current.messages.length === 0) {
                    let messagesPromise = GetMessages(current.ID.toString())
                    messagesPromise.then( response => { dispatch({type: actionTypes.SET_MESSAGES, payload: {messages: response, group: current.ID}}) } )
                    toggleToggler();
                }
            }
        )();
    }, [current, dispatch]);

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
                                            {state.groups.length!==0?state.groups.map(item => {return <GroupLabel key={item.ID} setCurrent={setCurrent} group={item}/>}):null}
                                        </ul>
                                    </div>
                                </div>
                                <Chat group={current} ws={ws} setCurrent={setCurrent} toggler={toggler}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
          <ModalCreateGroup show={props.showCrGroup} toggle={props.toggleCrGroup}/>
        </div>
    )
}
export default Main;