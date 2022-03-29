import React, { useEffect, useState } from "react";
import {v4 as uuidv4} from "uuid";
import GroupMenu from "./GroupMenu";
import { ModalAddUser } from "./modals/AddUser";
import { ModalDeleteGroup } from "./modals/DeleteGroup";
import { ModalMembers } from "./modals/GroupMembers";

const Chat = (props) => {

    const [member, setMember] = useState({});
    const [msg, setMsg] = useState("");

    // add user to group modal
    const [addUserShow, setAddUserShow] = useState(false);
    const toggleAddUser = () => {
        setAddUserShow(!addUserShow);
    };
    // delete group modal
    const [delGrShow, setDelGrShow] = useState(false);
    const toggleDelGroup = () => {
        setDelGrShow(!delGrShow);
    };
    // members modal
    const [membersShow, setMembersShow] = useState(false);
    const toggleMembers = () => {
        setMembersShow(!membersShow);
    };

    // getting group membership
    useEffect(()=>{
        (
            async () => {
                if (props.group.ID === undefined) {
                    return
                }
                const response = await fetch("http://localhost:8080/api/group/membership?group=" + props.group.ID.toString(), {
                    headers: {"Content-Type": "application/json"},
                    credentials: "include",
                });
                const responseJSON = await response.json();
                setMember(responseJSON);
            }
        )();
    }, [props.group]);

    // function for sending message when submit
    const sendMessage = (e) => {
        e.preventDefault();
        if (msg === "") return false;
        props.ws.send(JSON.stringify({
            "group": props.group.ID,
            "member": member.ID,
            "text": msg,
            "nick": member.nick
        }));
        document.getElementById("text-area").value = "";
        document.getElementById("text-area").focus();
    }

    let nomessages = false;
    if (props.messages === []) {
        nomessages = true;
    }

    let load;
    if (props.group.ID === undefined) {
        load = <h1>Select a group to chat!</h1>;
    } else {
        load = (
            <div className="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
                <div className="selected-user row">
                    <span className="mr-auto mt-4">To: <span className="name">{props.group.name}</span></span>
                    <div className="dropdown">
                        <button className="btn btn-primary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            Settings
                        </button>
                        <GroupMenu toggleDel={toggleDelGroup} toggleAdd={toggleAddUser} toggleMembers={toggleMembers} member={member}/>
                    </div>
                </div>
                <div className="chat-container">
                    <ul className="chat-box chatContainerScroll" style={{height: '70vh', overflow: 'scroll'}}>
                        {nomessages?null:props.messages.map(item => {return <Message key={uuidv4()} time={item.created} message={item.text} name={item.nick} member={item.member} user={member.ID}/>})}
                    </ul>
                    <form id="chatbox" className="form-group mt-3 mb-0 d-flex column justify-content-center" onSubmit={sendMessage}>
                        <textarea autoFocus  id="text-area" className="form-control mr-1" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
                        <input className="btn btn-primary" type="submit" value="Send" />
                    </form>
                </div>
                <ModalDeleteGroup show={delGrShow} toggle={toggleDelGroup} group={props.group} setGroups={props.setGroups} groups={props.groups} setCurrent={props.setCurrent}/>
                <ModalAddUser show={addUserShow} toggle={toggleAddUser} group={props.group}/>
                <ModalMembers show={membersShow} toggle={toggleMembers} group={props.group} member={member} />
            </div>
        );
    }
    return load;
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
            <div className="chat-hour">{props.time} <span className="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {props.member===props.user?right:left}
        </div>
    )
}

export default Chat;