package main

import (
	"fmt"
	"log"
	"os"
	"time"

	agofig "github.com/akutz/gofig"
	lsctx "github.com/emccode/libstorage/api/context"
	lstypes "github.com/emccode/libstorage/api/types"
	lsutils "github.com/emccode/libstorage/api/utils"
	lsclient "github.com/emccode/libstorage/client"
)

const service = "virtualbox"
const volName = ""
const volID = "c3932aee-ee52-44fe-84a6-fcb426724ca5"

var ctx lstypes.Context
var nextDevice = ""

func init() {
	ctx = lsctx.WithValue(lsctx.Background(), lsctx.ServiceKey, service)
}

func main() {
	lsClient, err := makeClient()
	if err != nil {
		log.Fatal("error creating LibStorage client: ", err)
	}

	// are we already attached
	vols, err := lsClient.Storage().Volumes(
		ctx,
		&lstypes.VolumesOpts{Attachments: false},
	)

	var vol *lstypes.Volume
	for _, vol = range vols {
		if vol.Name == os.Args[1] {
			break
		}
	}

	// let's get volid
	vol, err = lsClient.Storage().VolumeInspect(
		ctx,
		vol.ID,
		&lstypes.VolumeInspectOpts{
			Attachments: true,
			Opts:        lsutils.NewStore(),
		},
	)

	attachedVol, token, err := lsClient.Storage().VolumeAttach(
		ctx,
		vol.ID,
		&lstypes.VolumeAttachOpts{
			Force: false,
			Opts:  lsutils.NewStore(),
		},
	)
	if err != nil {
		log.Fatal("failed to attach:", err)
	}

	success, devices, err := lsClient.Executor().WaitForDevice(
		ctx,
		&lstypes.WaitForDeviceOpts{
			Token:   token,
			Timeout: 13 * time.Second,
		},
	)

	if err != nil {
		log.Fatal("WaitForDevice failed: ", err)
	}
	var device string
	if success {
		for k, v := range devices.DeviceMap {
			fmt.Println(k, "-->", v)
			if token == k {
				device = v
			}
		}
	}
	fmt.Println("attached to ", device)
	// list existing known vols
	printVols([]*lstypes.Volume{attachedVol})
}

func makeClient() (lstypes.Client, error) {
	cfg := agofig.New()
	cfg.Set("libstorage.host", "tcp://:7979")
	return lsclient.New(ctx, cfg)
}

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
