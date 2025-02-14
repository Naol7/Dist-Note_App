import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import "../styles/UserProfile.css";

const UserProfile = () => {
  const [username, setUsername] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      setMessage("You must be logged in to view this page.");
      return;
    }

    fetchUserProfile(token);
  }, []);

  const fetchUserProfile = async (token) => {
    try {
      const response = await fetch("http://localhost:8080/protected", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to load user profile");
      }

      const data = await response.json();
      setUsername(data.username);
    } catch (error) {
      console.error("Error fetching user profile:", error);
      setMessage("Failed to load user profile");
    }
  };

  const handlePasswordChange = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("token");

    try {
      const response = await fetch("http://localhost:8080/update", {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ password: newPassword }),
      });

      const data = await response.json();
      if (response.ok) {
        setMessage("Password updated successfully!");
        setNewPassword("");
      } else {
        setMessage(data.message || "Failed to update password");
      }
    } catch (error) {
      console.error("Error updating password:", error);
      setMessage("Error updating password");
    }
  };

  const handleDelete = async () => {
    const token = localStorage.getItem("token");
    try {
      const response = await fetch("http://localhost:8080/delete", {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok) {
        localStorage.removeItem("token");
        navigate("/");
      } else {
        setMessage("Failed to delete account");
      }
    } catch (error) {
      console.error("Error deleting account:", error);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <div className="userprofile-container">
      <nav className="navbar">
        <div className="navbar-brand">NoteApp</div>
        <div className="navbar-links">
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

      <div className="userprofile-body">
        <div className="card">
          <h2>User Profile</h2>
          {message && (
            <div
              className={`alert ${
                message.includes("success") ? "alert-success" : "alert-error"
              }`}
            >
              {message}
            </div>
          )}

          <div className="profile-info">
            <p>
              Username: <strong>{username}</strong>
            </p>
          </div>

          <form onSubmit={handlePasswordChange}>
            <div className="form-group">
              <label>New Password:</label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
                minLength="6"
              />
            </div>
            <button type="submit" className="update-btn">
              Change Password
            </button>
          </form>

          <div className="danger-zone">
            <h3>Danger Zone</h3>
            <button onClick={handleDelete} className="delete-btn">
              Delete Account
            </button>
          </div>
        </div>
      </div>

      <footer className="home-footer">{/* Footer remains the same */}</footer>
    </div>
  );
};

export default UserProfile;
