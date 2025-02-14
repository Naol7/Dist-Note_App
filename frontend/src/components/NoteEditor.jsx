import React, { useEffect, useState, useCallback } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import "../styles/NoteEditor.css";

const NoteEditor = () => {
  const { noteId } = useParams();
  const [content, setContent] = useState("");
  const [ws, setWs] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const websocket = new WebSocket("ws://localhost:8082/ws");

    websocket.onopen = () => {
      console.log("Connected to WebSocket");
      // Join the specific note room
      websocket.send(
        JSON.stringify({
          type: "join_note",
          noteId: noteId,
        })
      );
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === "delta" && message.payload.id === noteId) {
        setContent((prevContent) => applyDelta(prevContent, message.payload));
      }
    };

    websocket.onclose = () => {
      console.log("Disconnected from WebSocket");
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [noteId]);

  const applyDelta = (content, delta) => {
    return (
      content.slice(0, delta.start) + delta.text + content.slice(delta.end)
    );
  };

  const handleChange = useCallback(
    (e) => {
      const newContent = e.target.value;
      setContent(newContent);

      if (ws && ws.readyState === WebSocket.OPEN) {
        const delta = {
          id: noteId,
          start: 0,
          end: content.length,
          text: newContent,
          version: 0, // You can fetch the current version from the server
        };
        ws.send(
          JSON.stringify({
            type: "delta",
            payload: delta,
          })
        );
      }
    },
    [ws, noteId, content]
  );

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <div className="noteeditor-container">
      <nav className="navbar">
        <div className="navbar-brand">NoteApp</div>
        <div className="navbar-links">
          <Link to="/profile" className="nav-link">
            User Profile
          </Link>
          <Link to="/createnote" className="nav-link">
            Create Note
          </Link>
          <button onClick={handleLogout} className="nav-link logout-btn">
            Logout
          </button>
        </div>
      </nav>

      <div className="noteeditor-body">
        <div className="card">
          <h2>Edit Note</h2>
          <textarea
            value={content}
            onChange={handleChange}
            placeholder="Edit your note here..."
            rows="10"
          />
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

export default NoteEditor;
