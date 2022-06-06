import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";
import { RespondInvite } from "../Requests";

const Invite = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const Respond = (answer) => {
        let result = RespondInvite(props.invite.ID, answer);
        result.then(response => {
            if (response === null) {
                alert("couldn't respond to invte");
                return;
            }
            dispatch({type: actionTypes.NEW_GROUP, payload: response});
            dispatch({type: actionTypes.DELETE_NOTIFICATION, payload: props.invite.ID});
            console.log(response);
        });
    };

    return (
        <div className="dropdown-item invite">
            <div className="list-group-item list-group-item-info d-flex row justify-content-around">
                <div className="chat-avatar image-holder-invite">
                    <img className="rounded-circle img-thumbnail"
                        src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.invite.issuer.pictureUrl}
                        onError={({ currentTarget }) => {
                            currentTarget.onerror = null; 
                            currentTarget.src="https://erasmuscoursescroatia.com/wp-content/uploads/2015/11/no-user.jpg";
                        }}
                    />
                </div>
                <div className="chat-name align-self-center">{props.invite.issuer.username}</div>
                <div className="align-self-center">invited you to </div>
                <div className="chat-avatar image-holder-invite">
                    <img className="rounded-circle img-thumbnail"
                        src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.invite.group.pictureUrl}
                        onError={({ currentTarget }) => {
                            currentTarget.onerror = null; 
                            currentTarget.src="https://cdn.icon-icons.com/icons2/3005/PNG/512/people_group_icon_188185.png";
                        }}
                    />
                </div>
                <div className="chat-name align-self-center">{props.invite.group.name}</div>
                <button className="btn-primary h-50 align-self-center" type="button" onClick={() => {Respond(true)}}>Accept</button>
                <button className="btn-secondary h-50 align-self-center" type="button" onClick={() => {Respond(false)}}>Decline</button>
            </div>
        </div>
    )
};

export default Invite;