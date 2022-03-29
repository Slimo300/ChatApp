import React from "react";

const GroupMenu = (props) => {

    return (
        <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
            <button className="dropdown-item">Options</button>
            <button className="dropdown-item" onClick={props.toggleMembers}>Members</button>
            <button className="dropdown-item" onClick={props.toggleAdd} disabled={!props.member.adding}>Add User</button>
            <div className="dropdown-divider"></div>
            <button className="dropdown-item" onClick={props.toggleDel} disabled={!props.member.creator}>Delete Group</button>
        </div>
    );
};

export default GroupMenu;