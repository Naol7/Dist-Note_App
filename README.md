

---

# **Distributed Note App**

--



1. [Features](#features)
2. [Technologies Used](#technologies-used)
3. [Setup Instructions](#setup-instructions)
4. [API Endpoints](#api-endpoints)
5. [Frontend Pages](#frontend-pages)

---

---

## **Features**

**User Features:**

- **Real-Time Collaboration:** Instant live updates for note editing.
- **Note Management:** Create, update, delete, and list notes effortlessly.
- **User Authentication:** Secure signup, login, and profile management using JWT.
- **Delta Updates:** Optimized partial note updates for reduced bandwidth usage.

---

## **Technologies Used**

- **Backend:** Go, Gorilla WebSocket, MongoDB
- **Frontend:** React, Vite, Axios
- **Authentication:** JWT-based security

---

## **Setup Instructions**

### **Prerequisites**

- Go (v1.16 or later)
- Node.js and npm
- MongoDB (or MongoDB Atlas)
- Git

### **Steps**

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/your-username/distributed-note-app.git
   cd distributed-note-app
   ```

2. **Backend Setup:**

   - **Authentication/Note API:**
     ```bash
     cd auth
     go mod tidy
     go run main.go
     ```
   - **WebSocket Server:**
     ```bash
     cd ../websocket
     go mod tidy
     go run main.go
     ```

3. **Frontend Setup:**

   ```bash
   cd ../frontend
   npm install
   npm run dev
   ```

   The application will be available at [http://localhost:5173](http://localhost:5173).

---

## **API Endpoints**

**User Authentication:**

- `POST /api/users/signup` – Register a new user.
- `POST /api/users/login` – Authenticate an existing user.
- `PUT /api/users/profile` – Update user profile information.

**Note Management:**

- `GET /api/notes` – Retrieve all notes.
- `GET /api/notes/{id}` – Retrieve a specific note.
- `POST /api/notes` – Create a new note.
- `PUT /api/notes/{id}` – Update an existing note.
- `DELETE /api/notes/{id}` – Delete a note.

**WebSocket:**

- `/ws` – Establish a WebSocket connection for real-time note collaboration.

---

## **Frontend Pages**

**Public Pages:**

- **Login:** `/login`
- **Signup:** `/signup`
- **Home:** `/`

**Note Pages:**

- **Create Note:** `/create`
- **Note Editor:** `/edit/{noteId}`
- **Note List:** `/notes`
- **Update Note:** `/update/{noteId}`
- **User Profile:** `/profile`

---

---
