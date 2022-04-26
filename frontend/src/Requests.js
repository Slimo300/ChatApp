export async function GetUser() {
    const response = await fetch('http://localhost:8080/api/user', {
        method: 'GET',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include'});
    if (response.status !== 200) {
        throw new Error("couldn't get user");
    }
    const promise = response.json();
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
        const promise = response.json();
        return promise
    }
}