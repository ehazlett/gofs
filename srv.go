package gofs

import (
	"fmt"
	"io"

	"github.com/mortdeus/go9p"
	go9ps "github.com/mortdeus/go9p/srv"
)

type Server struct {
	h Handler
}

type FileEntry struct {
	h      Handler
	openro io.ReadSeeker
	openwo io.WriteSeeker
	openrw io.ReadWriteSeeker
}

func (n *Server) Attach(req *go9ps.Req) {
	fmt.Printf("Attach fid=%d\n", req.Tc.Fid)
	qids := newQids(req.Tc.Mode)
	req.RespondRattach(&qids[0])
}

func (n *Server) Flush(req *go9ps.Req) {
	fmt.Printf("Flush fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	req.RespondError(go9ps.Enoent)
}

func (n *Server) Open(req *go9ps.Req) {
	fmt.Printf("Open fid=%d wname=%v\n", req.Tc.Fid, req.Tc.Wname)
	if req.Fid.Aux == nil {
		req.RespondError(go9ps.Enoent)
		return
	}
	e, ok := req.Fid.Aux.(*FileEntry)
	if !ok {
		req.RespondError(go9ps.Enoent)
		return
	}
	if e.h == nil {
		req.RespondError(go9ps.Enoent)
		return
	}
	if e.h.IsDir() {
		// FIXME: wrap dir operations in a readwriter to the 9p dir entry serialized format.
	} else {
		var err error
		e.openrw, err = e.h.OpenRW()
		if err != nil {
			req.RespondError(err)
			return
		}
	}
	qids := newQids(req.Tc.Mode)
	req.RespondRopen(&qids[0], 0)
}

func (n *Server) Create(req *go9ps.Req) {
	fmt.Printf("Create fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	req.RespondError(go9ps.Enoent)
}

func (n *Server) Remove(req *go9ps.Req) {
	fmt.Printf("Remove fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	req.RespondError(go9ps.Enoent)
}

/*
type Dir struct {
    Size   uint16 // size-2 of the Dir on the wire
    Type   uint16
    Dev    uint32
    Qid           // file's Qid
    Mode   uint32 // permissions and flags
    Atime  uint32 // last access time in seconds
    Mtime  uint32 // last modified time in seconds
    Length uint64 // file length in bytes
    Name   string // file name
    Uid    string // owner name
    Gid    string // group name
    Muid   string // name of the last user that modified the file

    Ext     string // special file's descriptor
    Uidnum  uint32 // owner ID
    Gidnum  uint32 // group ID
    Muidnum uint32 // ID of the last user that modified the file
}
    Dir describes a file
*/

func (n *Server) Stat(req *go9ps.Req) {
	fmt.Printf("Stat fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
}

func (n *Server) Wstat(req *go9ps.Req) {
	fmt.Printf("Wstat fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	req.RespondError(go9ps.Enoent)
}

func (n *Server) Clunk(req *go9ps.Req) {
	fmt.Printf("Clunk fid=%d\n", req.Tc.Fid)
	req.RespondRclunk()
}

func (n *Server) Read(req *go9ps.Req) {
	fmt.Printf("Read fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	req.RespondError(go9ps.Enoent)
}

func (n *Server) Write(req *go9ps.Req) {
	fmt.Printf("Write fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	req.RespondError(go9ps.Enoent)
}

func (n *Server) Walk(req *go9ps.Req) {
	fmt.Printf("Walk fid=%d wname=%v newfid=%d\n", req.Tc.Fid, req.Tc.Wname, req.Tc.Newfid)
	e, ok := req.Fid.Aux.(*FileEntry)
	if !ok {
		req.RespondError(go9ps.Enoent)
		return
	}
	if e.h == nil {
		req.RespondError(go9ps.Enoent)
		return
	}
	newH, err := e.h.WalkDir(req.Tc.Wname...)
	if err != nil {
		req.RespondError(err)
		return
	}
	req.Newfid.Aux = &FileEntry{
		h: newH,
	}
	req.RespondRwalk(newQids(req.Tc.Mode))
}

var qnext uint64

func newQids(modes ...uint8) []go9p.Qid {
	qids := make([]go9p.Qid, 0, len(modes))
	for _, mode := range modes {
		qids = append(qids, go9p.Qid{
			Type:    mode,
			Version: 0,
			Path:    qnext,
		})
		fmt.Printf("[newqid] allocated qid=%d\n", qnext)
		qnext++
	}
	return qids
}
