import React from "react";

export default class RegisterForm extends React.Component {
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
            <div className="mt-5">
              <div className="mt-5 row">
                <form>
                  <div className="display-1 mb-4 text-center text-primary"> Register</div>
                  <div class="mb-3 text-center">
                    <label for="username" className="form-label">Username</label>
                    <input type="text" className="form-control" id="username" onChange={this.handleChange}/>
                  </div>
                  <div class="mb-3 text-center">
                    <label for="email" className="form-label">Email address</label>
                    <input type="email" className="form-control" id="email" onChange={this.handleChange}/>
                  </div>
                  <div class="mb-3 text-center">
                    <label for="password" class="form-label">Password</label>
                    <input type="password" class="form-control" id="password" onChange={this.handleChange}/>
                  </div>
                  <div class="mb-3 text-center">
                    <label for="r-password" class="form-label">Repeat Password</label>
                    <input type="password" class="form-control" id="r-password" onChange={this.handleChange}/>
                  </div>
                  <div class="text-center">
                    <button type="submit" class="btn btn-primary text-center">Submit</button>
                  </div>
                  <div className="display-5 mt-4 text-center text-primary"><a href="#">or Log in</a></div>
                </form>
              </div>
            </div>
          )
    }
}