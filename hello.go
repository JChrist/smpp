package main

import (
	"fmt"
	fxsmpp "github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/jchrist/smppserver/smpp"
)

func main() {
	srv := smpp.NewServer()
	srv.Handler = func(c smpp.Conn, m pdu.Body) {
		fmt.Println("received pdu", m, "decoded body:", m.Fields()[pdufield.ShortMessage])
		// Real servers will reply with at least the same sequence number from the request:
		resp := pdu.NewSubmitSMResp()
		resp.Header().Seq = m.Header().Seq
		resp.Fields().Set(pdufield.MessageID, "1234")
		resp.Header().Status = 0x0
		c.Write(resp)
	}
	defer srv.Close()

	fmt.Println("started server at", srv.Addr())

	tx := &fxsmpp.Transceiver{
		Addr:   srv.Addr(),
		User:   "client",
		Passwd: "secret",
	}

	conStatus := <-tx.Bind()

	if conStatus.Error() != nil {
		fmt.Println("received error connecting to server", conStatus.Error())
		return
	}

	fmt.Println("conStatus reported successful bind:", conStatus.Status())

	if sm, err := tx.Submit(&fxsmpp.ShortMessage{Text: pdutext.GSM7("yo yo yo")}); err != nil {
		fmt.Println("error submitting short message", err)
	} else {
		fmt.Println("submitted short message", sm)
	}

	if err := tx.Close(); err != nil {
		fmt.Println("error closing connection", err)
	}

	fmt.Println("finished!")
}
