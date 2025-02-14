import "./styles/global.css";
import Home from "./components/Home";
import NoteList from "./components/NoteList";
import CreateNote from "./components/CreateNote";
import Login from "./components/Login";
import Signup from "./components/Signup";
import UpdateNote from "./components/UpdateNote";
import UserProfile from "./components/UserProfile";
import { Route, Routes } from "react-router-dom";
import NoteEditor from "./components/NoteEditor";

function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/notes" element={<NoteList />} />
        <Route path="/createnote" element={<CreateNote />} />
        <Route path="/update/:id" element={<UpdateNote />} />
        <Route path="/profile" element={<UserProfile />} />

        <Route path="/edit" element={<NoteEditor />} />
      </Routes>
    </div>
  );
}

export default App;
