import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import "../styles/UpdateNote.css";

const UpdateNote = () => {
  const { id } = useParams();
  const [content, setContent] = useState("");
  const [ws, setWs] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const websocket = new WebSocket("ws://localhost:8082/ws");

    websocket.onopen = () => {
      console.log("Connected to WebSocket");
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === "init") {
        const note = message.payload.find((note) => note.id === id);
        if (note) {
          setContent(note.content);
        }
      }
    };

    websocket.onclose = () => {
      console.log("Disconnected from WebSocket");
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [id]);

  const handleUpdate = async (e) => {
    e.preventDefault();
    const updatedNote = { id, content, version: 0 }; // Fetch the current version from the server

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          type: "update",
          payload: updatedNote,
        })
      );

      alert("Note updated successfully!");
      navigate("/notes");
    } else {
      console.error("WebSocket is not connected");
      alert("Failed to update note. Please try again.");
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <div className="updatenote-container">
      <nav className="navbar">
        <div className="navbar-brand">NoteApp</div>
        <div className="navbar-links">
          <Link to="/profile" className="nav-link">
            User Profile
          </Link>
          <Link to="/createnote" className="nav-link">
            Create Note
          </Link>
          <Link to="/notes" className="nav-link">
            Note List
          </Link>
          <button onClick={handleLogout} className="nav-link logout-btn">
            Logout
          </button>
        </div>
      </nav>

      <div className="updatenote-body">
        <div className="card">
          <h2>Update Note</h2>
          <form onSubmit={handleUpdate}>
            <div className="form-group">
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Update your note content"
                rows="10"
              ></textarea>
            </div>
            <button type="submit">Update Note</button>
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

export default UpdateNote;
