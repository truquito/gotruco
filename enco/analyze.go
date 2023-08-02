package enco

func Analyze(pkts []Packet) map[CodMsg]bool {
	out := make(map[CodMsg]bool)

	for _, pkt := range pkts {
		out[pkt.Message.Cod()] = true
	}

	return out
}
