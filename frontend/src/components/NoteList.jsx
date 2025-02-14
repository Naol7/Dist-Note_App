import { useEffect, useState, useCallback } from "react";
import { Link, useNavigate } from "react-router-dom";
import "../styles/NoteList.css";

const NoteList = () => {
  const [notes, setNotes] = useState([]);
  const [ws, setWs] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const websocket = new WebSocket("ws://localhost:8082/ws");

    websocket.onopen = () => {
      console.log("Connected to WebSocket");
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data);

      switch (message.type) {
        case "init":
          setNotes(message.payload);
          break;
        case "create":
          setNotes((prev) => [...prev, message.payload]);
          break;
        case "update":
          setNotes((prev) =>
            prev.map((note) =>
              note.id === message.payload.id ? message.payload : note
            )
          );
          break;
        case "delete":
          setNotes((prev) =>
            prev.filter((note) => note.id !== message.payload.id)
          );
          break;
        default:
          console.log("Unknown message type:", message.type);
      }
    };

    websocket.onclose = () => {
      console.log("Disconnected from WebSocket");
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, []);

  const handleDelete = useCallback(
    (id) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(
          JSON.stringify({
            type: "delete",
            payload: { id },
          })
        );
      } else {
        console.error("WebSocket is not connected");
      }
    },
    [ws]
  );

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <div className="notelist-container">
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

      <div className="notelist-body">
        <h2>List of Notes</h2>
        <div className="note-list">
          {notes.map((note) => (
            <div key={note.id} className="note-item">
              <p>{note.content}</p>
              <div className="note-actions">
                <Link to={`/update/${note.id}`} className="update-btn">
                  Update
                </Link>
                <button
                  className="delete-btn"
                  onClick={() => handleDelete(note.id)}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
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

export default NoteList;
