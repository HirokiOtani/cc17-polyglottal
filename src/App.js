import "./App.css";
import { Complaints } from "./Complaints";
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Redirect,
  Link,
} from "react-router-dom";
import { InputComplaint } from "./InputComplaint";

function App() {
  return (
    <div className="App">
      <Router>
        <body>
          <Link to="/">Home</Link>
          <br />
          <Link to="/patient">Patient</Link>
          <br />
          <Link to="/doctor">Doctor</Link>
          <Route exact path="/">
            <h1>Summarize Complaints</h1>
          </Route>
          <Route exact path="/patient" component={InputComplaint} />
          <Route exact path="/doctor" component={Complaints} />
        </body>
      </Router>
    </div>
  );
}
export default App;
