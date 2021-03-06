import React, {useContext, useState} from 'react';
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import { actionTypes, StorageContext } from '../../ChatStorage';
import { DeleteGroup } from '../../Requests';

export const ModalDeleteGroup = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const [err, setErr] = useState("");

    const submit = async() => {
        let promise = DeleteGroup(props.group.ID);
        let flag = false;
        promise.then( response => { 
            if (response === null){
                dispatch({type: actionTypes.DELETE_GROUP, payload: props.group.ID})
                setErr("Group Deleted");
                flag = true;
            } else {
                setErr(response.message);
            }
            setTimeout(function () {    
                props.toggle();
                setErr("");
                if (flag) {
                    props.setCurrent({});
                }
            }, 1000);
         } );
    }

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Delete Group
                </ModalHeader>
                <ModalBody>
                    <div>
                        {err!==""?<h5 className="mb-4 text-danger">{err}</h5>:null}
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
