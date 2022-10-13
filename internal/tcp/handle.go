package tcp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Server interface {
	Start() error
}

type Client interface {
	SendFile(fid string, pkey, signmsg, sign []byte) error
}

type ConMgr struct {
	conn     NetConn
	dir      string
	fileName string

	// 发送的文件列表
	sendFiles []string

	// 在发送重要消息的时候，需要同步等待消息的状态，返回是否正确
	waitNotify chan bool
	stop       chan struct{}
}

func NewServer(conn NetConn, dir string) Server {
	return &ConMgr{
		conn: conn,
		dir:  dir,
		stop: make(chan struct{}),
	}
}

func (c *ConMgr) Start() error {
	c.conn.HandlerLoop()
	// 处理接收的消息
	return c.handler()
}

func (c *ConMgr) handler() error {
	var fs *os.File
	var err error

	defer func() {
		if fs != nil {
			_ = fs.Close()
		}
	}()

	for !c.conn.IsClose() {
		m, ok := c.conn.GetMsg()
		if !ok {
			return fmt.Errorf("close by connect")
		}
		if m == nil {
			continue
		}

		switch m.MsgType {
		case MsgHead:
			// 创建文件
			if m.FileName != "" {
				c.fileName = m.FileName
			} else {
				c.fileName = GenFileName()
			}

			fmt.Println("recv head fileName is", c.fileName)
			fs, err = os.OpenFile(filepath.Join(c.dir, c.fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				fmt.Println("os.Create err =", err)
				c.conn.SendMsg(NewNotifyMsg(c.fileName, Status_Err))
				return err
			}
			fmt.Println("send head is ok")

			c.conn.SendMsg(NewNotifyMsg(c.fileName, Status_Ok))
		case MsgFile:
			if fs == nil {
				fmt.Println(c.fileName, "file is not open !")
				c.conn.SendMsg(NewCloseMsg(c.fileName, Status_Err))
				return nil
			}
			// 写入文件
			_, err = fs.Write(m.Bytes)
			if err != nil {
				fmt.Println("file.Write err =", err)
				c.conn.SendMsg(NewCloseMsg(c.fileName, Status_Err))
				return err
			}
		case MsgEnd:
			// 操作完成
			info, _ := fs.Stat()
			if info.Size() != int64(m.FileSize) {
				err = fmt.Errorf("file.size %v rece size %v \n", info.Size(), m.FileSize)
				c.conn.SendMsg(NewCloseMsg(c.fileName, Status_Err))
				return err
			}

			fmt.Printf("save file %v is success \n", info.Name())
			c.conn.SendMsg(NewNotifyMsg(c.fileName, Status_Ok))

			fmt.Printf("close file %v is success \n", c.fileName)
			_ = fs.Close()
			fs = nil
		case MsgNotify:
			c.waitNotify <- m.Bytes[0] == byte(Status_Ok)
		case MsgClose:
			fmt.Printf("revc close msg ....\n")
			if m.Bytes[0] != byte(Status_Ok) {
				return fmt.Errorf("server an error occurred")
			}
			return nil
		}
	}

	return err
}

func NewClient(conn NetConn, dir string, files []string) Client {
	return &ConMgr{
		conn:       conn,
		dir:        dir,
		sendFiles:  files,
		waitNotify: make(chan bool, 1),
		stop:       make(chan struct{}),
	}
}

func (c *ConMgr) SendFile(fid string, pkey, signmsg, sign []byte) error {
	var err error
	c.conn.HandlerLoop()
	// 处理接收的消息
	go func() {
		_ = c.handler()
	}()
	err = c.sendFile(fid, pkey, signmsg, sign)
	return err
}

func (c *ConMgr) sendFile(fid string, pkey, signmsg, sign []byte) error {
	defer func() {
		_ = c.conn.Close()
	}()

	var err error
	for _, file := range c.sendFiles {
		err = c.sendSingleFile(filepath.Join(c.dir, file), fid, pkey, signmsg, sign)
		if err != nil {
			return err
		}
	}

	c.conn.SendMsg(NewCloseMsg(c.fileName, Status_Ok))
	return err
}

func (c *ConMgr) sendSingleFile(filePath string, fid string, pkey, signmsg, sign []byte) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file err %v \n", err)
		return err
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	fileInfo, _ := file.Stat()

	fmt.Println("client ready to write ", filePath)
	fmt.Println("    file id: ", fid)
	m := NewHeadMsg(fileInfo.Name(), fid, pkey, signmsg, sign)
	// 发送文件信息
	c.conn.SendMsg(m)

	// 等待服务器返回通知消息
	timer := time.NewTimer(10 * time.Second)
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send err")
		}
	case <-timer.C:
		return fmt.Errorf("wait server msg timeout")
	}

	for !c.conn.IsClose() {
		// 发送文件数据
		readBuf := BytesPool.Get().([]byte)

		n, err := file.Read(readBuf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		c.conn.SendMsg(NewFileMsg(c.fileName, readBuf[:n]))
	}

	c.conn.SendMsg(NewEndMsg(c.fileName, uint64(fileInfo.Size())))
	waitTime := fileInfo.Size() / 1024 / 10
	if waitTime < 5 {
		waitTime = 5
	}
	// 等待服务器返回通知消息
	timer = time.NewTimer(time.Second * time.Duration(waitTime))
	select {
	case ok := <-c.waitNotify:
		if !ok {
			return fmt.Errorf("send err")
		}
	case <-timer.C:
		return fmt.Errorf("wait server msg timeout")
	}

	fmt.Println("client send " + filePath + " file success...")
	return nil
}

// PathExists 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GenFileName() string {
	u := uuid.New()
	return u.String()

}
