/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import (
	"net"
	"time"
)

type TServerSocket struct { // 服务端socket，实现了TServerTransport接口
	listener      net.Listener  // 监听的Listener结构
	addr          net.Addr      // 网络监听地址
	clientTimeout time.Duration // 客户端超时时间
	interrupted   bool          // 标志位，当置位时表示应该打断当前堵塞的accept和listen
}

func NewTServerSocket(listenAddr string) (*TServerSocket, error) { // 创建一个新的TServer socket
	return NewTServerSocketTimeout(listenAddr, 0)
}

func NewTServerSocketTimeout(listenAddr string, clientTimeout time.Duration) (*TServerSocket, error) { // 加入设置客户端超时
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &TServerSocket{addr: addr, clientTimeout: clientTimeout}, nil
}

func (p *TServerSocket) Listen() error { // 启动Listen
	if p.IsListening() {
		return nil
	}
	l, err := net.Listen(p.addr.Network(), p.addr.String()) // 开始执行监听
	if err != nil {
		return err
	}
	p.listener = l // 设置Listener结构
	return nil
}

func (p *TServerSocket) Accept() (TTransport, error) {
	if p.interrupted {
		return nil, errTransportInterrupted
	}
	if p.listener == nil {
		return nil, NewTTransportException(NOT_OPEN, "No underlying server socket")
	}
	conn, err := p.listener.Accept() // 接收新请求，返回conn结构
	if err != nil {
		return nil, NewTTransportExceptionFromError(err)
	}
	return NewTSocketFromConnTimeout(conn, p.clientTimeout), nil // 创建对应客户端的TSocket
}

// Checks whether the socket is listening.
func (p *TServerSocket) IsListening() bool { // 检查server socket是否处于listen状态
	return p.listener != nil
}

// Connects the socket, creating a new socket object if necessary.
func (p *TServerSocket) Open() error {
	if p.IsListening() {
		return NewTTransportException(ALREADY_OPEN, "Server socket already open")
	}
	if l, err := net.Listen(p.addr.Network(), p.addr.String()); err != nil {
		return err
	} else {
		p.listener = l
	}
	return nil
}

func (p *TServerSocket) Addr() net.Addr { // 获取TServerSocket的地址
	return p.addr
}

func (p *TServerSocket) Close() error { // 关闭TServerSocket
	defer func() {
		p.listener = nil
	}()
	if p.IsListening() {
		return p.listener.Close()
	}
	return nil
}

func (p *TServerSocket) Interrupt() error {
	p.interrupted = true
	return nil
}
