import React, {createContext, useReducer} from "react";

export const StorageContext = createContext({});

export const actionTypes = {
    LOGIN: "LOGIN",
    SET_GROUPS: "SET_GROUPS",
    NEW_GROUP: "NEW_GROUP",
    DELETE_GROUP: "DELETE_GROUP",
    ADD_MEMBER: "ADD_MEMBER",
    DELETE_MEMBER: "DELETE_MEMBER"
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
            break;
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


// const handleGroupDelete = (id) => {
//     let newGroups = groups.filter((item)=>{return item.ID !== id})
//     setGroups(newGroups);
// };

// const handleMemberDelete = (member) => {
//     let newGroups = groups;
//     for (let i = 0; i < newGroups.length; i++) {
//         if (newGroups[i].ID === member.group_id) {
//             newGroups[i].Members = newGroups[i].Members.filter((item)=>{return item.ID !== member.ID});
//             setGroups(newGroups);
//         }
//     }
// }

// const handleMemberAdd = (member) => {
//     let newGroups = groups;
//     for (let i = 0; i < groups.length; i++) {
//         if (newGroups[i].ID === member.group_id) {
//             newGroups[i].Members.push(member);
//             console.log(newGroups);
//             setGroups(newGroups);
//         }
//     }
// }