import { useState } from "react";
import axios from "axios";
import { Link, useNavigate } from "react-router-dom";
import "../styles/Signup.css";

const Signup = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const navigate = useNavigate();
  const register = async (e) => {
    e.preventDefault();
    try {
      await axios.post("http://localhost:8080/register", {
        username,
        password,
      });
      alert("User registered successfully!");
      navigate("/");
    } catch (error) {
      alert("Registration failed. Please try again.");
      console.error("Registration error:", error);
    }
  };

  return (
    <div className="signup-container">
      {/* Navbar */}
      <nav className="navbar">
        <div className="navbar-brand">NoteApp</div>
        <div className="navbar-links">
          <Link to="/" className="nav-link">
            Home
          </Link>
          <Link to="/login" className="nav-link">
            Login
          </Link>
        </div>
      </nav>

      {/* Body Section */}
      <div className="signup-body">
        <div className="card">
          <h2>Signup</h2>
          <form onSubmit={register}>
            <div className="form-group">
              <input
                type="text"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>
            <div className="form-group">
              <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <button type="submit">Register</button>
          </form>
          <p>
            Already have an account? <Link to="/login">Login here</Link>
          </p>
        </div>
      </div>

      {/* Footer */}
      <footer className="home-footer">
        <div className="footer-social">
          <a href="https://twitter.com" target="_blank" rel="noopener noreferrer">
            Twitter
          </a>
          <a href="https://facebook.com" target="_blank" rel="noopener noreferrer">
            Facebook
          </a>
          <a href="https://linkedin.com" target="_blank" rel="noopener noreferrer">
            LinkedIn
          </a>
        </div>
        <div className="footer-help">
          <a href="mailto:help@noteapp.com">help@noteapp.com</a>
        </div>
        <div className="footer-copyright">
          &copy; {new Date().getFullYear()} NoteApp. All rights reserved.
        </div>
      </footer>
    </div>
  );
};

export default Signup;