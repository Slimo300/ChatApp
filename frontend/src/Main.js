import React from "react";
import {Navigate} from "react-router-dom";

export default function AuthMain(){

    console.log(this.props.name)
    return (
        <div>
            {this.props.name === ""? <Navigate to="/login" />:<Main/>}
        </div>
    );
}

class Main extends React.Component {

    render() {
        return (
            <h1>Main page</h1>
        )
    }
}