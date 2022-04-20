import React, {createContext, useReducer} from "react";

export const StorageContext = createContext({});

export const actionTypes = {
    LOGIN: "LOGIN",
    SET_GROUPS: "SET_GROUPS",
    NEW_GROUP: "NEW_GROUP",
    DELETE_GROUP: "DELETE_GROUP",
    ADD_MEMBER: "ADD_MEMBER",
    DELETE_MEMBER: "DELETE_MEMBER",
    SET_MESSAGES: "SET_MESSAGES",
    ADD_MESSAGE: "ADD_MESSAGE",
    ADD_MESSAGES: "ADD_MESSAGES"
}

function reducer(state, action) {
    let newState;
    switch (action.type) {
        case actionTypes.LOGIN:
            newState = {...state};
            newState.user = action.payload;
            return newState;
        case actionTypes.SET_GROUPS:
            newState = {...state};
            newState.groups = action.payload;
            return newState;
        case actionTypes.NEW_GROUP:
            newState = {...state};
            newState.groups = [...newState.groups, action.payload];
            return newState;
        case actionTypes.DELETE_GROUP:
            newState = {...state};
            newState.groups = newState.groups.filter( (item) => {return item.ID !== action.payload} );
            return newState;
        case actionTypes.ADD_MEMBER:
            newState = {...state};
            for (let i = 0; i < newState.groups.length; i++) {
                if (newState.groups[i].ID === action.payload.group_id) {
                    newState.groups[i].Members.push(action.payload);
                    return newState;
                }
            }
            break;
        case actionTypes.DELETE_MEMBER:
            newState = {...state};
            for (let i = 0; i < newState.groups.length; i++) {
                if (newState.groups[i].ID === action.payload.group_id) {
                    newState.groups[i].Members = newState.groups[i].Members.filter((item)=>{return item.ID !== action.payload.ID});
                    return newState
                }
            }
            throw new Error("Member to be deleted not found");
        case actionTypes.SET_MESSAGES:
            newState = {...state};
            for (let i = 0; i < newState.groups.length; i++) {
                if (newState.groups[i].ID === action.payload.group) {
                    newState.groups[i].messages = action.payload.messages
                    return newState;
                }
            }
            throw new Error("Received messages don't belong to any of your groups");
        case actionTypes.ADD_MESSAGE:
            newState = {...state};
            for (let i = 0; i < newState.groups.length; i++) {
                if (newState.groups[i].ID === action.payload.group) {
                    newState.groups[i].messages = [...newState.groups[i].messages, action.payload];
                    return newState;
                }
            }
            throw new Error("Received message don't belong to any of your groups");
        case actionTypes.ADD_MESSAGES:
            newState = {...state};
            for (let i = 0; i < newState.groups.length; i++) {
                if (newState.groups[i].ID === action.payload.group) {
                    newState.groups[i].messages = [...action.payload.messages, ...newState.groups[i].messages];
                    return newState;
                }
            }
            throw new Error("Received messages don't belong to any of your groups");
        default:
            throw new Error("Action not specified");
    }
}

const ChatStorage = ({children}) => {

    const [state, dispatch] = useReducer(reducer, {groups: []});

    return (
        <StorageContext.Provider value={[state, dispatch]}>
            {children}
        </StorageContext.Provider>
    );
}

export default ChatStorage;
