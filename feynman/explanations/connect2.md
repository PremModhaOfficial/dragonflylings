
- redis Clinet is just the Instance of Clinet struct that lazily connects to the redis server at first request

1. so the Client not nil means the clinet was created ; but the network connection (TCP)  is yet to happen
2. in case or redis client the Client is just a paper on which we wrote the address when the person reaches the address they find out if the place is actually open or not; even the path to that address is open/active!
3. distributed systems talk to eachother but in this complex netwoek no one knos the state of the next in line service ; it may be down or working in load or be just OK ; so we need to see if it is running (a quick health-check to know the state of the service) so we need to do that  everytiem to chek of availability
