
Moteminder is a very simple program intended for Go team developers (with gomote
access) to ping existing gomotes, keeping them alive for some specific duration.
Note that active commands and ssh sessions will already keep a mote from being
GC'd; the main use here is for motes that are idle but still needed in the
future.

