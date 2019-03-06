package protocol

import "bufio"

type Text struct {
}

func (t Text) ReadString(reader *bufio.Reader) (interface{}, error) {
	msg,err := reader.ReadString('\n')
	if err != nil {
		return msg,err
	}
	msg = msg[1:]
	return msg,err
}

func (t Text) WriteString(msg interface{}) []byte {
	strMsg := msg.(string) + "\n"
	return []byte(strMsg)
}