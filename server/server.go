// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"html/template"
	"bytes"
)



type Peer struct{
	ID string
	Addr string
}

type Event struct {
	ID string
	Peers map[string]*Peer
}
var(
	Peers map[string]*Peer
	Events map[string]*Event
	//frontPageTmpl = template.Must(template.ParseFiles("/Users/dong/go/src/github.com/lishaodong/MazeServer/server/server.html"))
	frontPageTmpl = template.Must(template.ParseFiles("server.html"))
)


func init() {
	Peers = make(map[string]*Peer)
	Events = make(map[string]*Event)

	http.HandleFunc("/",frontPageHandler)
	http.HandleFunc("/peer", handlePeer)
	http.HandleFunc("/event",handleEvent)
}

func main(){
	http.ListenAndServe(":8081",nil)
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	b := new(bytes.Buffer)
	if err := frontPageTmpl.Execute(b, nil); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "tmpl.Execute failed: %v", err)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	b.WriteTo(w)
}

func handlePeer(w http.ResponseWriter, r *http.Request) {
	addr:=r.FormValue("addr")

	id:=r.FormValue("id")
	if r.Method=="GET"{
		fmt.Fprint(w,findPeer(id))
		return
	}
	if r.Method == "POST"{
		announcePeer(id, addr)
		fmt.Fprint(w,"peer announced: "+id+"adddr:"+addr)
	}

}
func findPeer(ID string) (string){
	peer ,ok:= Peers[ID]
	if ok{
		b,err:=json.Marshal(peer)
		if err!=nil{
			fmt.Println("can't marshal peer")
			return "null"
		}
		return string(b)
	}else{

	fmt.Println("don't find:"+ID)
	return "null"
	}
}

func announcePeer(id string, addr string){

	Peers[id]=&Peer{ID:id,Addr:addr}
	fmt.Println("announced:"+Peers[id].Addr)
}


//Events
func handleEvent(w http.ResponseWriter, r *http.Request){
	eid:=r.FormValue("eid")
	pid:=r.FormValue("pid")
	addr:=r.FormValue("addr")
	if r.Method=="GET"{
		fmt.Fprint(w,findEvent(eid))
		return
	}
	if r.Method=="POST"{
		announceEvent(eid,pid,addr)
		fmt.Fprint(w,"event announced: "+eid)
	}
}

func findEvent(eid string) string{
	event,ok:=Events[eid]
	if ok{
		b,err:=json.Marshal(event)
		if err!=nil{
			fmt.Println("can't marshal peer")
			return "null"
		}
		return string(b)
	}else{
		fmt.Println("don't find:"+eid)
		return "null"
	}
}

func announceEvent(eid,pid,addr string){
	event ,ok:= Events[eid]
	if ok{
		event.Peers[pid]=&Peer{ID:pid,Addr:addr}
		return
	}
	Events[eid]=&Event{ID:eid,Peers:make(map[string]*Peer)}
	event.Peers[pid]=&Peer{ID:pid,Addr:addr}
	return
}

func output(){
	for _,peer:=range Peers{
		fmt.Println(peer.ID)
	}
}
