import React, {useState} from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';

export const ModalAddFriend = (props) => {

    const [username, setUsername] = useState("");
    const [err, setErr] = useState("");

    const submit = async() => {
        const response = await fetch('http://localhost:8080/api/group/add', {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: {
                "name": username,
            }
        });
        const responseJSON = await response.json();

        if (responseJSON.err !== undefined) {
            setErr(responseJSON.err)
        }
    }

    let message = null;
    if (err !== "") {
        message = <h5 className="mb-4 text-danger">{err}</h5>;
    }

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Add Friend
                </ModalHeader>
                <ModalBody>
                    <div>
                        {message}
                        <form onSubmit={submit}>
                            <div className="form-group">
                                <label htmlFor="email">Username:</label>
                                <input name="name" type="text" className="form-control" id="gr_name" onChange={(e)=>{setUsername(e.target.value)}}/>
                            </div>
                            <div className="form-row text-center">
                                <div className="col-12 mt-2">
                                    <button type="submit" className="btn btn-dark btn-large text-primary">Add Friend</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 

// const AddFriendForm = () => {
//     const [username, setUsername] = useState("");
//     const [err, setErr] = useState("");

//     const submit = async() => {
//         const response = await fetch('http://localhost:8080/api/group/add', {
//             method: "POST",
//             headers: {'Content-Type': 'application/json'},
//             credentials: 'include',
//             body: {
//                 "name": username,
//             }
//         });
//         const responseJSON = await response.json();

//         if (responseJSON.err !== undefined) {
//             setErr(responseJSON.err)
//         }
//     }

//     let message = null;
//     if (err !== "") {
//         message = <h5 className="mb-4 text-danger">{err}</h5>;
//     }

//     return (
//         <div>
//             {message}
//             <form onSubmit={submit}>
//                 <div className="form-group">
//                     <label htmlFor="email">Username:</label>
//                     <input name="name" type="text" className="form-control" id="gr_name" onChange={(e)=>{setUsername(e.target.value)}}/>
//                 </div>
//                 <div className="form-row text-center">
//                     <div className="col-12 mt-2">
//                         <button type="submit" className="btn btn-dark btn-large text-primary">Add Friend</button>
//                     </div>
//                 </div>
//             </form>
//         </div>
//     );
// }
