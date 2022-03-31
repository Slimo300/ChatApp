package ws

import (
	"github.com/Slimo300/ChatApp/backend/src/communication"
	"github.com/Slimo300/ChatApp/backend/src/models"
)

// Group created updates user groups list to which he listens
func (h *Hub) GroupCreated(userID, groupID int) {
	for client := range h.clients {
		if client.id == userID {
			client.groups = append(client.groups, int64(groupID))
		}
	}
}

// Deletes group from every user that is subscribed to it and sends information via websocket to user
func (h *Hub) GroupDeleted(groupID int) {
	for client := range h.clients {
		for i, group := range client.groups {
			if group == int64(groupID) {
				client.groups = append(client.groups[:i], client.groups[:i+1]...)
				client.send <- &communication.Action{Action: "DELETE_GROUP", Group: groupID}
			}
		}
	}
}

// Adds subscription to member groups and sends info to other members in group
func (h *Hub) MemberAdded(member models.Member) {
	for client := range h.clients {
		if client.id == int(member.UserID) {
			client.groups = append(client.groups, int64(member.GroupID))
			continue
		}
		for _, group := range client.groups {
			if member.GroupID == uint(group) {
				client.send <- &communication.Action{Action: "ADD_MEMBER", Member: member}
			}
		}
	}
}

// Deletes member subscription and sends info about it to other members in group
func (h *Hub) MemberDeleted(member models.Member) {
	for client := range h.clients {
		for i, group := range client.groups {
			// if user is a member of group
			if group == int64(member.GroupID) {
				// if user is the one to be deleted
				if client.id == int(member.UserID) {
					client.groups = append(client.groups[:i], client.groups[:i+1]...)
				} else {
					client.send <- &communication.Action{Action: "DELETE_MEMBER", Member: member}
				}
			}
		}
	}
}
