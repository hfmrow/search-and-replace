package files

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

/* Const for FileMode
Usage to Create any directories needed to put this file in them:
     var dir_file_mode os.FileMode
     dir_file_mode = os.ModeDir | (OS_USER_RWX | OS_ALL_R)
     os.MkdirAll(dir_str, dir_file_mode)

	fmt.Printf("%o\n%o\n%o\n%o\n",
		os.ModePerm&OS_ALL_RWX,
		os.ModePerm&OS_USER_RW|OS_GROUP_R|OS_OTH_R,
		os.ModePerm&OS_USER_RW|OS_GROUP_RW|OS_OTH_R,
		os.ModePerm&OS_USER_RWX|OS_GROUP_RWX|OS_OTH_R)
*/
const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_ALL_X
)

// DispRights: display right.
//i.e: g.DispRights(g.OS_GROUP_RW | g.OS_USER_RW | g.OS_ALL_R) -> 664
func DispRights(value int) {
	fmt.Printf("%o\n", value)
}

// ChangeFileOwner: set file owner to real user instead of root.
func ChangeFileOwner(filename string) (err error) {
	var realUser *user.User
	cmd := exec.Command("id", "-u")
	output, _ := cmd.Output()
	if string(output[:len(output)-1]) == "0" {
		if realUser, err = user.Lookup(os.Getenv("SUDO_USER")); err == nil {
			if uid, err := strconv.Atoi(realUser.Uid); err == nil {
				if gid, err := strconv.Atoi(realUser.Gid); err == nil {
					err = os.Chown(filename, uid, gid)
				}
			}
		}
	}
	return err
}
