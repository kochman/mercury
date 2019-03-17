# Mercury

## About

Mercury is a decentralized, peer-to-peer, end-to-end encrypted messaging system that relies on chaining together people (nodes) to deliver messages. Nodes that can't directly connect to each other can exchange messages through mutual peers.

No internet connection is required, and all message exchanges happen directly between nodes. This provides privacy benefits since it keeps user information from leaving the network. Nobody can read messages except for the recipient. It also provides resiliency and flexibility (perhaps in disaster or conflict zones) since ad-hoc network connections can be created by any device and used to exchange messages without any centralized infrastructure.

Peers are discovered on the local network using DNS-SD and Zeroconf. Additionally, if internet is available, an App Engine service running on Google Cloud Platform is utilized to find other users globally and attempt to initiate connections to them. It serves as an additional registry of peers in addition to Zeroconf, but it is not necessary to Mercury's operation.

If two peers can communicate, the user is prompted to add them to the contacts book. Then users can send encrypted messages to each other. Sent messages are broadcast to all connected nodes, and if the node receives a message for itself, then it presents it to the user. Messages destined for other nodes are forwarded so that they eventually arrive at the recipient.

## Security

Mercury utilizes RSA to generate a unique public private keypair when it is first started. The address book allows users to associate public keys with friendly names. This ensures that only the recipient of a message can decrypt it—only the peer that has the corresponding private key is able to. This means that any node can help forward encrypted messages to other nodes, but it doesn't know what the messages contain unless it is the intended recipient and it has the corresponding private key.

## API

`/api/self` - displays user's own public key to share with contacts

`/api/messages` - JSON of messages that were sent to user (decrypted)

`/api/contacts/all` - JSON of all added contacts

`/api/contacts/peers` - JSON of peers that can be added

## User endpoints

`/messages` - card display view of messages received ordered by latest -> earliest

`/contacts` - page for user to view current contacts and discover new ones

## Peer endpoints

`/messages` - all encrypted messages that are forwarded to peers

`/peers` - list of peers that this peer knows about

`/pubkey` - this peer's friendly name and public key, for easy discovery and contact creation

## Building and running

Requires Go 1.11 or greater.

```
git clone git@github.com:kochman/mercury.git
cd mercury
go build .
./mercury
```

## Usage

Once the application is built and running, interact with the web application on `localhost:3000`

If there are others running Mercury on the same network, they will show up in `/contacts`.

Once a contact has been added, user can send messages to the new contact from the root page.

Messages will be persisted, and can be passed from network to network as machines move across networks.