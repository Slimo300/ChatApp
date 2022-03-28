import React from "react";

const GroupMenu = (props) => {

    return (
        <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
            <button className="dropdown-item" href="#">Options</button>
            <button className="dropdown-item" href="#">Members</button>
            <button className="dropdown-item" href="#">Add User</button>
            <div className="dropdown-divider"></div>
            <button className="dropdown-item" onClick={props.toggle}>Delete Group</button>
        </div>
    );
};

export default GroupMenu;