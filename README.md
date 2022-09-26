# DISYS_MandatoryHandIn2

a) What are packages in your implementation? What data structure do
you use to transmit data and meta-data?

We used our own datastructure "packet" which contains three variables, syn, ack, hash and data. Syn, ack and hash are ints, while data is a string.
The meta-data is therefore syn, ack and hash, while the data we wish to transmit is the data string.
Hash is used to store the hashcode of the data, in order to implement checksums.


b) Does your implementation use threads or processes? Why is it not
realistic to use threads?

Our implementation revolves around threads. It is not realistic since it runs concurrenlty on the same device/server. Therefore it all runs in the same process. Channels also wait for each other, meaning it really is not possible for the client to send data without the server responding. If this was to happen, the client would just be stuck until the program ends.


c) How do you handle message re-orde

Message re-ordering is handled by the client sending a sequential index to the server, which then returns the index plus one. At the end of the data transmission, it is checked whether the server has recieved the correct packets in the correct order. If not, the last correct packet "syn" number is returned to the client, and an error message is provided. One could also implement a way to attempt to send the packet again.


d) How do you handle message loss?

Message loss is handled through the hash function. If parts of the message get corrupted or lost, the hashcode will very likely not be the same, as the one provided as message data. If that is the case, the last correct "syn" number is provided to the client. Here one could implement a way to send it again over the network, however we decided to just provide the user with an error message.


e) Why is the 3-way handshake important?

The 3-way handshake is important to synchronize and establish a reliable connection to the server. Here, the client makes sure the server knows that it is communicating with the client, and the server makes sure the client knows that the server knows that they are communicating. 