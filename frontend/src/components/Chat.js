import React, { useEffect, useState } from "react";
import {v4 as uuidv4} from "uuid";
import {ws} from "../services/ws"

const Chat = (props) => {

    const [member, setMember] = useState(0);
    const [msg, setMsg] = useState("");

    useEffect(()=>{
        (
            async () => {
                if (props.group === 0) {
                    return
                }
                const response = await fetch("http://localhost:8080/api/group/membership?group=" + props.group.toString(), {
                    headers: {"Content-Type": "application/json"},
                    credentials: "include",
                });
                const responseJSON = await response.json();
                setMember(responseJSON.ID);
            }
        )();
    }, [props.group]);

    const submit = (e) => {
        e.preventDefault();
        if (msg === "") return false;
        ws.send(JSON.stringify({
            "group": props.group,
            "member": member,
            "text": msg
        }));
        console.log("sent");
        document.getElementById("text-area").value = "";
        document.getElementById("text-area").focus();
    }

    let nomessages = false;
    console.log(props.messages);
    if (props.messages === []) {
        nomessages = true;
    }

    let load;

    if (props.group === 0) {
        load = <h1>Select a group to chat!</h1>;
    } else {
        load = (
            <div className="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
                <div className="selected-user">
                    <span className="mr-auto">To: <span className="name">{props.groupname}</span></span>
                    <button type="button" className="btn btn-primary float-right mt-3">Settings</button>
                </div>
                <div className="chat-container">
                    <ul className="chat-box chatContainerScroll">
                        {nomessages?null:props.messages.map(item => {return <Message key={uuidv4()} time={item.created} message={item.text} name={item.nick} member={item.member} user={member}/>})}
                    </ul>
                    <form id="chatbox" className="form-group mt-3 mb-0" onSubmit={submit}>
                        <textarea autoFocus  id="text-area" className="form-control" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
                        <input type="submit" value="Send" />
                    </form>
                </div>
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