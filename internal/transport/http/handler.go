package http

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"

  "github.com/bgercken/go-rest-api-course/internal/comment"
  "github.com/gorilla/mux"
)

// Handler - stores the pointer to our comments service
type Handler struct {
  Router *mux.Router
  Service *comment.Service
}

// Response - an object to store responses from our API
type Response struct {
  Message string
  Error string
}

// NewHandler - returns a pointer to a handler
func NewHandler(service *comment.Service) *Handler {
  return &Handler{
    Service: service,
  }
}

// SetupRoutes - setups all the routes for our application
func (h *Handler)SetupRoutes() {
  fmt.Println("Setting Up Routes")
  h.Router = mux.NewRouter()

  h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
  h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
  h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
  h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
  h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

  h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
    if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive"}); err != nil {
      sendErrorResponse(w, "Error encoding message.", err)
      return
    }
  })
}

// GetComment - retrieve a comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
  fmt.Println("GetComment")
  vars := mux.Vars(r)
  id := vars["id"]

  i, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    sendErrorResponse(w, "Unable tp parse UINT from ID", err)
    return
  }
  comment, err := h.Service.GetComment(uint(i))
  if err != nil {
    sendErrorResponse(w, "Error Retrieving Comment By ID", err)
    return
  }

  writeStatusOK(w)
  writeCommentResponse(w, comment)  // ? what happens to the object comment? -> comment.Comment
}

// GetAllComments - retrieves all comments from the comment service
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
  fmt.Println("GetAllComments")
  comments, err := h.Service.GetAllComments()
  if err != nil {
    sendErrorResponse(w, "Failed to retrieve all comments", err)
    return
  }

  writeStatusOK(w)
  writeCommentsResponse(w, comments) 
}

// PostComent - adds a new comment
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
  fmt.Println("PostComment")
  var comment comment.Comment
  if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
    sendErrorResponse(w, "Failed to decode JSON Body", err)
    return
  }

  comment, err := h.Service.PostComment(comment)
  if err != nil {
    sendErrorResponse(w, "Failed to post new comment", err)
    return
  }

  writeStatusOK(w)
  writeCommentResponse(w, comment)
}

// UpdateComment - updates a comment by ID
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
  fmt.Println("UpdateComment")

  vars := mux.Vars(r)
  id := vars["id"]
  commentID, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    sendErrorResponse(w, "Failed to parse uint from ID", err)
    return
  }

  var comment comment.Comment
  if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
    sendErrorResponse(w, "Failed to decode JSON Body", err)
    return
  }

  // Update the comment in the database.
  comment, err = h.Service.UpdateComment(uint(commentID), comment)
  if err != nil {
    sendErrorResponse(w, "Failed to update comment", err)
    return
  }

  // Encode writes the JSON encoding of v to the stream, followed by a newline character.
  writeStatusOK(w)
  if err := json.NewEncoder(w).Encode(comment); err != nil {
    panic(err)
    return
  }
}

// DeleteComment - deletes a comment by ID
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
  fmt.Println("DeleteComment")
  vars := mux.Vars(r)
  id := vars["id"]
  commentID, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    sendErrorResponse(w, "Failed to parse UINT from ID", err)
  }

  err = h.Service.DeleteComment(uint (commentID))
  if err != nil {
    sendErrorResponse(w, "Failed to delete comment by comment ID", err)
    return
  }
  
  writeStatusOK(w)
  writeStringResponse(w, "Comment successfully deleted")
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
  w.WriteHeader(http.StatusInternalServerError)
  if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
    panic(err)
  }
}

// writeStatusOK - helper function to send http.StatusOK in JSON format
func writeStatusOK(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
}

// writeStringResponse - helper function to write results in JSON format
func writeStringResponse(w http.ResponseWriter, str string) {
  if err := json.NewEncoder(w).Encode(str); err != nil {
    panic(err)
  }
}

// writeCommentResponse - helper function to write comment in JSON format
func writeCommentResponse(w http.ResponseWriter, c comment.Comment) {
  if err := json.NewEncoder(w).Encode(c); err != nil {
    panic(err)
  }
}

// writeCommentsResponse - helper function to write comment(s) in JSON format
func writeCommentsResponse(w http.ResponseWriter, c []comment.Comment) {
  if err := json.NewEncoder(w).Encode(c); err != nil {
    panic(err)
  }
}
