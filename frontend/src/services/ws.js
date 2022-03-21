export const ws = new WebSocket("ws://localhost:8080/ws")
ws.onopen = () => {
    console.log("Websocket openned");
};
ws.onclose = () => {
    console.log("closed");
};

// export const connect = () => {
//     let ws = new WebSocket("ws://localhost:8080/ws");
//     ws.onopen = () => {
//         console.log("Websocket openned");
//     };
//     ws.onclose = () => {
//         console.log("closed");
//     };

//     return ws
// };