package offheap

func (p *ChunkPool) grawMmapBytesList() error {
	mmapBytes, err := AllocMmapBytes(int(p.perMmapBytesSize))
	if err != nil {
		return err
	}
	p.mmapBytesList = append(p.mmapBytesList, &mmapBytes)
	p.currentMmapBytes = p.mmapBytesList[len(p.mmapBytesList)-1]

	return nil
}
