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

type TServer interface { // server接口
	ProcessorFactory() TProcessorFactory       // 返回TProcessorFactory
	ServerTransport() TServerTransport         // 返回TServerTransport
	InputTransportFactory() TTransportFactory  // 返回输入的TransportFactory
	OutputTransportFactory() TTransportFactory // 返回输出的TransportFactory
	InputProtocolFactory() TProtocolFactory
	OutputProtocolFactory() TProtocolFactory

	// Starts the server
	Serve() error // 启动Server
	// Stops the server. This is optional on a per-implementation basis. Not
	// all servers are required to be cleanly stoppable.
	Stop() error // 停止Server
}
