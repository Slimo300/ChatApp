import React from "react";

export const GroupLabel = (props) => {
    const change = () => {
        props.setCurrent(props.group);
        let newCounter = props.counter;
        newCounter[props.group.ID] = 0;
        props.setCounter(newCounter);
    };
    console.log(props.counter[props.group.ID]);
    return (
        <li className="person" onClick={change}>
            <div className="user">
                <img src="https://www.bootdey.com/img/Content/avatar/avatar3.png" alt="Retail Admin"/>
            </div>
            <p className="name-time">
                <span className="name">{props.group.name}</span>
            </p>
            {props.counter[props.group.ID]>0?<span className="badge badge-primary float-right">{props.counter[props.group.ID]}</span>:null}
        </li>
    );
}
