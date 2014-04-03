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
	"log"
)

// Simple, non-concurrent server for testing.
type TSimpleServer struct {
	stopped bool

	processorFactory       TProcessorFactory
	serverTransport        TServerTransport // server端的transport
	inputTransportFactory  TTransportFactory
	outputTransportFactory TTransportFactory
	inputProtocolFactory   TProtocolFactory
	outputProtocolFactory  TProtocolFactory
}

func NewTSimpleServer2(processor TProcessor, serverTransport TServerTransport) *TSimpleServer {
	return NewTSimpleServerFactory2(NewTProcessorFactory(processor), serverTransport)
}

func NewTSimpleServer4(processor TProcessor, serverTransport TServerTransport, transportFactory TTransportFactory, protocolFactory TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory4(NewTProcessorFactory(processor), // 默认用NewTProcessorFactory包装
		serverTransport,
		transportFactory,
		protocolFactory,
	)
}

func NewTSimpleServer6(processor TProcessor, serverTransport TServerTransport, inputTransportFactory TTransportFactory, outputTransportFactory TTransportFactory, inputProtocolFactory TProtocolFactory, outputProtocolFactory TProtocolFactory) *TSimpleServer {
	return NewTSimpleServerFactory6(NewTProcessorFactory(processor),
		serverTransport,
		inputTransportFactory,
		outputTransportFactory,
		inputProtocolFactory,
		outputProtocolFactory,
	)
}

func NewTSimpleServerFactory2(processorFactory TProcessorFactory, serverTransport TServerTransport) *TSimpleServer {
	return NewTSimpleServerFactory6(processorFactory,
		serverTransport,
		NewTTransportFactory(),
		NewTTransportFactory(),
		NewTBinaryProtocolFactoryDefault(),
		NewTBinaryProtocolFactoryDefault(),
	)
}

func NewTSimpleServerFactory4(processorFactory TProcessorFactory, serverTransport TServerTransport, transportFactory TTransportFactory, protocolFactory TProtocolFactory) *TSimpleServer { // transportFactory和protocolFactory相同
	return NewTSimpleServerFactory6(processorFactory,
		serverTransport,
		transportFactory,
		transportFactory,
		protocolFactory,
		protocolFactory,
	)
}

func NewTSimpleServerFactory6(processorFactory TProcessorFactory, serverTransport TServerTransport, inputTransportFactory TTransportFactory, outputTransportFactory TTransportFactory, inputProtocolFactory TProtocolFactory, outputProtocolFactory TProtocolFactory) *TSimpleServer { // 最终都会归结为对NewTSimpleServerFactory6的调用
	return &TSimpleServer{processorFactory: processorFactory,
		serverTransport:        serverTransport,
		inputTransportFactory:  inputTransportFactory,
		outputTransportFactory: outputTransportFactory,
		inputProtocolFactory:   inputProtocolFactory,
		outputProtocolFactory:  outputProtocolFactory,
	} // 分别赋值6个域
}

func (p *TSimpleServer) ProcessorFactory() TProcessorFactory {
	return p.processorFactory
}

func (p *TSimpleServer) ServerTransport() TServerTransport {
	return p.serverTransport
}

func (p *TSimpleServer) InputTransportFactory() TTransportFactory {
	return p.inputTransportFactory
}

func (p *TSimpleServer) OutputTransportFactory() TTransportFactory {
	return p.outputTransportFactory
}

func (p *TSimpleServer) InputProtocolFactory() TProtocolFactory {
	return p.inputProtocolFactory
}

func (p *TSimpleServer) OutputProtocolFactory() TProtocolFactory {
	return p.outputProtocolFactory
}

func (p *TSimpleServer) Serve() error { // 启动server执行服务
	p.stopped = false
	err := p.serverTransport.Listen() // 调用serverTransport的listen，开始监听
	if err != nil {
		return err
	}
	for !p.stopped { // 如果server未停止
		client, err := p.serverTransport.Accept() // 接收新服务，返回一个对应client端的transport
		if err != nil {
			log.Println("Accept err: ", err)
		}
		if client != nil { // 如果client非空，创建一个goroutine执行processRequest
			go func() {
				if err := p.processRequest(client); err != nil {
					log.Println("error processing request:", err)
				}
			}()
		}
	}
	return nil
}

func (p *TSimpleServer) Stop() error { // 停止server
	p.stopped = true
	p.serverTransport.Interrupt()
	return nil
}

func (p *TSimpleServer) processRequest(client TTransport) error {
	processor := p.processorFactory.GetProcessor(client)
	inputTransport := p.inputTransportFactory.GetTransport(client) // 获得输入的transport
	outputTransport := p.outputTransportFactory.GetTransport(client) // 获得输出的transport
	inputProtocol := p.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := p.outputProtocolFactory.GetProtocol(outputTransport)
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}
	for {
		ok, err := processor.Process(inputProtocol, outputProtocol)
		if err, ok := err.(TTransportException); ok && err.TypeId() == END_OF_FILE {
			return nil
		} else if err != nil {
			return err
		}
		if !ok || !inputProtocol.Transport().Peek() {
			break
		}
	}
	return nil
}
