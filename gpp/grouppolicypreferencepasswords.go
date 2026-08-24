package gpp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/TheManticoreProject/Manticore/crypto/gppp"
	smbclient "github.com/TheManticoreProject/Manticore/network/smb/client"
	"github.com/TheManticoreProject/Manticore/windows/fileflags"
)

// transferChunk bounds how many bytes are buffered in memory per read iteration
// when downloading a remote file. The SMB client further splits each call to fit
// the negotiated MaxBufferSize, so this only caps our own buffering.
const transferChunk = 0xFF00

// readRemoteFile streams the remote file at the share-relative path on the
// client's current tree into w.
func readRemoteFile(client *smbclient.Client, path string, w io.Writer) error {
	h, err := client.OpenFile(path, smbclient.OpenOptions{
		DesiredAccess:     fileflags.GENERIC_READ,
		ShareAccess:       fileflags.FILE_SHARE_READ,
		CreateDisposition: fileflags.FILE_OPEN,
		CreateOptions:     fileflags.FILE_NON_DIRECTORY_FILE,
	})
	if err != nil {
		return err
	}
	defer client.CloseFile(h)

	var offset uint64
	for {
		chunk, rerr := client.ReadFile(h, offset, transferChunk)
		if len(chunk) > 0 {
			if _, werr := w.Write(chunk); werr != nil {
				return werr
			}
			offset += uint64(len(chunk))
		}
		// The client surfaces an error at end of file; a short read also marks
		// the end. Either way we stop after consuming whatever was returned.
		if rerr != nil || uint32(len(chunk)) < transferChunk {
			break
		}
	}
	return nil
}

// XML structure for Properties
type User_Properties struct {
	Action    string `xml:"action,attr"`
	NewName   string `xml:"newName,attr"`
	UserName  string `xml:"userName,attr"`
	CPassword string `xml:"cpassword,attr"`
}

// XML structure for User
type User struct {
	Properties User_Properties `xml:"Properties"`
}

// XML structure for Groups
type Groups struct {
	Users []User `xml:"User"`
}

// XML structure for Trigger
type Trigger struct {
	Interval     string `xml:"interval,attr"`
	Type         string `xml:"type,attr"`
	StartHour    string `xml:"startHour,attr"`
	StartMinutes string `xml:"startMinutes,attr"`
	BeginYear    string `xml:"beginYear,attr"`
	BeginMonth   string `xml:"beginMonth,attr"`
	BeginDay     string `xml:"beginDay,attr"`
	HasEndDate   string `xml:"hasEndDate,attr"`
	RepeatTask   string `xml:"repeatTask,attr"`
	Week         string `xml:"week,attr"`
	Days         string `xml:"days,attr"`
	Months       string `xml:"months,attr"`
}

// XML structure for Triggers
type Triggers struct {
	Trigger []Trigger `xml:"Trigger"`
}

// XML structure for Task Properties
type TaskProperties struct {
	DeleteWhenDone         string   `xml:"deleteWhenDone,attr"`
	StartOnlyIfIdle        string   `xml:"startOnlyIfIdle,attr"`
	StopOnIdleEnd          string   `xml:"stopOnIdleEnd,attr"`
	NoStartIfOnBatteries   string   `xml:"noStartIfOnBatteries,attr"`
	StopIfGoingOnBatteries string   `xml:"stopIfGoingOnBatteries,attr"`
	SystemRequired         string   `xml:"systemRequired,attr"`
	Action                 string   `xml:"action,attr"`
	Name                   string   `xml:"name,attr"`
	AppName                string   `xml:"appName,attr"`
	Args                   string   `xml:"args,attr"`
	StartIn                string   `xml:"startIn,attr"`
	Comment                string   `xml:"comment,attr"`
	RunAs                  string   `xml:"runAs,attr"`
	CPassword              string   `xml:"cpassword,attr"`
	Enabled                string   `xml:"enabled,attr"`
	Triggers               Triggers `xml:"Triggers"`
}

// XML structure for Task
type Task struct {
	Clsid      string         `xml:"clsid,attr"`
	Name       string         `xml:"name,attr"`
	Image      string         `xml:"image,attr"`
	Changed    string         `xml:"changed,attr"`
	UID        string         `xml:"uid,attr"`
	Properties TaskProperties `xml:"Properties"`
}

// XML structure for ScheduledTasks
type ScheduledTasks struct {
	Clsid string `xml:"clsid,attr"`
	Tasks []Task `xml:"Task"`
}

type CPasswordEntry struct {
	RunAs     string
	UserName  string
	NewName   string
	CPassword string
	Password  string
}

type GroupPolicyPreferencePasswordsFound struct {
	Entries map[string][]*CPasswordEntry
}

func (r *GroupPolicyPreferencePasswordsFound) CallbackFunctionCPassword(client *smbclient.Client, share string, pathToFile string) error {
	elements := strings.Split(pathToFile, ".")
	extension := strings.ToLower(elements[len(elements)-1])

	if strings.EqualFold(extension, "xml") {
		uncPathToFile := fmt.Sprintf("\\\\%s\\%s\\%s", client.ServerIdentity().DNSComputerName, share, pathToFile)

		buffer := bytes.NewBuffer([]byte{})

		err := readRemoteFile(client, pathToFile, buffer)
		if err != nil {
			return err
		}

		cpasswords, err := ExtractCPasswordsFromRawXML(buffer)
		if err != nil {
			return fmt.Errorf("error extracting GPP passwords from %s: %w", pathToFile, err)
		}

		if len(cpasswords) != 0 {
			if _, ok := r.Entries[uncPathToFile]; !ok {
				r.Entries[uncPathToFile] = make([]*CPasswordEntry, 0)
			}
			r.Entries[uncPathToFile] = append(r.Entries[uncPathToFile], cpasswords...)
		}
	}

	return nil
}

func ExtractCPasswordsFromRawXML(buffer *bytes.Buffer) ([]*CPasswordEntry, error) {
	// Create an instance of Groups to hold the parsed data
	foundCpasswords := make([]*CPasswordEntry, 0)

	if strings.Contains(buffer.String(), "</ScheduledTasks>") {
		// Parse the XML data to search for ScheduledTasks
		scheduledtasks := ScheduledTasks{}

		if err := xml.NewDecoder(buffer).Decode(&scheduledtasks); err != nil {
			return nil, fmt.Errorf("error parsing ScheduledTasks XML: %w", err)
		}

		// Extract and print the desired properties
		for _, task := range scheduledtasks.Tasks {
			if len(task.Properties.CPassword) != 0 {
				password, err := gppp.GPPPDecryptBase64(task.Properties.CPassword)
				if err != nil {
					return nil, fmt.Errorf("error decrypting cpassword for scheduled task %q: %w", task.Name, err)
				}
				entry := CPasswordEntry{
					RunAs:     task.Properties.RunAs,
					CPassword: task.Properties.CPassword,
					Password:  password,
				}
				foundCpasswords = append(foundCpasswords, &entry)
			}
		}

	} else if strings.Contains(buffer.String(), "</Groups>") {
		// Parse the XML data to search for Users
		groups := Groups{}

		if err := xml.NewDecoder(buffer).Decode(&groups); err != nil {
			return nil, fmt.Errorf("error parsing Groups XML: %w", err)
		}

		// Extract and print the desired properties
		for _, user := range groups.Users {
			if len(user.Properties.CPassword) != 0 {
				password, err := gppp.GPPPDecryptBase64(user.Properties.CPassword)
				if err != nil {
					return nil, fmt.Errorf("error decrypting cpassword for user %q: %w", user.Properties.UserName, err)
				}
				entry := CPasswordEntry{
					UserName:  user.Properties.UserName,
					NewName:   user.Properties.NewName,
					CPassword: user.Properties.CPassword,
					Password:  password,
				}
				foundCpasswords = append(foundCpasswords, &entry)
			}
		}
	}

	return foundCpasswords, nil
}
