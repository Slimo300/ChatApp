
const port = "8080";

let hostname;

if (window.location.host ==="localhost:3000") {
    hostname="localhost";
}else{
    hostname="backend";
}

export async function GetUser() {
    const response = await fetch('http://'+hostname+':'+port+'/api/user', {
        method: 'GET',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include'});
    if (response.status !== 200) {
        throw new Error("couldn't get user");
    }
    const promise = await response.json();
    return promise;  
}

export async function GetInvites() {
    const response = await fetch('http://'+hostname+':'+port+'api/invites', {
        headers: {'Content-Type': 'application/json'},
        credentials: 'include'
    });
    if (response.status !== 200 && response.status !== 204) {
        throw new Error("Invalid response when requesting user invites");
    }
    if (response.status === 200) {
        const promise = await response.json();
        return promise
    }
}

export async function GetGroups() {
    const response = await fetch('http://'+hostname+':'+port+'/api/group', {
        headers: {'Content-Type': 'application/json'},
        credentials: 'include'});
    if (response.status !== 200 && response.status !== 204 ) {
        throw new Error("Invalid response when requesting user groups");
    }
    const promise = await response.json();
    return promise;
}

export async function GetWebsocket() {
    let socket = new WebSocket('ws://'+hostname+':'+port+'/ws/')
    socket.onopen = () => {
        console.log("Websocket openned");
    };
    socket.onclose = () => {
        console.log("closed");
    };
    return socket
}

export async function Login(email, password) {
    let response;
    let error;
    if (email.trim() === "") {
        throw new Error("Email cannot be blank");
    }
    if (password.trim() === "") {
        throw new Error("Password cannot be blank");
    }

    response = await fetch('http://'+hostname+':'+port+'/api/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify({
            email,
            password,
        })
    });

    if (response.status !== 200) {
        error = "Error occured, status code: "+response.status;
    }
   
    let responseJSON = await response.json();
    return {response: responseJSON, error: error};
}

export async function Logout() {
    try {
        let response = await fetch('http://'+hostname+':'+port+'/api/signout', {
            method: "POST",
            credentials: "include",
            headers: {"Content-Type": "application/json"}
        });
        if (response.status !== 200) {
            throw new Error("Logout unsuccesful");
        }
    } catch(err) {
        return err;
    }
}

export async function Register(email, username, password, rpassword) {
    if (password.trim() === "") {
        throw new Error("Password can't be blank");
    }
    if (password !== rpassword) {
        throw new Error("Passwords don't match");
    }
    let response = await fetch('http://'+hostname+':'+port+'/api/register', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify({
            username, 
            email,
            password
        })
    });
    let responseJSON = await response.json();

    return responseJSON;
}

export async function LoadMessages(groupID, offset) {
    let messages;
    try {
        let response = await fetch('http://'+hostname+':'+port+'/api/group/'+groupID+'/messages?num=8&offset=' + offset, {
            headers: {"Content-Type": "application/json"},
            credentials: "include",
        });
        if (response.status === 200) {
            messages = await response.json();
        }
        else if (response.status === 204) {
            messages = [];
        } 
        else {
            throw new Error("getting messages failed with status code: " + response.status.toString());
        } 
    } catch(err) {
        return err;
    }
    return messages;
}

export async function CreateGroup(name, desc) {
    let response;
    try {
        response = await fetch('http://'+hostname+':'+port+'/api/group', {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                "name": name,
                "desc": desc,
            })
        });
        if (response.status !== 201) {
            throw new Error("Invalid status code: " + response.status.toString());
        }
    } catch(err) {
        return err;
    }
    let newGroup = await response.json();
    return newGroup;
}

export async function DeleteGroup(groupID) {
    let response;
    try {
        response = await fetch('http://'+hostname+':'+port+'/api/group/'+groupID, {
            method: "DELETE",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
        });

        if (response.status !== 200) {
            throw new Error("Invalid status code: " + response.status.toString())
        }
    } catch(err) {
        return err;
    }
    return null;
}

export async function DeleteMember(memberID) {
    let response;
    try {
        response = await fetch('http://'+hostname+':'+port+'/api/member/'+memberID, {
            method: 'DELETE',
            headers: {"Content-Type": "application/json"},
            credentials: "include"
        });
        if (response.status !== 200) {
            throw new Error('invalid status code: ' + response.status.toString());
        }
    } catch(err) {
        return err;
    }
    let responseJSON = await response.json();
    return responseJSON;
}

export async function SetRights(memberID, adding, deleting, setting) {
    let response;
    try {
        response = await fetch('http://'+hostname+':'+port+'/api/member/'+memberID, {
            method: 'PUT',
            headers: {"Content-Type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                "adding": adding,
                "deleting": deleting,
                "setting": setting,
            })
        });
        if (response.status !== 200) {
            throw new Error("Invalid status code: " + response.status.toString());
        }
    } catch(err) {
        return err;
    }
    let responseJSON = await response.json();
    return responseJSON;
}

export async function RespondInvite(inviteID, answer) {
    let response;

    response = await fetch('http://'+hostname+':'+port+'/api/invites/'+inviteID, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        credentials: "include",
        body: JSON.stringify({
            "answer": answer,
        }),
    });
    if (response.status !== 200) {
        let responseJSON = await response.json();
        console.log(responseJSON.err);
        return null;
    }
    let responseJSON = await response.json();
    return responseJSON;
}  