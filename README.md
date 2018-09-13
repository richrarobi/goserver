# goserver
simple golang server for distributed local linux systems including Raspberry Pi.#
From an original post by eliben.
I have extended the server to allow a form of RPC (but NOT using RPC as such).

Using websockets. A web based client also works, as per the example derived from that provided by Eli.

I am planning to grow the whole environment so any go or python3 client can access specific functions on other systems.
This will facilitate any Raspberry pi or other (Unix/Linux preferred) systems to use e.g maths functions on faster pc's,
or to do things like set led's remotely on one or more Raspi systems.

currently able to access system type on any (?) linux, and cpuTemp(erature) on Pi.

Next is to incorporate Blinkt functionality, and then look at maybe prime numbers, etc.

Rich R
