import React from "react";
import {Navigate} from "react-router-dom";

export default class AuthMain extends React.Component {
    constructor(props){
        super(props)
        this.state = {
            user: {
                loggedin: false,
                name: "",
            }
        }
    }

    render(){
        return (
            <div>
                {!this.state.user.loggedin? <Navigate to="/login" />:<Main/>}
            </div>
        );
    }
}

class Main extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            aa: "as"
        };
    }

    render() {
        return (
            <h1>Main page</h1>
        )
    }
}