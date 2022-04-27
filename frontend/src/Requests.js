export async function GetUser() {
    const response = await fetch('http://localhost:8080/api/user', {
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
    const response = await fetch('http://localhost:8080/api/invites', {
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
    const response = await fetch('http://localhost:8080/api/group/get', {
        headers: {'Content-Type': 'application/json'},
        credentials: 'include'});
    if (response.status !== 200 && response.status !== 204 ) {
        throw new Error("Invalid response when requesting user groups");
    }
    const promise = await response.json();
    return promise;
}

export async function GetMessages(groupID) {
    const response = await fetch("http://localhost:8080/api/group/messages?group=" + groupID + "&num=8", {
        headers: {"Content-Type": "application/json"},
        credentials: "include",
    });
    let messages;
    if (response.status === 200) {
        messages = await response.json();
    }
    else if (response.status === 204) {
        messages = [];
    } 
    else {
        throw new Error("getting messages failed with status code: " + response.status.toString());
    }
    return messages;
}

export async function GetWebsocket() {
    let socket = new WebSocket("ws://localhost:8080/ws")
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
    try {
        if (email.trim() === "") {
            throw new Error("Email cannot be blank");
        }
        if (password.trim() === "") {
            throw new Error("Password cannot be blank");
        }

        response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                email,
                password,
            })
        });

        if (response.status !== 200) {
            throw new Error("Invalid status code: "+ response.status.toString());
        }
    } catch(err) {
        return err;
    }
    let promise = await response.json();
    return promise;
}

export async function Logout() {
    try {
        let response = await fetch("http://localhost:8080/api/signout", {
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
    try {
        if (password.trim() === "") {
            throw new Error("Password can't be blank");
        }
        if (password !== rpassword) {
            throw new Error("Passwords don't match");
        }
        let response = await fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                username, 
                email,
                password
            })
        });
        if ( response.status !==  201) {
            throw new Error("Invalid status code: " + response.status.toString());
        }
    }
    catch(err) {
        return err;
    }
    return null;
}

export async function LoadMessages(groupID, offset) {
    let messages;
    try {
        let response = await fetch("http://localhost:8080/api/group/messages?group=" + groupID + "&num=8&offset=" + offset, {
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
        response = await fetch('http://localhost:8080/api/group/create', {
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
        response = await fetch('http://localhost:8080/api/group/delete', {
            method: "DELETE",
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                "group": groupID
            })
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
        response = await fetch('http://localhost:8080/api/group/remove', {
            method: 'PUT',
            headers: {"Content-Type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                "member": memberID
            })
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
        response = await fetch('http://localhost:8080/api/group/rights', {
            method: 'PUT',
            headers: {"Content-Type": "application/json"},
            credentials: "include",
            body: JSON.stringify({
                "member": memberID,
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