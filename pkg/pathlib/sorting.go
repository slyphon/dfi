package pathlib

type LexOrderPosix []PosixPath

func (o LexOrderPosix) Len() int           { return len(o) }
func (o LexOrderPosix) Less(i, j int) bool { return o[i].String() < o[j].String() }
func (o LexOrderPosix) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

type LexOrderPure []PurePath

func (o LexOrderPure) Len() int           { return len(o) }
func (o LexOrderPure) Less(i, j int) bool { return o[i].String() < o[j].String() }
func (o LexOrderPure) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
