package etcdutil

import (
	"path"

	"github.com/coreos/go-etcd/etcd"
)

// Creates the given dir (and all of its parent directories if they don't
// already exist). Will not return an error if the given directory already
// exists
func MkDirP(ec *etcd.Client, dir string) error {
    parts := make([]string, 0, 4)
    for {
        parts = append(parts, dir)
        dir = path.Dir(dir)
        if dir == "/" {
            break
        }
    }

    for i := range parts {
        ai := len(parts) - i - 1
        _, err := ec.CreateDir(parts[ai], 0)
        if err != nil && err.(*etcd.EtcdError).ErrorCode != 105 {
            return err
        }
    }
    return nil
}

// Returns the contents of a directory as a list of absolute paths
func Ls(ec *etcd.Client, dir string) ([]string, error) {
    r, err := ec.Get(dir, false, false)
    if err != nil {
        return nil, err
    }

    dirNode := r.Node
    ret := make([]string, len(dirNode.Nodes))
    for i, node := range dirNode.Nodes {
        ret[i] = node.Key
    }

    return ret, nil
}
