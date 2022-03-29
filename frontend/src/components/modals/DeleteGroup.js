import React, {useState} from 'react';
import { Modal, ModalHeader, ModalBody } from 'reactstrap';

export const ModalDeleteGroup = (props) => {

    const [err, setErr] = useState("");

    const submit = async() => {
        const response = await fetch('http://localhost:8080/api/group/delete', {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                "group": props.group.ID
            })
        });

        const responseJSON = await response.json()

        let flag = false;

        if (responseJSON.message === "ok") {
            props.setGroups(props.groups.filter((item)=> { return item.ID !== props.group.ID }));
            setErr("Group Deleted");
            flag = true;
        }
        else {
            setErr(responseJSON.err);
        }
        setTimeout(function () {    
            props.toggle();
            setErr("");
            if (flag) {
                props.setCurrent({});
            }
        }, 1000);
    }

    var message = null;
    if (err !== "") {
        message = <h5 className="mb-4 text-danger">{err}</h5>
    }

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Delete Group
                </ModalHeader>
                <ModalBody>
                    <div>
                        {message}
                        <div className='form-group'>
                            <label>Are you sure you want to delete group {props.group.name}?</label>
                        </div>
                        <div className="form-row text-center">
                            <div className="col-12 mt-2">
                                <button className="btn btn-dark btn-large text-primary" onClick={submit}>Delete</button>
                            </div>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 
