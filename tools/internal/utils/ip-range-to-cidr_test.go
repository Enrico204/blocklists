package utils

import (
	"github.com/praserx/ipconv"
	"log"
	"net"
	"testing"
)

func TestIPrangeToCIDR_v4(t *testing.T) {
	t.Run("simple range 24", func(t *testing.T) {
		var start = net.IPv4(1, 1, 1, 0)
		var end = net.IPv4(1, 1, 1, 255)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "1.1.1.0/24" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple range 27", func(t *testing.T) {
		var start = net.IPv4(1, 1, 1, 0)
		var end = net.IPv4(1, 1, 1, 31)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "1.1.1.0/27" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple 2 range 27", func(t *testing.T) {
		var start = net.IPv4(1, 0, 0, 0)
		var end = net.IPv4(1, 0, 0, 31)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "1.0.0.0/27" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple range 32", func(t *testing.T) {
		var start = net.IPv4(1, 1, 1, 0)
		var end = net.IPv4(1, 1, 1, 0)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "1.1.1.0/32" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple range 8", func(t *testing.T) {
		var start = net.IPv4(1, 0, 0, 0)
		var end = net.IPv4(1, 255, 255, 255)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error ", len(networks))
		} else if networks[0].String() != "1.0.0.0/8" {
			t.Fatal("cidr error ", networks[0].String())
		}
	})
	t.Run("range 1", func(t *testing.T) {
		var start = net.IPv4(1, 1, 1, 111)
		var end = net.IPv4(1, 1, 1, 120)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 3 {
			t.Fatal("length error ", len(networks))
		} else if networks[0].String() != "1.1.1.111/32" {
			t.Fatal("cidr error1 ", networks[0].String())
		} else if networks[1].String() != "1.1.1.112/29" {
			t.Fatal("cidr error2 ", networks[1].String())
		} else if networks[2].String() != "1.1.1.120/32" {
			t.Fatal("cidr error3 ", networks[2].String())
		}
	})
	t.Run("bogon range 8", func(t *testing.T) {
		var start = net.IPv4(0, 0, 0, 0)
		var end = net.IPv4(0, 255, 255, 255)
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "0.0.0.0/8" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("test hangs with 255.255.255.255 as end IP", func(t *testing.T) {
		var start = net.IPv4(255, 255, 255, 234)
		var end = net.IPv4(255, 255, 255, 255)
		networks := IPrangeToCIDR(start, end)
		log.Println(networks)
	})
}

func TestIPrangeToCIDR_v6(t *testing.T) {
	t.Run("simple range 64", func(t *testing.T) {
		var start = net.ParseIP("2001:0db8:85a3:0000:0000:0000:0000:0000")
		var end = net.ParseIP("2001:0db8:85a3:0000:ffff:ffff:ffff:ffff")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "2001:db8:85a3::/64" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple range 59", func(t *testing.T) {
		var start = net.ParseIP("2001:0db8:85a3:0000:0000:0000:0000:0000")
		var end = net.ParseIP("2001:0db8:85a3:001f:ffff:ffff:ffff:ffff")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "2001:db8:85a3::/59" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple 2 range 59", func(t *testing.T) {
		var start = net.ParseIP("2001:0000:85a3:0000:0000:0000:0000:0000")
		var end = net.ParseIP("2001:0000:85a3:001f:ffff:ffff:ffff:ffff")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "2001:0:85a3::/59" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple range 128", func(t *testing.T) {
		var start = net.ParseIP("2001:0000:85a3:001f:2001:0000:85a3:001f")
		var end = net.ParseIP("2001:0000:85a3:001f:2001:0000:85a3:001f")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "2001:0:85a3:1f:2001:0:85a3:1f/128" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("simple range 8", func(t *testing.T) {
		var start = net.ParseIP("2000:0000:0000:0000:0000:0000:0000:0000")
		var end = net.ParseIP("20ff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error ", len(networks))
		} else if networks[0].String() != "2000::/8" {
			t.Fatal("cidr error ", networks[0].String())
		}
	})
	t.Run("range 1", func(t *testing.T) {
		var start = net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:732f")
		var end = net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7338")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 3 {
			t.Fatal("length error ", len(networks))
		} else if networks[0].String() != "2001:db8:85a3::8a2e:370:732f/128" {
			t.Fatal("cidr error1 ", networks[0].String())
		} else if networks[1].String() != "2001:db8:85a3::8a2e:370:7330/125" {
			t.Fatal("cidr error2 ", networks[1].String())
		} else if networks[2].String() != "2001:db8:85a3::8a2e:370:7338/128" {
			t.Fatal("cidr error3 ", networks[2].String())
		}
	})
	t.Run("bogon range 8", func(t *testing.T) {
		var start = net.ParseIP("::")
		var end = net.ParseIP("0000:0000:07ff:ffff:ffff:ffff:ffff:ffff")
		networks := IPrangeToCIDR(start, end)
		if len(networks) != 1 {
			t.Fatal("length error", len(networks))
		} else if networks[0].String() != "::/37" {
			t.Fatal("cidr error", networks[0].String())
		}
	})
	t.Run("test hangs with ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff as end IP", func(t *testing.T) {
		var start = net.ParseIP("ffff:ffff:ffff:ffff:0000:0000:0000:0000")
		var end = net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
		networks := IPrangeToCIDR(start, end)
		log.Println(networks)
	})
}

func FuzzIPrangeToCIDR_v4(f *testing.F) {
	f.Fuzz(func(t *testing.T, start uint32, end uint32) {
		var startIP = ipconv.IntToIPv4(start)
		var endIP = ipconv.IntToIPv4(end)
		networks := IPrangeToCIDR(startIP, endIP)
		if len(networks) == 0 {
			t.Errorf("empty networks with %s %s", startIP.String(), endIP.String())
		}
	})
}

func FuzzIPrangeToCIDR_v6(f *testing.F) {
	f.Fuzz(func(t *testing.T, startHI uint64, startLO uint64, endHI uint64, endLO uint64) {
		var startIP = ipconv.IntToIPv6(startHI, startLO)
		var endIP = ipconv.IntToIPv6(endHI, endLO)
		networks := IPrangeToCIDR(startIP, endIP)
		if len(networks) == 0 {
			t.Errorf("empty networks with %s %s", startIP.String(), endIP.String())
		}
	})
}
