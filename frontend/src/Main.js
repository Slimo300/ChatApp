import React, { useEffect, useState } from "react";
import {Navigate} from "react-router-dom";

const Main = (props) => {

    console.log(props.name)
    return (
        <div>
            {props.name === ""? <Navigate to="/login" />:<AuthMain name={props.name}/>}
        </div>
    );
}

const AuthMain = (props) => {

    const [groups, setGroups] = useState([]);
    const [current, setCurrent] = useState(0);
    const [messages, setMessages] = useState([]);
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
    useEffect(()=>{
        (
            async () => {
                if (current !== 0) {
                    const response = await fetch("http://localhost:8080/api/group/messages?group=" + current.toString(), {
                        headers: {"Content-Type": "application/json"},
                        credentials: "include",
                    });
                    const responseJSON = await response.json();
                    console.log(responseJSON);
                    setMessages(responseJSON);
                    // setCurrent(0);
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
                                            {groups.map(item => {return <GroupLabel name={item.name} key={item.ID} id={item.ID} setCurrent={setCurrent}/>})}
                                        </ul>
                                    </div>
                                </div>
                                <Chat messages={messages} group={current}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

const GroupLabel = (props) => {
    const change = () => {props.setCurrent(props.id)};
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

const Chat = (props) => {
    const [member, setMember] = useState(0);
    const getMember = async () => {
        const response = await fetch("http://localhost:8080/api/group/membership?group=" + props.group.toString(), {
            headers: {"Content-Type": "application/json"},
            credentials: "include",
        });
        const responseJSON = await response.json();
        if (responseJSON.ID !== member) {
            setMember(responseJSON.ID);
            console.log(responseJSON);
        }
    }
    let load;

    if (props.group === 0) {
        load = <h1>Select a group to chat!</h1>;
    } else {
        getMember();
        
        load = (
            <div className="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
                <div className="selected-user">
                    <span>To: <span className="name">{props.name}</span></span>
                </div>
                <div className="chat-container">
                    <ul className="chat-box chatContainerScroll">
                        {props.messages.map(item => {return <Message key={item.ID} time={item.posted} message={item.text} name={item.Member.Nick} member={item.member_id} user={member}/>})}
                    </ul>
                    <div className="form-group mt-3 mb-0">
                        <textarea className="form-control" rows="3" placeholder="Type your message here..."></textarea>
                    </div>
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
            <div className="chat-hour">{props.time} <span class="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {props.member===props.user?left:right}
        </div>
    )
}

export default Main;