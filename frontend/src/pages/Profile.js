import React, { useContext, useState } from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { ChangePassword, DeleteProfilePicture, UpdateProfilePicture } from "../Requests";
import { actionTypes, StorageContext } from '../ChatStorage';

export const ModalUserProfile = (props) => {

    const [state, dispatch] = useContext(StorageContext);

    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [repeatPassword, setRepeatPassword] = useState("");

    const [file, setFile] = useState({});

    const [message, setMessage] = useState("");

    const changePassword = (e) => {
            e.preventDefault()
            let response = ChangePassword(oldPassword, newPassword, repeatPassword);
            response.then(resp => {
                document.getElementById("oldpassword").value = "";
                document.getElementById("newpassword").value = "";
                document.getElementById("rpassword").value = "";
                if (resp.err !== undefined) {
                    setMessage(resp.err);
                } else {
                    setMessage(resp.message);
                }
                setTimeout(function() {
                    setMessage("");
                }, 3000);
            })
    }

    const updatePicture = (e) => {
        e.preventDefault();
        let data = new FormData();
        data.append("avatarFile", file);
    
        let response = UpdateProfilePicture(data);
        response.then(resp => {
            if (resp.err !== undefined) {
                setMessage(resp.err);
            } else {
                setMessage("image uploaded");
                dispatch({type: actionTypes.NEW_PROFILE_PICTURE, payload: resp.newUrl});
                let timestamp = new Date().getTime();
                document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+state.user.pictureUrl+"?"+timestamp;
                document.getElementById("customFile").value= null;
            }
            setTimeout(function() {
                setMessage("");
            }, 3000);
        })
    };

    const deletePicture = (e) => {
        let response = DeleteProfilePicture()
        response.then(resp => {
            if (resp.err !== undefined) {
                setMessage(resp.err);
            } else {
                setMessage("image deleted");
                dispatch({type: actionTypes.DELETE_PROFILE_PICTURE})
                let timestamp = new Date().getTime();
                document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+state.user.pictureUrl+"?"+timestamp;
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
                                            src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+state.user.pictureUrl}
                                            onError={({ currentTarget }) => {
                                                currentTarget.onerror = null; 
                                                currentTarget.src="https://erasmuscoursescroatia.com/wp-content/uploads/2015/11/no-user.jpg";
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
                                    <form className="mt-4">
                                        <h3> Change Password </h3>
                                        <div className="mb-3 text-center">
                                            <label htmlFor="pass" className="form-label">Old Password</label>
                                            <input name="oldpassword" type="password" className="form-control" id="oldpassword" onChange={(e) => setOldPassword(e.target.value)} />
                                        </div>
                                        <div className="mb-3 text-center">
                                            <label htmlFor="pass" className="form-label">New Password</label>
                                            <input name="newpassword" type="password" className="form-control" id="newpassword" onChange={(e) => setNewPassword(e.target.value)} />
                                        </div>
                                        <div className="mb-3 text-center">
                                            <label htmlFor="pass" className="form-label">Repeat Password</label>
                                            <input name="rpassword" type="password" className="form-control" id="rpassword" onChange={(e) => setRepeatPassword(e.target.value)} />
                                        </div>
                                            
                                        <div className="text-center">
                                            <button className="btn btn-primary text-center w-100" onClick={changePassword}>Change password</button>
                                        </div>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 