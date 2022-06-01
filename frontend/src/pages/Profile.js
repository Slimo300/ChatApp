import React, { useState } from "react";
import { Navigate } from "react-router-dom";

const AuthProfile = (props) => {

    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [repeatPassword, setRepeatPassword] = useState("");

    return (
        <div class="container">
            <div className="row d-flex justify-content-center">
                <div className="text-center card-box">
                    <div className="member-card pt-2 pb-2">
                        <div className="mx-auto profile-image-holder">
                            <img className="rounded-circle img-thumbnail"
                                src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.name+".jpeg"}
                                onError={({ currentTarget }) => {
                                    currentTarget.onerror = null; // prevents looping
                                    currentTarget.src="https://erasmuscoursescroatia.com/wp-content/uploads/2015/11/no-user.jpg";
                                }}
                            />
                        </div>
                        <div>
                            <h4>{props.name}</h4>
                        </div>
                        <hr />
                        <h3>Change profile picture</h3>
                        <form>
                            <input type="file" className="form-control" id="customFile" />   
                            <div className="text-center mt-4">
                                <button type="submit" className="btn btn-primary text-center w-100">Upload</button>
                            </div>
                        </form>
                        <div className="text-center mt-4">
                            <button type="submit" className="btn btn-danger text-center w-100">Delete Picture</button>
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
                                <button type="submit" className="btn btn-primary text-center w-100">Submit</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>);
};


const Profile = (props) => {

    return (
        <div>
            {props.name === ""? <Navigate to="/login" />:<AuthProfile {...props}/>}
        </div>
    );
};

export default Profile;