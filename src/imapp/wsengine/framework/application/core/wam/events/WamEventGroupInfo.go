package events

import eventSerialize "ws/framework/plugin/event_serialize"

type WamEventGroupInfo struct {
	WAMessageEvent

	//ExitGroup                 float64 //0x2
	//GroupAddParticipants      float64 //0x3
	//GroupAudioCall            float64 //0x4
	//GroupClearChat            float64 //0x5
	//GroupDescription          float64 //0x6
	//GroupDisappearingMessages float64 //0x7
	//GroupEncryption           float64 //0x8
	//GroupExportChat           float64 //0x9
	GroupInfoVisit float64 //0x1
	//GroupMedia                float64 //0xa
	//GroupMembers              float64 //0xb
	//GroupMuteClick            float64 //0xc
	//GroupName                 float64 //0xd
	//GroupPhoto                float64 //0xe
	//GroupSearch               float64 //0xf
	//GroupShare                float64 //0x10
	//GroupStarredMessages      float64 //0x11
	//GroupVideoCall            float64 //0x12
	//GroupWallpaperAndSound    float64 //0x13
	//InviteToGroupGiaLink      float64 //0x14
	//ReportGroup               float64 //0x15
}

func (event *WamEventGroupInfo) InitFields(option interface{}) {
	event.GroupInfoVisit = 1
}

func (event *WamEventGroupInfo) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	//buffer.Body().
	//SerializeNumber(0x2, event.ExitGroup).
	//SerializeNumber(0x3, event.GroupAddParticipants).
	//SerializeNumber(0x4, event.GroupAudioCall).
	//SerializeNumber(0x5, event.GroupClearChat).
	//SerializeNumber(0x6, event.GroupDescription).
	//SerializeNumber(0x7, event.GroupDisappearingMessages).
	//SerializeNumber(0x8, event.GroupEncryption).
	//SerializeNumber(0x9, event.GroupExportChat).
	//SerializeNumber(0x1, event.GroupInfoVisit)
	//SerializeNumber(0xa, event.GroupMedia).
	//SerializeNumber(0xb, event.GroupMembers).
	//SerializeNumber(0xc, event.GroupMuteClick).
	//SerializeNumber(0xd, event.GroupName).
	//SerializeNumber(0xe, event.GroupPhoto).
	//SerializeNumber(0xf, event.GroupSearch).
	//SerializeNumber(0x10, event.GroupShare).
	//SerializeNumber(0x11, event.GroupStarredMessages).
	//SerializeNumber(0x12, event.GroupVideoCall).
	//SerializeNumber(0x13, event.GroupWallpaperAndSound).
	//SerializeNumber(0x14, event.InviteToGroupGiaLink).
	//SerializeNumber(0x15, event.ReportGroup)
	buffer.Footer().SerializeNumber(0x1, event.GroupInfoVisit)
}
