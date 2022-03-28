import React from "react";

export const GroupLabel = (props) => {
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
