#ifndef TUNTAP_H
#define TUNTAP_H

int tuntap_setup(int fd, unsigned char *name, int mode, int packet_info);

#endif
