import React, {useState} from 'react';
import {v4 as uuidv4} from "uuid";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';

export const ModalMembers = (props) => {

    const [err, setErr] = useState("");

    let message = null;
    if (err !== "") {
        message = <h5 className='text-danger'>{err}</h5>
    }

    let nogroup = false;
    if (props.group.Members === null) {
        nogroup = true
    }
    return (
        <Modal id="buy" tabIndex="-1" size='lg' role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Group Members
                </ModalHeader>
                <ModalBody>
                    <div>
                        {message}
                        <div className='form-group'>
                            <table className="table">
                                <tbody>
                                    {nogroup?null:props.group.Members.map((item) => {return <Member key={uuidv4()} member={item} setErr={setErr} toggle={props.toggle} user={props.member}/>})}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 

const Member = (props) => {
    const [adding, setAdding] = useState(props.member.adding);
    const toggleAdding = () => {
        setAdding(!adding);
    }
    const [deleting, setDeleting] = useState(props.member.deleting);
    const toggleDeleting = () => {
        setDeleting(!deleting);
    }
    const [setting, setSetting] = useState(props.member.setting);
    const toggleSetting = () => {
        setSetting(!setting);
    }

    const deleteMember = async() => {

        const response = await fetch('http://localhost:8080/api/group/remove', {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                "member": props.member.ID
            })
        });

        const responseJSON = await response.json();

        if (responseJSON.message === "ok") {
            props.setErr("Member deleted");
        } else {
            props.setErr(responseJSON.err);
            console.log(responseJSON.err);
        }
        setTimeout(function () {    
            props.toggle();
            props.setErr("");
        }, 2000);

    }

    const setRights = async() => {

        if (adding === props.member.adding && deleting === props.member.deleting && setting === props.member.setting) {
            return
        }

        const response = await fetch('http://localhost:8080/api/group/member/rights', {
            method: 'POST',
            headers: {"Content-Type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                "member": props.member.ID,
                "adding": adding,
                "deleting": deleting,
                "setting": setting,
            })
        });

        const responseJSON = await response.json();

        if (responseJSON.message === "ok") {
            props.setErr("Rights changed");
        } else {
            props.setErr(responseJSON.err);
        }
        setTimeout(function () {
            props.setErr("");
        }, 2000);

    }

    return (
        <tr className="chat-avatar">
            <td className='pr-3'><img src="https://www.bootdey.com/img/Content/avatar/avatar3.png" alt="Retail Admin"/></td>
            <td className="chat-name pr-3 align-middle">{props.member.nick}</td>
            {props.user.setting?<td className='align-middle'>
                <input className="form-check-input" type="checkbox" id="inlineCheckbox1" checked={adding} disabled={props.member.creator} onChange={toggleAdding}/>
                <label className="form-check-label" htmlFor="inlineCheckbox1">Adding</label>
            </td>:null}
            {props.user.setting?<td className='align-middle'>
                <input className="form-check-input" type="checkbox" id="inlineCheckbox2" checked={deleting} disabled={props.member.creator} onChange={toggleDeleting}/>
                <label className="form-check-label" htmlFor="inlineCheckbox2">Deleting</label>
            </td>:null}
            {props.user.setting?<td className='align-middle'>
                <input className="form-check-input" type="checkbox" id="inlineCheckbox3" checked={setting} disabled={props.member.creator} onChange={toggleSetting}/>
                <label className="form-check-label" htmlFor="inlineCheckbox3">Setting</label>
            </td>:null}
            {props.user.setting?<td className='pr-3'><button className='btn-primary btn' disabled={props.member.creator} onClick={setRights}>Set rights</button></td>:null}
            {props.user.deleting?<td className='pr-3'><button className='btn-primary btn' disabled={props.member.creator} onClick={deleteMember}>Delete</button></td>:null}
        </tr>
    );
};