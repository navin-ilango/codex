# Tech Fundamentals
## Encryption

```
Encryption is a process that takes a plaintext and converts it into a CipherText using an algorithm and and a Key.
```

`Encryption approaches`
1. Encryption At Rest
2. Encryption In Transit.

`Encryption Concepts`
1. Plain Text
2. Algorithms - AWS, RC4, DES, RC5, RC6.
3. Key
4. CipherText

`Types of Encryption`
1. `Symmetric Key Encryption` - Same Key is used to both encrypt a plaintext and decrypt a ciphertest.
2. `Asymmetric Key Encryption` - The encryption key is formed of two parts - Public Key and Private Key. The Public key can only encrypt and the private key can only decrypt.

## Network Starter Pack

`OSI Model - 7 Layers`
7. Application
6. Presentation
5. Session
4. Transport
3. Network
2. DataLink
1. Physical

`Layer 1 - Physical` 

The physical layer specifications define the transmission and reception of RAW Bit Streams between a device and a shared physical medium. It defines things like voltage levels, timing, rates, distances, modulation and connectors(hub). 

Layer 1 Limitations:
1. No device addressing, all data is processed by all devices.
2. Anything received on any port is transmitted on ever other port by the physical hub.
3. If multiple devices transmit at once, a collition occurs.
4. L1 has no media access control and no collision detection.

Layer 1 Summary:
1. Physical shared medium.
2. Standards fo transmitting onto the medium
3. Standards for receiving from the medium.
4. No access control.
5. No uniqueky identified devices.
6. No device-to-device communications.

`Layer 2 - Data Link` 

1. Devices at L2 have a Unique Hardware MAC Address. This allows network frames to be addressed to a destination or broadcast.
2. A `FRAME` consists of `SFD` - Starting Frame Delimites which allows to find the start of a frame, `MAC Header` - Contains the Source and Destination MAC Address, `PAYLOAD` - is the data the frame carries from source to destination. Its generally provided by L3 and the `Ether Type` attribute defines which L3 protocol is used (ex: IP), `FCS` - Frame Check Sequence is used the check if error exists in the frame.
3. L2 provides Unique MAC Address for each devices which enables device to device communication, media access control and performs carrier checks to avoid collision which 2 devices are trying to transmit at the same time.
4. `Carrier Checks` is a process where L2 checks of L1 is already transmitting or receiving any data. If so, L2 will wai for L1 to complete its data transmission before sending new data - thus collision is avoided.
5. Layer 2 uses `SWITCH` instead if Hub. Switch understand frames and MAC Addresses. They maintain a MAC address table which starts off empty. As the switch receives frames on its ports, it learns which devices are connected and populates the MAC address table. This makes a SWITCH intelligent to route frames to its destination only instead of broadcasting to all the connected devices like a hub.

`Layer 3 - Network`

1. Packets - IP Address
2. Routers - Routes and Route Tables
3. Subnet Mask
4. Address Resolution Protocol - ARP

Summary:

1. IP Addresses - Cross Network Addressing.
2. ARP - Find the MAC Address for this IP.
3. Routes - Where to forward this Packet.
4. Route Tables - Multiple Routes.
5. Router - Moves packets from source to destination with encapsulation in L2 on the way.
6. Device - to - Device communication over internet.
7. SRC IP to DEST IP Only - No other method of communication.
8. Packets can be delivered out of order.


`Layer 4 - Transport & Layer 5 - Session`

1. `TCP / UDP` - TCP is connection based protocol. A connection is established between two devices using a random port on a client and a known port on the server. Once established the connection is bi-directional. The connection is a reliable connection provided via the segments encapsulated in IP Packets.
2. `TCP Segment` consists of, `SRC PORT / DST PORT`, `SEQUENCE NUMBER`, `ACK`, `FLAGS`, `WINDOW`, `CHECKSUM`, `URGENT POINTER`.

`Firewals`

A `STATELESS FIREWALL` (AWS Network ACLs) would see 2 things: OUTBOUND & INBOUND
A `STATEFULL FIREWALL` (AWS SG) would see 1 thing: OUTBOUND. So allowing the outboud implicity allows the inbound response.

`NAT - Network Address Translation`

1. NAT is designed to overcome IPv4 Shortages.
2. Translated Private IPv$ address to Public.
3. Also provides some security benefits.
4. No need of MAT for IPv6.

NAT Types
1. STATIC NAT - 1 Private IP Address to 1 Fixed Public Address. Ex: Internet Gateway.
2. DYNAMIC NAT - 1 Private Addres to 1 available Public Address
3. PORT ADDRESS TRANSLATION (PAT) - Many Private address to 1 Public Address. Ex: NAT Gateway

`IP Addressing and Subnetting`
1. IPv4 standard created in 1981 (RFC 791).
2. 4,294,967,296 (4.3 Billion) Addresses.
3. Address Space: 0.0.0.0 --> 255.255.255.255
4. All public IPv4 addressing is Allocated.
5. Part of the address space is private and can be used/re-used freely.
6. `Private Address Space`: 10.0.0.0 - 10.255.255.255 (1 x Class A Network), 172.16.0.0 - 172.31.255.255 (16 x Class B Network), 192.168.0.0 - 192.168.255.0 (256 x Class C Network)

`Distributed Denial of Service`
1. Attacks designed to overload websites.
2. Compete against Legitimate Connections.
3. Distributed - Hard to block individual IPs / Ranges.
4. Involves large armies of compromised machines(bots).

Attacks
1. Application Layer Attack - HTTP Flood
2. Protocol Attack - SYN Flood
3. Volumetric Attacks - DNS Amplification

`SSL and TLS`
1. Privacy and Data Integrity between a client and server.
2. Comminications are encrypted in transit.
3. Identity (Server or Client & Server) verified.
4. Reliable Connection - Protects against alteration.

`Hashing`
1. Data + Hash Function = HASH
2. Same Data = Same Hash, Diff Data = Diff Hash.
3. Downloaded data can be verified by using the Hash.

`DNS 101 - Key Terms`
1. `DNS Zone` - A database containing records. Ex: netflix.com
2. `Zone File` - The files storing the zone on disk.
3. `Name Server` - A DNS Server whic hosts one or more zones and stores its Zone Files.
4. `Authoritative` - Contains real / genuine records.
5. `Non Authoritative / Cached` - Copies of records / zones stored elsewhere to t speed DNS resolution. 
