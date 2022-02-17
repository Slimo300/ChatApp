import React from "react";

export default class SignInForm extends React.Component {
    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.state = {
            errormsg: ""
        };
    }
    handleChange(event) {
        const name = event.target.name;
        const value = event.target.value;
        this.setState({
            [name]: value
        });
    }

    handleSubmit(event) {
        event.preventDefault();
        console.log(JSON.stringify(this.state));
    }

    render() {
        let message = null;
        if (this.state.errormsg.length !== 0) {
            message = <h5 className="mb-4 text-danger">{this.state.errormsg}</h5>;
        }
        return (
            <div className="mt-5 d-flex justify-content-center">
              <div className="mt-5 row">
                {message}
                <form onSubmit={this.handleSubmit}>
                  <div className="display-3 mb-4 text-center text-primary"> Log In</div>
                  <div className="mb-3 text-center">
                    <label htmlFor="email" className="form-label">Email address</label>
                    <input name="email" type="email" className="form-control" id="email" onChange={this.handleChange}/>
                  </div>
                  <div className="mb-3 text-center">
                    <label htmlFor="password" className="form-label">Password</label>
                    <input name="password" type="password" className="form-control" id="password" onChange={this.handleChange}/>
                  </div>
                  <div className="text-center">
                    <button type="submit" className="btn btn-primary text-center">Submit</button>
                  </div>
                  <div className="display-5 mt-4 text-center text-primary"><a href="/register">or Register</a></div>
                </form>
              </div>
            </div>
        )
    }
}