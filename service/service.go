package service

import (
	"crypto/md5"
	"fmt"
	"hash"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/1939323749/drcom_go/conf"
)

const RETRYTIMES = 5

const (
	_codeIn       = byte(0x03)
	_codeOut      = byte(0x06)
	_type         = byte(0x01)
	_eof          = byte(0x00)
	_controlCheck = byte(0x20)
	_adapterNum   = byte(0x05)
	_ipDog        = byte(0x01)
)

var (
	_delimiter   = []byte{0x00, 0x00, 0x00, 0x00}
	_emptyIP     = []byte{0, 0, 0, 0}
	_primaryDNS  = []byte{10, 10, 10, 10}
	_dhcpServer  = []byte{0, 0, 0, 0}
	_authVersion = []byte{0x6a, 0x00}
	_magic1      = big.NewInt(1968)
	_magic2      = big.NewInt(int64(0xffffffff))
	_magic3      = big.NewInt(int64(711))
)

type Service struct {
	c              *conf.Config
	md5Ctx         hash.Hash
	salt           []byte // [4:8]
	clientIP       []byte // [20:24]
	md5a           []byte
	tail1          []byte
	tail2          []byte
	keepAliveVer   []byte // [28:30]
	conn           *net.UDPConn
	ChallengeTimes int
	Count          int
	LogoutCh       chan struct{}

	logger  Logger
	Restart chan int
}

// Error exit message
type Error struct {
	Err error
	Msg string
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		md5Ctx:         md5.New(),
		md5a:           make([]byte, 16),
		tail1:          make([]byte, 16),
		tail2:          make([]byte, 4),
		keepAliveVer:   []byte{0xdc, 0x02},
		clientIP:       make([]byte, 4),
		salt:           make([]byte, 4),
		ChallengeTimes: 0,
		Count:          0,
		LogoutCh:       make(chan struct{}, 1),
		logger: Logger{
			loggers: map[string]*log.Logger{
				"info":  log.New(os.Stdout, "[INFO]", log.LstdFlags),
				"error": log.New(os.Stderr, "[ERROR]", log.LstdFlags),
			},
		},
		Restart: make(chan int, 1),
	}
	var (
		err     error
		udpAddr *net.UDPAddr
	)
	if udpAddr, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", c.AuthServer, c.Port)); err != nil {
		log.Fatalf("net.ResolveUDPAddr(udp4, %s) error(%v) ", fmt.Sprintf("%s:%s", c.AuthServer, c.Port), err)
	}
	if s.conn, err = net.DialUDP("udp", nil, udpAddr); err != nil {
		log.Fatalf("net.DialUDP(udp, %v, %v) error(%v)", nil, udpAddr, err)
	}
	return
}

func (s *Service) ReStart() error {
	err := s.Start()
	if err != nil {
		return err
	}
	return nil
}

// Start start drcom client.
func (s *Service) Start() error {
	s.logger.Info("start ...")
	s.logger.Info("challenge ...")
	if err := s.Challenge(s.ChallengeTimes); err != nil {
		s.logger.Error(fmt.Sprintf("drcomSvc.Challenge(%d) error(%v)", s.ChallengeTimes, err))
		return err
	}
	s.ChallengeTimes++
	s.logger.Info("ok")
	s.logger.Info("login ...")
	if err := s.Login(); err != nil {
		s.logger.Error(fmt.Sprintf("drcomSvc.Login() error(%v)", err))
		return err
	}
	s.logger.Info("alive daemon start ...")
	go s.aliveproc()
	s.logger.Info("alive daemon started")

	s.logger.Info("logout daemon start ...")
	go s.logoutproc()
	s.logger.Info("logout daemon started")

	s.logger.Info("connect daemon start ...")
	go s.checkConnect()
	s.logger.Info("connect daemon started")
	return nil
}

func (s *Service) aliveproc() {
	count := 0
	for {
		select {
		case _, ok := <-s.LogoutCh:
			if !ok {
				s.logger.Info("keep-aliveproc goroutine get a logout signal, exit")
				return
			}
		default:
		}
		count++
		s.logger.Info(fmt.Sprintf("keep-aliveproc ... %d", count))
		if err := s.Alive(); err != nil {
			s.logger.Error(fmt.Sprintf("drcomSvc.Alive() error(%v)", err))
			time.Sleep(time.Second * 5)
			continue
		}
		s.logger.Info("ok")
		time.Sleep(time.Second * 20)
	}
}

func (s *Service) logoutproc() {
	if _, ok := <-s.LogoutCh; !ok {
		s.logger.Info("logout ...")
		if err := s.Challenge(s.ChallengeTimes); err != nil {
			s.logger.Error(fmt.Sprintf("drcomSvc.Challenge(%d) error(%v)", s.ChallengeTimes, err))
			return
		}
		s.ChallengeTimes++
		if err := s.Logout(); err != nil {
			s.logger.Error(fmt.Sprintf("service.Logout() error(%v)", err))
			return
		}
		s.logger.Info("ok")
	}
}

func (s *Service) checkConnect() {
	retriedTimes := 0
	for {
		time.Sleep(2 * time.Second)
		go func() {
			ok, err := s.checkConnectivity()
			if !ok {
				retriedTimes++
				s.logger.Error(fmt.Sprintf("connectivity check failed, retry %d times, %s", retriedTimes, err))
			}
			return
		}()
		if retriedTimes >= RETRYTIMES {
			err := s.Close()
			if err != nil {
				return
			}
			s.logger.Error("connectivity check failed, logout")
			return
		}
	}
}

// Close close service.
func (s *Service) Close() error {
	close(s.LogoutCh)
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
