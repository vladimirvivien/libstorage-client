package main

import (
	"fmt"
	"log"
	"os"

	agofig "github.com/akutz/gofig"
	lsctx "github.com/emccode/libstorage/api/context"
	lstypes "github.com/emccode/libstorage/api/types"
	lsclient "github.com/emccode/libstorage/client"
)

const service = "virtualbox"

var ctx lstypes.Context

func init() {
	ctx = lsctx.WithValue(lsctx.Background(), lsctx.ServiceKey, service)
}

func main() {
	lsClient, err := makeClient()
	if err != nil {
		log.Fatal("error creating LibStorage client: ", err)
	}

	vol, err := lsClient.Storage().VolumeInspect(
		ctx,
		os.Args[1],
		&lstypes.VolumeInspectOpts{Attachments: true},
	)
	if err != nil {
		log.Fatal("Unable to do VolumesByService(): ", err)
	}

	printVols([]*lstypes.Volume{vol})
}

func makeClient() (lstypes.Client, error) {

	cfg := agofig.New()
	cfg.Set("libstorage.host", "tcp://:7979")
	return lsclient.New(ctx, cfg)
}

// func printVols(vols lstypes.VolumeMap) {
// 	for id, vol := range vols {
// 		fmt.Printf("Volume ID: %s, name: %s, mountPoint:%s\n", id, vol.VolumeName(), vol.MountPoint())
// 	}
// }

func printVols(vols []*lstypes.Volume) {
	for _, vol := range vols {

		fmt.Printf("Volume ID: %s\nname: %s\nmountPoint:%s\n",
			vol.ID,
			vol.VolumeName(),
			vol.MountPoint(),
		)
		fmt.Print("Attachments: ", len(vol.Attachments), "\n")
		for _, attach := range vol.Attachments {
			fmt.Println(" - DeviceName:", attach.DeviceName)
			fmt.Println(" - InstanceID;", attach.InstanceID)
			fmt.Println(" - Status:", attach.Status)
		}
		fmt.Println("--------------")
	}
}
