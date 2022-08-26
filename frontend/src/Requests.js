const port = "8080";
const hostname = "localhost";

export class API{
    constructor() {
        this.axios = require('axios').default;
        this.axios.defaults.baseURL = 'http://'+hostname+':'+port+'/api/';
        this.axios.defaults.headers.common['Content-Type'] = "application/json";
        this.axios.defaults.withCredentials = true;
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

    async GetWebsocket() {
        let socket = new WebSocket('ws://'+hostname+':'+port+'/ws/');
        socket.onopen = () => {
            console.log("Websocket openned");
        };
        socket.onclose = () => {
            console.log("closed");
        };
        return socket;
    }

    async Logout() {
        return await this.axios.post("/signout", {}, {
            withCredentials: true,
        });
    }

    async LoadMessages(groupID, offset) {
        return await this.axios.get("/group/"+groupID+"/messages?num=8&offset="+offset);
    }

    async CreateGroup(name, desc) {
        return await this.axios.post("/group", {
            "name": name,
            "desc": desc,
        })
    }

    async DeleteGroup(groupID) {
        return await this.axios.delete("/group/"+groupID);
    }

    async SendGroupInvite(username, groupID) {
        return await this.axios.post("/invites", {
            "target": username,
            "group": groupID
        });
    }

    async RespondGroupInvite(inviteID, answer) {
        return await this.axios.put("/invites/"+inviteID, {
            "answer": answer
        })
    }

    async DeleteMember(memberID) {
        return await this.axios.delete("member/"+memberID);
    }

    async SetRights(memberID, adding, deleting, setting) {
        return await this.axios.put("/member/"+memberID, {
            "adding": adding,
            "deleting": deleting,
            "setting": setting
        });
    }

    async ChangePassword(oldPassword, newPassword, repeatPassword) {
        if (newPassword === "") {
            throw new Error("password cannot be blank");
        }
        if (newPassword.length <  6) {
            throw new Error("password must be at least 6 characters long");
        }
        if (repeatPassword !== newPassword) {
            throw new Error("Passwords don't match");
        }

        return await this.axios.put("/change-password", {
            "oldPassword": oldPassword,
            "newPassword": newPassword,
        });
    }

    async UpdateProfilePicture(image) {
        return await this.axios.post("/set-image", image, {
            headers: {
                'Content-Type': 'multipart/form-data',
            }
        })
    }
    
    async DeleteProfilePicture() {
        return await this.axios.delete("/delete-image");
    }

    async UpdateGroupProfilePicture(imageForm, groupID) {
        return await this.axios.post("/group/"+groupID+"/image", imageForm, {
            headers: {
                'Content-Type': 'multipart/form-data',
            }
        });
    }

    async DeleteGroupProfilePicture(groupID) {
        return await this.axios.delete("group/"+groupID+"/image");
    }
}


const APICaller = new API();

export default APICaller;


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
