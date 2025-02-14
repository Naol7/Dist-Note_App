import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import "../styles/CreateNote.css";

const CreateNote = () => {
  const [content, setContent] = useState("");
  const [ws, setWs] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const websocket = new WebSocket("ws://localhost:8082/ws");

    websocket.onopen = () => {
      console.log("Connected to WebSocket");
    };

    websocket.onclose = () => {
      console.log("Disconnected from WebSocket");
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const newNote = { id: Date.now().toString(), content, version: 0 };

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          type: "create",
          payload: newNote,
        })
      );

      alert("Note created successfully!");
      setContent("");
      navigate("/notes");
    } else {
      console.error("WebSocket is not connected");
      alert("Failed to create note. Please try again.");
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <div className="createnote-container">
      <nav className="navbar">
        <div className="navbar-brand">NoteApp</div>
        <div className="navbar-links">
          <Link to="/notes" className="nav-link">
            Note List
          </Link>
          <Link to="/profile" className="nav-link">
            User Profile
          </Link>
          <button onClick={handleLogout} className="nav-link logout-btn">
            Logout
          </button>
        </div>
      </nav>

      <div className="createnote-body">
        <div className="card">
          <h2>Create a New Note</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Enter note content"
                rows="5"
              ></textarea>
            </div>
            <button type="submit">Create Note</button>
          </form>
        </div>
      </div>

      <footer className="home-footer">
        <div className="footer-social">
          <a
            href="https://twitter.com"
            target="_blank"
            rel="noopener noreferrer"
          >
            Twitter
          </a>
          <a
            href="https://facebook.com"
            target="_blank"
            rel="noopener noreferrer"
          >
            Facebook
          </a>
          <a
            href="https://linkedin.com"
            target="_blank"
            rel="noopener noreferrer"
          >
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

export default CreateNote;
