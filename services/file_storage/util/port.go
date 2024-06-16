package util

import "net"

func GetAvailablePort() (int, error) {
  var tcp *net.TCPAddr
  var err error

  if tcp, err = net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
    return 0, err
  }

  var listener *net.TCPListener
  if listener, err = net.ListenTCP("tcp", tcp); err != nil {
    return 0, err
  }
  defer listener.Close()

  return listener.Addr().(*net.TCPAddr).Port, nil
}
