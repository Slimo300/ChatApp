import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";
import { RespondInvite } from "../Requests";

const Invite = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const Respond = (answer) => {
        let result = RespondInvite(props.inviteID, answer);
        result.then(response => {
            if (response === null) {
                alert("couldn't respond to invte");
                return;
            }
            dispatch({type: actionTypes.NEW_GROUP, payload: response});
            dispatch({type: actionTypes.DELETE_NOTIFICATION, payload: props.inviteID});
            console.log(response);
        });
    };

    return (
        <div className="dropdown-item">
            <div className="list-group-item list-group-item-info d-flex row justify-content-around">
                <div>{props.issID} invited you to {props.groupID}</div>
                <button className="btn-primary" type="button" onClick={() => {Respond(true)}}>Accept</button>
                <button className="btn-secondary" type="button" onClick={() => {Respond(false)}}>Decline</button>
            </div>
        </div>
    )
};

export default Invite;