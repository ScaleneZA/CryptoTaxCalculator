package syncer

type Syncer interface {
	Sync(b Backends) error
}
