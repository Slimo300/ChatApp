const port = "8080";
const hostname = "localhost";

export class API{
    constructor() {
        this.axios = require('axios').default;
        this.axios.defaults.baseURL = 'http://'+hostname+':'+port+'/api/';
        this.axios.defaults.headers.common['Content-Type'] = "application/json";
    }

    SetAccessToken(accessToken) {
        this.axios.defaults.headers.common['Authorization'] = "Bearer " + accessToken;
    }

    async Register(email, username, password, rpassword){
        if (email.trim() === "") {
            throw new Error("Email can't be blank");
        }
        if (username.trim() === "") {
            throw new Error("Username can't be blank");
        }
        if (password.trim() === "") {
            throw new Error("Password can't be blank");
        }
        if (password !== rpassword) {
            throw new Error("Passwords don't match");
        }
        return await this.axios.post("/register", {
                username: username, 
                email: email,
                password: password,
            });
    }

    async Login(email, password) {
        if (email.trim() === "") {
            throw new Error("Email cannot be blank");
        }
        if (password.trim() === "") {
            throw new Error("Password cannot be blank");
        }
        return await this.axios.post("/login", {
            email: email,
            password: password,
        });
    }

    async GetUser() {
        return await this.axios.get('/user');
    }

    async GetInvites() {
        return await this.axios.get("/invites");
    }

    async GetGroups() {
        return await this.axios.get("/group");
    }
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

export async function ChangePassword(oldPassword, newPassword, repeatPassword) {
    if (newPassword === "") {
        return {"err": "password cannot be blank"}
    }
    if (newPassword.length <  6) {
        return {"err": "password must be at least 6 characters long"}
    }
    if (repeatPassword !== newPassword) {
        return {"err": "Passwords don't match"}
    }

    let response = await fetch('http://'+hostname+':'+port+"/api/change-password", {
        method: "PUT",
        credentials: "include",
        body: JSON.stringify({
            "oldPassword": oldPassword,
            "newPassword": newPassword,
        }),
    })

    let responseJSON = await response.json();

    return responseJSON;
}

export async function UpdateProfilePicture(data) {
    let response = await fetch('http://'+hostname+':'+port+"/api/set-image", {
        method: "POST",
        credentials: "include",
        body: data,
    })
    let responseJSON = await response.json();
    return responseJSON;
}

export async function DeleteProfilePicture() {
    let response = await fetch('http://'+hostname+":"+port+"/api/delete-image", {
        method: "DELETE",
        credentials: "include",
    })
    let responseJSON = await response.json();
    return responseJSON;
}

export async function UpdateGroupProfilePicture(data, groupID) {
    let response = await fetch('http://'+hostname+':'+port+"/api/group/"+groupID+"/image", {
        method: "POST",
        credentials: "include",
        body: data,
    })
    let responseJSON = await response.json();
    return responseJSON;
}

export async function DeleteGroupProfilePicture(groupID) {
    let response = await fetch('http://'+hostname+":"+port+"/api/group/"+groupID+"/image", {
        method: "DELETE",
        credentials: "include",
    })
    let responseJSON = await response.json();
    return responseJSON;
}

const APICaller = new API();

export default APICaller;


// export async function Register(email, username, password, rpassword) {
//     if (password.trim() === "") {
//         throw new Error("Password can't be blank");
//     }
//     if (password !== rpassword) {
//         throw new Error("Passwords don't match");
//     }
//     let response = await fetch('http://'+hostname+':'+port+'/api/register', {
//         method: 'POST',
//         headers: {'Content-Type': 'application/json'},
//         credentials: 'include',
//         body: JSON.stringify({
//             username: username, 
//             email: email,
//             password: password
//         })
//     });
//     let responseJSON = await response.json();

//     return responseJSON;
// }

// export async function Login(email, password) {
//     if (email.trim() === "") {
//         throw new Error("Email cannot be blank");
//     }
//     if (password.trim() === "") {
//         throw new Error("Password cannot be blank");
//     }

//     let response = await fetch('http://'+hostname+':'+port+'/api/login', {
//         method: 'POST',
//         headers: {'Content-Type': 'application/json'},
//         credentials: 'include',
//         body: JSON.stringify({
//             email: email,
//             password: password,
//         })
//     });

//     let responseJSON = await response.json();

//     return responseJSON;
// }


// export async function GetUser() {
//     const response = await fetch('http://'+hostname+':'+port+'/api/user', {
//         method: 'GET',
//         headers: {'Content-Type': 'application/json'},
//         credentials: 'include'});
//     if (response.status !== 200) {
//         throw new Error("couldn't get user");
//     }
//     const promise = await response.json();
//     return promise;  
// }

// export async function GetInvites() {
//     const response = await fetch('http://'+hostname+':'+port+'/api/invites', {
//         headers: {'Content-Type': 'application/json'},
//         credentials: 'include'
//     });
//     if (response.status !== 200 && response.status !== 204) {
//         throw new Error("Invalid response when requesting user invites");
//     }
//     if (response.status === 200) {
//         const promise = await response.json();
//         return promise
//     }
// }

// export async function GetGroups() {
//     const response = await fetch('http://'+hostname+':'+port+'/api/group', {
//         headers: {'Content-Type': 'application/json'},
//         credentials: 'include'});
//     if (response.status === 204) {
//         return {"err": "no groups"};
//     }
//     const promise = await response.json();
//     return promise;
// }