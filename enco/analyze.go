package enco

func Analyze(pkts []Envelope) map[CodMsg]bool {
	out := make(map[CodMsg]bool)

	for _, pkt := range pkts {
		out[pkt.Message.Cod()] = true
	}

	return out
}
