package stats

import (
	"strconv"

	"github.com/elastic/gosigar/sys/linux"
	"github.com/yuuki/lstf/netutil"
)

// AddrPort are <addr>:<port>
type AddrPort struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

// Stat represents a socket statistics.
type Stat struct {
	Proto    uint8          `json:"proto"`
	RecvQ    uint32         `json:"recvq"`
	SendQ    uint32         `json:"sendq"`
	Local    *AddrPort      `json:"local"`
	Foreign  *AddrPort      `json:"foreign"`
	State    string         `json:"state"`
	Program  string         `json:"program"`
}

// Stats represents a group of Stat.
type Stats []*Stat

// GetStats gets socket statistics by Linux netlink API.
func GetStats(all bool, program bool) (Stats, error) {
	conns, err := netutil.NetlinkConnections()
	if err != nil {
		return nil, err
	}
	var stats Stats
	for _, conn := range conns {
		if !all {
			if linux.TCPState(conn.State) == linux.TCP_LISTEN {
				continue
			}
		}
		// FIXME
		pidprogs := ""
		if program {
			// get program name/pid from inode and procfs
			pidprogs = "PID/ProgramName"
		}
		stats = append(stats, &Stat{
			Proto:   conn.Family,
			RecvQ:   conn.RQueue,
			SendQ:   conn.WQueue,
			Local:   &AddrPort{
				Addr: conn.SrcIP().String(),
				Port: strconv.Itoa(conn.SrcPort()),
			},
			Foreign: &AddrPort{
				Addr: conn.DstIP().String(),
				Port: strconv.Itoa(conn.DstPort()),
			},
			State:   linux.TCPState(conn.State).String(),
			Program: pidprogs,
		})
	}
	return stats, nil
}

// ReplaceLookupedName replaces s.Local/Foreign.Addr into lookuped name.
func (s *Stat) ReplaceLookupedName() {
	s.Local.Addr = netutil.ResolveAddr(s.Local.Addr)
	s.Foreign.Addr = netutil.ResolveAddr(s.Foreign.Addr)
}

// ReplacePortName replaces s.Local/Foreign.Port "0" to "*".
func (s *Stat) ReplacePortName() {
	if s.Local.Port == "0" {
		s.Local.Port = "*"
	}
	if s.Foreign.Port == "0" {
		s.Foreign.Port = "*"
	}
}
