import React from "react";
import { Link } from "react-router-dom";
import "../styles/Home.css"; // Import the CSS file for styling

const Home = () => {
  return (
    <div className="home-container">
      {/* Navbar */}
      <nav className="navbar">
        <div className="navbar-brand">NoteApp</div>
        <div className="navbar-links">
          <Link to="/login" className="nav-link">
            Login
          </Link>
          <Link to="/signup" className="nav-link">
            Signup
          </Link>
        </div>
      </nav>

      {/* Body Section */}
      <div className="home-body">
        <h1 className="home-title">Welcome to NoteApp</h1>
        <p className="home-intro">
          NoteApp is a simple and intuitive platform for creating, managing, and
          organizing your notes in real-time. Whether you're jotting down ideas,
          making to-do lists, or collaborating with others, NoteApp has you
          covered.
        </p>
        <div className="home-actions">
          <Link to="/signup" className="home-button">
            Get Started
          </Link>
          <Link to="/login" className="home-button secondary">
            Login
          </Link>
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

export default Home;