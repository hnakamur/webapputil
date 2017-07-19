// Package problem provides a base struct for problem detail specified in RFC 7807,
// and a function for sending problem json response.
//
// Here is a server code example.
//
//	package main
//
//	import (
//		"encoding/json"
//		"flag"
//		"log"
//		"net/http"
//
//		"github.com/hnakamur/webapputil/problem"
//	)
//
//	func main() {
//		addr := flag.String("addr", ":8080", "server listen address")
//		flag.Parse()
//
//		http.HandleFunc("/", handleRoot)
//		s := &http.Server{
//			Addr: *addr,
//		}
//		log.Fatal(s.ListenAndServe())
//	}
//
//	type problemExample1 struct {
//		problem.Problem
//		Balance  int      `json:"balance"`
//		Accounts []string `json:"accounts"`
//	}
//
//	type problemExample2 struct {
//		problem.Problem
//		InvalidParams []invalidParam `json:"invalid-params"`
//	}
//
//	type invalidParam struct {
//		Name   string `json:"name"`
//		Reason string `json:"reason"`
//	}
//
//	func sendResponse(w http.ResponseWriter, statusCode int, detail interface{}) error {
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(statusCode)
//		enc := json.NewEncoder(w)
//		return enc.Encode(detail)
//	}
//
//	func handleRoot(w http.ResponseWriter, r *http.Request) {
//		err := r.ParseForm()
//		if err != nil {
//			log.Printf("failed to parse form; %v", err)
//			return
//		}
//		p := r.Form.Get("prob")
//		switch p {
//		case "ex1":
//			prob := problemExample1{
//				Problem: problem.Problem{
//					Type:     "https://example.com/probs/out-of-credit",
//					Title:    "You do not have enough credit.",
//					Detail:   "Your current balance is 30, but that costs 50.",
//					Instance: "/account/12345/msgs/abc",
//				},
//				Balance: 30,
//				Accounts: []string{
//					"/account/12345",
//					"/account/67890",
//				},
//			}
//			err = problem.SendProblem(w, http.StatusForbidden, prob)
//			if err != nil {
//				log.Printf("failed to send problem; %v", err)
//				return
//			}
//		case "ex2":
//			prob := problemExample2{
//				Problem: problem.Problem{
//					Type:  "https://example.net/validation-error",
//					Title: "Your request parameters didn't validate.",
//				},
//				InvalidParams: []invalidParam{
//					{Name: "age", Reason: "must be a positive integer"},
//					{Name: "color", Reason: "must be 'green', 'red' or 'blue'"},
//				},
//			}
//			err = problem.SendProblem(w, http.StatusBadRequest, prob)
//			if err != nil {
//				log.Printf("failed to send problem; %v", err)
//				return
//			}
//		default:
//			err = sendResponse(w, http.StatusOK, struct {
//				Msg string `json:"msg"`
//			}{
//				Msg: "Hello, world!",
//			})
//			if err != nil {
//				log.Printf("failed to send response; %v", err)
//				return
//			}
//		}
//	}
package problem

import (
	"encoding/json"
	"net/http"
)

const contentType = "application/problem+json"

// Problem is the base type of a problem detail specified in RFC 7807.
// You can embed Problem in your application specific problem struct.
type Problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// SendProblem writes a problem json response
// with Content-Type "application/problem+json".
func SendProblem(w http.ResponseWriter, statusCode int, problem interface{}) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	return enc.Encode(problem)
}
