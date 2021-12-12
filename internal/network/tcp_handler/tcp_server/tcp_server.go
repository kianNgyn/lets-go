package tcp_server

import (
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/proto"
	"github.com/nkien0204/projectTemplate/configs"
	
	"encoding/binary"
	"github.com/nkien0204/protobuf/build/proto/events"
	"io"
	"bufio"
	"go.uber.org/zap"
	"github.com/nkien0204/projectTemplate/internal/log"
	"net"
)

func NewTcpServer(cfg *configs.Cfg) *Server {
	return &Server {
		address: cfg.TcpClient.TcpServerUrl,
	}
}

func (s *Server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Logger().With(zap.Error(err)).Fatal("Error starting TCP server.")
	}
	defer listener.Close()
	logger := log.Logger().With(zap.String("address", s.address))
	logger.Info("tcp server is started")
	for {
		logger.Info("waiting new incoming client ...")
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("error while accepting connection", zap.Error(err))
			return
		}
		
		uId, _ := uuid.NewV4()
		client := &Client{
			conn:        conn,
			Server:      s,
			ReceivedBuf: make([]byte, DefaultPacketSize),
			ReceivedLen: 0,
			UUID: uId.String(),
		}
		logger.Info("new incoming client: accepted", zap.String("uuid", client.UUID))
		go client.listen()
	}
}

// Read client data from channel
func (c *Client) listen() {
	log.Logger().Info("begin read")
	reader := bufio.NewReader(c.conn)
	tempBuf := make([]byte, DefaultPacketSize)
	for {
		n, err := reader.Read(tempBuf)
		if err != nil {
			if err != io.EOF {
				log.Logger().With(zap.Error(err)).Info("read error: eof")
			}
			_ = c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}

		if n == 0 {
			log.Logger().Info("read failed!")
			_ = c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}

		c.Server.onNewMessage(c, tempBuf, n)
	}
}

func (s *Server) onClientConnectionClosed(c *Client, err error) {
	log.Logger().With(zap.String("err", err.Error())).Warn("client closed")
	event := events.InternalMessageEvent {
		EventType: events.EventType_LOST_CONNECTION,
		MsgOneOf: &events.InternalMessageEvent_LostConnectionEvent {
			LostConnectionEvent: &events.LostConnectionEvent {
				ClientName: c.Name,
				ClientUuid: c.UUID,
			},
		},
		Token: "",
	}
	s.dispatch(&event)
}

func (s *Server) dispatch(event *events.InternalMessageEvent) {
	logger := log.Logger()
	logger.Info("got message: ", zap.String("message_type", event.EventType.String()))
	switch event.GetEventType() {
	case events.EventType_LOST_CONNECTION:
		s.handleLostConnection(event)
	case events.EventType_HEART_BEAT:
		s.handleHeartBeat(event)
	default:
		log.Logger().Warn("this command is not support right now")
	}
}

func (s *Server) handleLostConnection (event *events.InternalMessageEvent) {
	logger := log.Logger()
	logger.Info("lost connection")
	// todo
}

func (s *Server) handleHeartBeat (event *events.InternalMessageEvent) {
	logger := log.Logger()
	logger.Info("heart beat event")
	// todo
}

func (s *Server) onNewMessage(client *Client, data []byte, byteLen int) {
	logger := log.Logger().With(zap.Int("byteLen", byteLen))
	logger.Info("received")
	client.ReceivedBuf = make([]byte, byteLen)
	copy(client.ReceivedBuf[client.ReceivedLen:], data)
	client.ReceivedLen += byteLen
	var eatenByte = 0
	for eatenByte < client.ReceivedLen {
		msgLen := binary.LittleEndian.Uint32(client.ReceivedBuf[eatenByte : eatenByte+4])
		if msgLen > 1500 { //saint check
			client.ReceivedLen = 0
			break
		}
		if eatenByte == client.ReceivedLen {
			break
		}

		msgLenEnd := eatenByte + int(msgLen) + int(4)
		if msgLenEnd > client.ReceivedLen {
			break
		}
		// decode protobuf message
		event := events.InternalMessageEvent{}
		err := proto.Unmarshal(client.ReceivedBuf[eatenByte+4:msgLenEnd], &event)
		if err != nil {
			logger.Error("unmarshal failed")
			client.ReceivedLen = 0
			break
		}
		eatenByte = msgLenEnd

		s.dispatch(&event)
	}
	if eatenByte != 0 && eatenByte < client.ReceivedLen {
		copy(client.ReceivedBuf[0:client.ReceivedLen-eatenByte], client.ReceivedBuf[eatenByte:client.ReceivedLen])
		log.Logger().Info("shrink memory buffer")
	}
	client.ReceivedLen = client.ReceivedLen - eatenByte
	if client.ReceivedLen < 0 {
		logger.Error("ReceivedLen error", zap.Int("receivedLen", client.ReceivedLen), zap.Int("eatenByte", eatenByte))
		client.ReceivedLen = 0
	}
	log.Logger().Info("after execute ", zap.Int("remain_size", client.ReceivedLen))
}
