/*
 * ZeroTier SDK - Network Virtualization Everywhere
 * Copyright (C) 2011-2018  ZeroTier, Inc.  https://www.zerotier.com/
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * --
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial closed-source software that incorporates or links
 * directly against ZeroTier software without disclosing the source code
 * of your own application.
 */

/**
 * @file
 *
 * Application-facing, partially-POSIX-compliant socket API
 */

#ifndef LIBZT_DEFINES_H
#define LIBZT_DEFINES_H

#define LIBZT_IPV4 1
#define LIBZT_IPV6 1

/**
 * Default port that libzt will use to support all virtual communication
 */
#define LIBZT_DEFAULT_PORT 9994

#define NO_STACK   0 // for layer-2 only (this will omit all userspace network stack code)

/**
 * Maximum MTU size for ZeroTier
 */
#define ZT_MAX_MTU	10000

/**
 * How fast service states are re-checked (in milliseconds)
 */
#define ZTO_WRAPPER_CHECK_INTERVAL	100

/**
 * Length of buffer required to hold a ztAddress/nodeID
 */
#define ZTO_ID_LEN	16

#if defined(_MSC_VER)
#include <BaseTsd.h>
typedef SSIZE_T ssize_t;
#endif

/****************************************************************************/
/* For SOCK_RAW support, it will initially be modeled after linux's API, so */
/* below are the various things we need to define in order to make this API */
/* work on other platforms. Mayber later down the road we will customize	 */
/* this for each different platform. Maybe.											*/
/****************************************************************************/

#if !defined(__linux__)
#define SIOCGIFINDEX  101
#define SIOCGIFHWADDR 102

// Normally defined in linux/if_packet.h, defined here so we can offer a linux-like
// raw socket API on non-linux platforms
struct sockaddr_ll {
		unsigned short	sll_family;	/* Always AF_PACKET */
		unsigned short	sll_protocol; /* Physical layer protocol */
		int				sll_ifindex;  /* Interface number */
		unsigned short	sll_hatype;	/* ARP hardware type */
		unsigned char	sll_pkttype;  /* Packet type */
		unsigned char	sll_halen;	 /* Length of address */
		unsigned char	sll_addr[8];  /* Physical layer address */
};

#endif

/****************************************************************************/
/* lwIP                                                                     */
/****************************************************************************/

// For LWIP configuration see: include/lwipopts.h

/* The following three quantities are related and govern how incoming frames are fed into the 
network stack's core:

Every LWIP_GUARDED_BUF_CHECK_INTERVAL milliseconds, a callback will be called from the core and 
will input a maximum of LWIP_FRAMES_HANDLED_PER_CORE_CALL frames before returning control back
to the core. Meanwhile, incoming frames from the ZeroTier wire will be allocated and their 
pointers will be cached in the receive frame buffer of the size LWIP_MAX_GUARDED_RX_BUF_SZ to 
await the next callback from the core */

#define LWIP_GUARDED_BUF_CHECK_INTERVAL	50 // in ms
#define LWIP_MAX_GUARDED_RX_BUF_SZ	1024 // number of frame pointers that can be cached waiting for receipt into core
#define LWIP_FRAMES_HANDLED_PER_CORE_CALL 16 // How many frames are handled per call from core

typedef signed char err_t;

#define ND6_DISCOVERY_INTERVAL 1000
#define ARP_DISCOVERY_INTERVAL ARP_TMR_INTERVAL

/**
 * Specifies the polling interval and the callback function that should
 * be called to poll the application. The interval is specified in
 * number of TCP coarse grained timer shots, which typically occurs
 * twice a second. An interval of 10 means that the application would
 * be polled every 5 seconds. (only for raw lwIP driver)
 */
#define LWIP_APPLICATION_POLL_FREQ 2

/**
 * TCP timer interval in milliseconds (only for raw lwIP driver)
 */
#define LWIP_TCP_TIMER_INTERVAL 25

/**
 * How often we check VirtualSocket statuses in milliseconds (only for raw lwIP driver)
 */
#define LWIP_STATUS_TMR_INTERVAL 500

// #define LWIP_CHKSUM <your_checksum_routine>, See: RFC1071 for inspiration

/****************************************************************************/
/* Defines                                                                  */
/****************************************************************************/

/**
 * Maximum number of sockets that libzt can administer
 */
#define ZT_MAX_SOCKETS 1024

/**
 * Maximum MTU size for libzt (must be less than or equal to ZT_MAX_MTU)
 */
#define ZT_SDK_MTU ZT_MAX_MTU

/**
 *
 */
#define ZT_LEN_SZ 4

/**
 *
 */
#define ZT_ADDR_SZ 128

/**
 * Size of message buffer for VirtualSockets
 */
#define ZT_SOCKET_MSG_BUF_SZ ZT_SDK_MTU + ZT_LEN_SZ + ZT_ADDR_SZ

/**
 * Polling interval (in ms) for file descriptors wrapped in the Phy I/O loop (for raw drivers only)
 */
#define ZT_PHY_POLL_INTERVAL 5

/**
 * State check interval (in ms) for VirtualSocket state
 */
#define ZT_ACCEPT_RECHECK_DELAY 50

/**
 * State check interval (in ms) for VirtualSocket state
 */
#define ZT_CONNECT_RECHECK_DELAY 50

/**
 * State check interval (in ms) for VirtualSocket state
 */
#define ZT_API_CHECK_INTERVAL 50

/**
 * Size of TCP TX buffer for VirtualSockets used in raw network stack drivers
 */
#define ZT_TCP_TX_BUF_SZ 1024 * 1024 * 128

/**
 * Size of TCP RX buffer for VirtualSockets used in raw network stack drivers
 */
#define ZT_TCP_RX_BUF_SZ 1024 * 1024 * 128

/**
 * Size of UDP TX buffer for VirtualSockets used in raw network stack drivers
 */
#define ZT_UDP_TX_BUF_SZ ZT_MAX_MTU

/**
 * Size of UDP RX buffer for VirtualSockets used in raw network stack drivers
 */
#define ZT_UDP_RX_BUF_SZ ZT_MAX_MTU * 10

/**
 * Send buffer size for the network stack
 * By default picoTCP sets them to 16834, this is good for embedded-scale
 * stuff but you might want to consider higher values for desktop and mobile
 * applications.
 */
#define ZT_STACK_TCP_SOCKET_TX_SZ ZT_TCP_TX_BUF_SZ

/**
 * Receive buffer size for the network stack
 * By default picoTCP sets them to 16834, this is good for embedded-scale
 * stuff but you might want to consider higher values for desktop and mobile
 * applications.
 */
#define ZT_STACK_TCP_SOCKET_RX_SZ ZT_TCP_RX_BUF_SZ

/**
 * Maximum size we're allowed to read or write from a stack socket
 * This is put in place because picoTCP seems to fail at higher values.
 * If you use another stack you can probably bump this up a bit.
 */
#define ZT_STACK_SOCKET_WR_MAX 4096

/**
 * Maximum size of read operation from a network stack
 */
#define ZT_STACK_SOCKET_RD_MAX 4096*4

/**
 * Maximum length of libzt/ZeroTier home path (where keys, and config files are stored)
 */
#define ZT_HOME_PATH_MAX_LEN 256

/**
 * Length of human-readable MAC address string
 */
#define ZT_MAC_ADDRSTRLEN 18

/**
 * Everything is ok
 */
#define ZT_ERR_OK 0

/**
 * Value returned during an internal failure at the VirtualSocket/VirtualTap layer
 */
#define ZT_ERR_GENERAL_FAILURE -88

/**
 * Whether sockets created will have SO_LINGER set by default
 */
#define ZT_SOCK_BEHAVIOR_LINGER	false
	
/**
 * Length of time that VirtualSockets should linger (in seconds)
 */
#define ZT_SOCK_BEHAVIOR_LINGER_TIME 3

/**
 * Maximum wait time for socket closure if data is still present in the write queue
 */
#define ZT_SDK_CLTIME 60

/**
 * Interval for performing background tasks (such as adding routes) on VirtualTap objects (in seconds)
 */
#define ZT_HOUSEKEEPING_INTERVAL 1

/****************************************************************************/
/* Socket API Signatures                                                    */
/****************************************************************************/

#define ZT_SETSOCKOPT_SIG int fd, int level, int optname, const void *optval, socklen_t optlen
#define ZT_GETSOCKOPT_SIG int fd, int level, int optname, void *optval, socklen_t *optlen
#define ZT_SENDMSG_SIG int fd, const struct msghdr *msg, int flags
#define ZT_SENDTO_SIG int fd, const void *buf, size_t len, int flags, const struct sockaddr *addr, socklen_t addrlen
#define ZT_RECV_SIG int fd, void *buf, size_t len, int flags
#define ZT_RECVFROM_SIG int fd, void *buf, size_t len, int flags, struct sockaddr *addr, socklen_t *addrlen
#define ZT_RECVMSG_SIG int fd, struct msghdr *msg,int flags
#define ZT_SEND_SIG int fd, const void *buf, size_t len, int flags
#define ZT_READ_SIG int fd, void *buf, size_t len
#define ZT_WRITE_SIG int fd, const void *buf, size_t len
#define ZT_SHUTDOWN_SIG int fd, int how
#define ZT_SOCKET_SIG int socket_family, int socket_type, int protocol
#define ZT_CONNECT_SIG int fd, const struct sockaddr *addr, socklen_t addrlen
#define ZT_BIND_SIG int fd, const struct sockaddr *addr, socklen_t addrlen
#define ZT_LISTEN_SIG int fd, int backlog
#define ZT_ACCEPT4_SIG int fd, struct sockaddr *addr, socklen_t *addrlen, int flags
#define ZT_ACCEPT_SIG int fd, struct sockaddr *addr, socklen_t *addrlen
#define ZT_CLOSE_SIG int fd
#define ZT_POLL_SIG struct pollfd *fds, nfds_t nfds, int timeout
#define ZT_SELECT_SIG int nfds, fd_set *readfds, fd_set *writefds, fd_set *exceptfds, struct timeval *timeout
#define ZT_GETSOCKNAME_SIG int fd, struct sockaddr *addr, socklen_t *addrlen
#define ZT_GETPEERNAME_SIG int fd, struct sockaddr *addr, socklen_t *addrlen
#define ZT_GETHOSTNAME_SIG char *name, size_t len
#define ZT_SETHOSTNAME_SIG const char *name, size_t len
#define ZT_FCNTL_SIG int fd, int cmd, int flags
#define ZT_IOCTL_SIG int fd, unsigned long request, void *argp
#define ZT_SYSCALL_SIG long number, ...

#define LIBZT_SERVICE_NOT_STARTED_STR "service not started yet, call zts_startjoin()"

#endif // _H
