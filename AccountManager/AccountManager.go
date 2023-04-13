package AccountManager

type NodeID = uint64

type AccountManager struct {
	WorkerNumberSet    map[uint32]string
	WorkerSetNumber    map[string]uint32
	VoterNumberSet	   map[uint32]string
	VoterSetNumber     map[string]uint32
	VoterSet           map[string]uint64
	WorkerSet          map[string]uint64
	WorkerCandidateSet map[string]uint64

	WorkerCandidateList []string
}
