import React, { useState } from "react";

const Chat = (props) => {

    const [member, setMember] = useState(0);
    const [msg, setMsg] = useState("");

    const getMember = async () => {
        const response = await fetch("http://localhost:8080/api/group/membership?group=" + props.group.toString(), {
            headers: {"Content-Type": "application/json"},
            credentials: "include",
        });
        const responseJSON = await response.json();
        if (responseJSON.ID !== member) {
            setMember(responseJSON.ID);
        }
    }

    const submit = (e) => {
        e.preventDefault();
        if (msg === "") return;
        if (!props.socket) {
            alert("Error: There is no socket connection.");
            return;
        }
        console.log("send");
        props.socket.send(JSON.stringify({
            "group": props.group,
            "member": member,
            "message": msg
        }));
        console.log("sent");
    }
    let load;

    if (props.group === 0) {
        load = <h1>Select a group to chat!</h1>;
    } else {
        getMember();
        
        load = (
            <div className="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
                <div className="selected-user">
                    <span>To: <span className="name">{props.groupname}</span></span>
                </div>
                <div className="chat-container">
                    <ul className="chat-box chatContainerScroll">
                        {props.messages.map(item => {return <Message key={item.ID} time={item.posted} message={item.text} name={item.Member.Nick} member={item.member_id} user={member}/>})}
                    </ul>
                    <form id="chatbox" className="form-group mt-3 mb-0" onSubmit={submit}>
                        <textarea className="form-control" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
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