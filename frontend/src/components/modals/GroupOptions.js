import React, { useContext, useState } from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { DeleteGroupProfilePicture, UpdateGroupProfilePicture } from "../../Requests";
import { actionTypes, StorageContext } from '../../ChatStorage';

export const ModalGroupOptions = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const [file, setFile] = useState({});

    const [message, setMessage] = useState("");

    const updatePicture = (e) => {
        e.preventDefault();
        let data = new FormData();
        data.append("avatarFile", file);
    
        let response = UpdateGroupProfilePicture(data, props.group.ID);
        response.then(resp => {
            if (resp.err !== undefined) {
                setMessage(resp.err);
            } else {
                setMessage("image uploaded");
                dispatch({type: actionTypes.NEW_GROUP_PROFILE_PICTURE, payload: {newURI: resp.newUrl, groupID: props.group.ID}});
                let timestamp = new Date().getTime();
                document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.group.pictureUrl+"?"+timestamp;
                document.getElementById("customFile").value= null;
            }
            setTimeout(function() {
                setMessage("");
            }, 3000);
        })
    };

    const deletePicture = (e) => {
        let response = DeleteGroupProfilePicture(props.group.ID)
        response.then(resp => {
            if (resp.err !== undefined) {
                setMessage(resp.err);
            } else {
                setMessage("image deleted");
                dispatch({type: actionTypes.DELETE_GROUP_PROFILE_PICTURE, payload: props.group.ID})
                let timestamp = new Date().getTime();
                document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.group.ID+"?"+timestamp;
            }
            setTimeout(function() {
                setMessage("");
            }, 3000);
        })
    };

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    User Profile
                </ModalHeader>
                <ModalBody>
                    <div class="container">
                        <div className="row d-flex justify-content-center">
                            <div className="text-center card-box">
                                <div className="member-card">
                                    {message}
                                    <div className="mx-auto profile-image-holder">
                                        <img id="profilePicture" className="rounded-circle img-thumbnail"
                                            src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.group.pictureUrl}
                                            onError={({ currentTarget }) => {
                                                currentTarget.onerror = null; 
                                                currentTarget.src="https://cdn.icon-icons.com/icons2/3005/PNG/512/people_group_icon_188185.png";
                                            }}
                                        />
                                    </div>
                                    <div>
                                        <h4>{props.name}</h4>
                                    </div>
                                    <hr />
                                    <h3>Change profile picture</h3>
                                    <form encType="multipart/form-data" action="">
                                        <input type="file" className="form-control" id="customFile" onChange={e => setFile(e.target.files[0])}  />
                                        <div className="text-center mt-2">
                                            <button className="btn btn-primary text-center w-100" onClick={updatePicture}>Upload</button>
                                        </div>
                                    </form>
                                    <div className="text-center mt-4">
                                        <button className="btn btn-danger text-center w-100" onClick={deletePicture}>Delete Picture</button>
                                    </div>
                                    <hr />
                                </div>
                            </div>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 