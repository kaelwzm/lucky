package main

import (
	"github.com/kataras/iris/v12/context"
	"lucky/core/ihttp"
	"lucky/core/iproto"
	"lucky/example/comm/msg"
	"lucky/example/comm/msg/code"
	"lucky/example/comm/protobuf"
	"lucky/log"
)

func main() {
	httpProcessor := iproto.NewPBProcessor()
	msg.SetEncrypt(httpProcessor)
	httpProcessor.RegisterHandler(code.Hello, &protobuf.Hello{}, func(args ...interface{}) {
		hello := args[iproto.Msg].(*protobuf.Hello)
		log.Debug(hello.Hello)
		ctx := args[iproto.Conn].(*context.Context)
		data, err := httpProcessor.WarpMsgNoHeader(hello)
		if err != nil {
			panic(err)
		}
		_, err = ctx.Write(data)
		if err != nil {
			panic(err)
		}
	})
	ihttp.EnableCrossOrigin()
	ihttp.Post("/", func(context *context.Context) {
		body, err := context.GetBody()
		if err != nil {
			log.Error("Read body error %v", err)
			return
		}
		httpProcessor.OnReceivedPackage(context, body)
	})
	log.Error("http run error %v", ihttp.Run(":3001"))
}
