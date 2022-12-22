package contactService

import (
	"fmt"
	"time"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/contact/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
)

var _ containerInterface.IContactService = &Contact{}

var keyTemplate = "contact_%s_%s"

// Contact .
type Contact struct {
	containerInterface.BaseService
}

func (s *Contact) createKey(id string) string {
	// 取后7位作为关键词
	var number string
	if len(id) < 7 {
		number = id[0:]
	} else {
		number = id[len(id)-7:]
	}

	return fmt.Sprintf(keyTemplate, s.JID.User, number)
}

// CreateContact 新增联系人
func (s *Contact) CreateContact(dstJID string, aliasPhoneNumber string) {
	cacheKey := s.createKey(dstJID)

	contact := contactDB.NewContactByDstJIDNumber(s.JID.User, dstJID)
	contact.DstPhoneNumber = aliasPhoneNumber

	_, err := databaseTools.Create(database.MasterDB(), &contact)
	if err != nil {
		s.Logger.Errorf("add contact failed. %s %v", dstJID, err)
	}

	s.AppIocContainer.ResolveMemoryCache().Cache(cacheKey, &contact)
}

// BatchCreateAddressBookContactByJIDList 批量添加到通讯录
func (s *Contact) BatchCreateAddressBookContactByJIDList(jidList []string, aliasPhoneNumber []string) {
	inContactList, _ := contactDB.FindList(database.MasterDB(), s.JID.User, jidList)

	mapping := make(map[string]int)
	for i := range inContactList {
		mapping[inContactList[i].JID] = i
	}

	now := time.Now().Unix()
	newContacts := make([]contactDB.Contact, 0)

	for i, jid := range jidList {
		contactIdx, in := mapping[jid]
		if in {
			inContactList[contactIdx].UpdateInAddressBook(true)
			inContactList[contactIdx].UpdateAddTime(now)

			_ = s.save(jid, &inContactList[contactIdx])
			continue
		}

		// 新建并添加
		contact := contactDB.NewContactByDstJIDNumber(s.JID.User, jidList[i])
		contact.DstPhoneNumber = aliasPhoneNumber[i]
		contact.AddTime = now
		contact.InAddressBook = true

		newContacts = append(newContacts, contact)
	}

	_, err := databaseTools.BatchCreate(database.MasterDB(), newContacts)
	if err != nil {
		s.Logger.Errorf("batch add contact failed. %v", err)
	}
}

// DeleteByJID .
func (s *Contact) DeleteByJID(dstJID string) error {
	contact := contactDB.NewContactByDstJIDNumber(s.JID.User, dstJID)

	err := contactDB.Delete(database.MasterDB(), &contact)
	if err != nil {
		s.Logger.Errorf("delete contact failed. %s %v", dstJID, err)
	}

	s.AppIocContainer.ResolveMemoryCache().UnCache(s.createKey(dstJID))

	return err
}

// CreateAddressBookContactByJID 添加到通讯录
func (s *Contact) CreateAddressBookContactByJID(dstJID string, aliasPhoneNumber string) error {
	cacheContact := s.FindByJID(dstJID)

	// 有过联系人记录
	if cacheContact != nil {
		// 设置添加属性
		cacheContact.UpdateInAddressBook(true)
		cacheContact.UpdateAddTime(time.Now().Unix())

		return s.save(dstJID, cacheContact)
	}

	// 没有的 就新建 并设置添加
	cacheKey := s.createKey(dstJID)

	contact := contactDB.NewContactByDstJIDNumber(s.JID.User, dstJID)
	contact.DstPhoneNumber = aliasPhoneNumber
	contact.AddTime = time.Now().Unix()
	contact.InAddressBook = true

	_, err := databaseTools.Create(database.MasterDB(), &contact)
	if err != nil {
		s.Logger.Errorf("add contact failed. %s %v", dstJID, err)
		return err
	}

	s.AppIocContainer.ResolveMemoryCache().Cache(cacheKey, &contact)

	return nil
}

// FindByPhoneNumber .
func (s *Contact) FindByPhoneNumber(phoneNumber string) *contactDB.Contact {
	cacheKey := s.createKey(phoneNumber)

	contact, ok := s.AppIocContainer.ResolveMemoryCache().FindInCache(cacheKey)
	if ok {
		return contact.(*contactDB.Contact)
	}

	c := contactDB.NewContactByDstPhoneNumber(s.JID.User, phoneNumber)

	err := databaseTools.FindByPrimaryKey(database.MasterDB(), &c)
	if err != nil {
		return nil
	}

	s.AppIocContainer.ResolveMemoryCache().Cache(cacheKey, &c)

	return &c
}

// FindByJID .
func (s *Contact) FindByJID(dstJID string) *contactDB.Contact {
	cacheKey := s.createKey(dstJID)

	contact, ok := s.AppIocContainer.ResolveMemoryCache().FindInCache(cacheKey)
	if ok {
		return contact.(*contactDB.Contact)
	}

	c := contactDB.NewContactByDstJIDNumber(s.JID.User, dstJID)

	err := databaseTools.FindByPrimaryKey(database.MasterDB(), &c)
	if err != nil {
		return nil
	}

	s.AppIocContainer.ResolveMemoryCache().Cache(cacheKey, &c)

	return &c
}

// ContextExecute .
func (s *Contact) ContextExecute(dstJID string, f func(*contactDB.Contact)) {
	contact := s.FindByJID(dstJID)
	if contact == nil {
		return
	}

	f(contact)

	_ = s.save(dstJID, contact)
}

func (s *Contact) save(dstJID string, contact *contactDB.Contact) error {
	cacheKey := s.createKey(dstJID)

	_, err := databaseTools.Save(database.MasterDB(), contact)
	if err != nil {
		s.Logger.Errorf("save contact failed. %s %v", dstJID, err)
	}

	s.AppIocContainer.ResolveMemoryCache().Cache(cacheKey, contact)

	return err
}

// CleanupAllData .
func (s *Contact) CleanupAllData() {
	_, _ = contactDB.DeleteByJID(database.MasterDB(), s.JID.User)
}
