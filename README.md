# go_tftpd
[![Build Status](https://travis-ci.org/mark-rushakoff/go_tftpd.png?branch=master)](https://travis-ci.org/mark-rushakoff/go_tftpd)

An implementation of a TFTP daemon written in Go.

## Usage

go_tftpd should currently be considered alpha status.
It serves files from the current working directory and it does not yet respect any TFTP options (e.g. larger block size).

If invoked as `go run main.go` it will bind to 127.0.0.1 and port 69 (the default TFTP server port).
The host and port can be overridden with `-host` and `-port` respectively.

## Implementation notes

This implementation aims to be two things:

1. A functional and useful TFTP daemon suitable for general use.
2. A Go project exemplary of clean code, common idioms, and code organization when designing a server that serves multiple clients simultaneously.

To summarize the principles used in designing this TFTP daemon:

* Domain objects:
  * As small as reasonably possible so that they can easily be unit tested
  * Avoid concurrency; prefer callbacks and interfaces (it's always easy to add concurrency later in Go)
  * Avoid creating other domain objects; prefer other objects to be passed as arguments (this simplifies unit testing)
* "Glue" code (e.g. [SafePacketProvider](safepacketprovider/safe_packet_provider.go)):
  * Concurrency encouraged where appropriate
  * Should accept minimal arguments necessary to create all needed domain objects
* The main function:
  * Coordinates glue code

### Standards-Compliant

This implementation aims to be standards-compliant, following the guidelines set forth in the following RFCs from IETF:

[RFC 1350](http://tools.ietf.org/html/rfc1350): The TFTP protocol (Rev 2)

- [x] Respond to read requests
- [x] Files transferred in 512-byte chunks
- [x] Data packets re-sent if no ack received in time
- [x] Send error for requests to files that do not exist
- [x] Send error for requests to files that exist but cannot be opened
- [x] Send error for very old ack
- [ ] Send error for file that has an error partway through reading
- [x] Handle netascii read requests
- [ ] Handle octet read requests

- [ ] Respond to write requests
- [ ] Ack packets re-sent if no data received in time

[RFC 1123, Section 4.2](http://tools.ietf.org/html/rfc1123#page-44): Requirements for internet hosts, TFTP

- [x] 4.2.2.1 Transfer mode "mail" is not supported
- [ ] 4.2.3.1 Sorcerer's Apprentice Syndrome addressed
- [ ] 4.2.3.2 Adaptive timeout (exponential backoff)
- [ ] 4.2.3.4 Access control (SHOULD include configurable access control of allowed pathnames; currently uses current working directory)
- [ ] 4.2.3.5 A TFTP request directed to a broadcast address SHOULD be silently ignored.

[RFC 2347](http://tools.ietf.org/html/rfc2347): TFTP Option Extension

- [x] Parses options

[RFC 2348](http://tools.ietf.org/html/rfc2348): TFTP Blocksize option

- [ ] Responds with OACK for block size
- [ ] Respects block size option

## License

go_tftpd is available under the terms of the MIT license.
