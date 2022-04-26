import React, {createContext, useReducer} from "react";

export const StorageContext = createContext({});

const initialState = {
    groups: [],
    notifications: [],
    user: {},
};

export const actionTypes = {
    LOGIN: "LOGIN",
    SET_GROUPS: "SET_GROUPS",
    NEW_GROUP: "NEW_GROUP",
    DELETE_GROUP: "DELETE_GROUP",
    ADD_MEMBER: "ADD_MEMBER",
    DELETE_MEMBER: "DELETE_MEMBER",
    SET_MESSAGES: "SET_MESSAGES",
    ADD_MESSAGE: "ADD_MESSAGE",
    ADD_MESSAGES: "ADD_MESSAGES",
    SET_NOTIFICATIONS: "SET_NOTIFICATIONS",
    ADD_NOTIFICATION: "ADD_NOTIFICATION",
    DELETE_NOTIFICATION: "DELETE_NOTIFICATION"
}

function reducer(state, action) {
    switch (action.type) {
        case actionTypes.LOGIN:
            return Login(state, action.payload);
        case actionTypes.SET_GROUPS:
            return SetGroups(state, action.payload);
        case actionTypes.NEW_GROUP:
            return NewGroup(state, action.payload);
        case actionTypes.DELETE_GROUP:
            return DeleteGroup(state, action.payload);
        case actionTypes.ADD_MEMBER:
            return AddMemberToGroup(state, action.payload);
        case actionTypes.DELETE_MEMBER:
            return DeleteMemberFromGroup(state, action.payload);
        case actionTypes.SET_MESSAGES:
            return SetMessages(state, action.payload);
        case actionTypes.ADD_MESSAGE:
            return AddMessage(state, action.payload);
        case actionTypes.ADD_MESSAGES:
            return AddMessages(state, action.payload);
        case actionTypes.SET_NOTIFICATIONS:
            return SetInvites(state, action.payload);
        case actionTypes.ADD_NOTIFICATION:
            return AddInvite(state, action.payload);
        case actionTypes.DELETE_NOTIFICATION:
            return DeleteInvite(state, action.payload);
        default:
            throw new Error("Action not specified");
    }
}

const ChatStorage = ({children}) => {

    const [state, dispatch] = useReducer(reducer, initialState);

    return (
        <StorageContext.Provider value={[state, dispatch]}>
            {children}
        </StorageContext.Provider>
    );
}
export default ChatStorage;

function Login(state, payload) {
    let newState = {...state};
    newState.user = payload;
    return newState;
}

function SetGroups(state, payload) {
    let newState = {...state};
    newState.groups = payload;
    return newState;
}

function NewGroup(state, payload) {
    let newState = {...state};
    newState.groups = [...newState.groups, payload];
    return newState;
}

function DeleteGroup(state, payload) {
    let newState = {...state};
    newState.groups = newState.groups.filter( (item) => { return item.ID !== payload } );
    return newState;
}

function AddMemberToGroup(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group_id) {
            newState.groups[i].Members.push(payload);
            return newState;
        }
    }
    throw new Error("Group not found");
}

function DeleteMemberFromGroup(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group_id) {
            newState.groups[i].Members = newState.groups[i].Members.filter((item)=>{return item.ID !== payload.ID});
            return newState
        }
    }
    throw new Error("Group not found");
}

function SetMessages(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group) {
            newState.groups[i].messages = payload.messages
            return newState;
        }
    }
    throw new Error("Received messages don't belong to any of your groups");
}

function AddMessage(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group) {
            newState.groups[i].messages = [...newState.groups[i].messages, payload];
            return newState;
        }
    }
    throw new Error("Received message don't belong to any of your groups");
}

function AddMessages(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group) {
            newState.groups[i].messages = [...payload.messages, ...newState.groups[i].messages];
            return newState;
        }
    }
    throw new Error("Received messages don't belong to any of your groups");
}

function SetInvites(state, payload) {
    let newState = {...state};
    newState.notifications = payload;
    return newState;
}

function AddInvite(state, payload) {
    let newState = {...state};
    newState.notifications = [...newState.notifications, payload];
    return newState;
}

function DeleteInvite(state, payload) {
    let newState = {...state};
    newState.notifications = newState.notifications.filter( (item) => { return item.ID !== payload } );
    return newState;
}