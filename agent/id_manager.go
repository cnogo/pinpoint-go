package agent

import (
	"sync"
	"sync/atomic"
)

type APIMeta struct {
	Api string
	ApiID int32
	Line int32
	Type int32
}

type SQLMeta struct {
	SqlID int32
	Sql   string
}

type StrMeta struct {
	StrID int32
	Str string
}


type IDManager struct {
	ID int32
	ID2MetaInfo map[int32]string
	MetaInfo2ID map[string]int32
	ID2MetaObject sync.Map
}

func NewIDManager() *IDManager {
	return &IDManager{
		ID: 0,
		ID2MetaInfo: make(map[int32]string),
		MetaInfo2ID: make(map[string]int32),
	}
}

func (p *IDManager) SetMetaObject(id int32, metaObject interface{}) {
	p.ID2MetaObject.Store(id, metaObject)
}

//返回apiID，和 是否是新的id
func (p *IDManager) GetOrCreateID(metaInfo string) (int32, bool) {
	newValue := false
	ID, ok := p.MetaInfo2ID[metaInfo]
	if !ok  {
		ID = atomic.AddInt32(&p.ID, 1)
		p.ID2MetaInfo[ID] = metaInfo
		p.MetaInfo2ID[metaInfo] = ID
		newValue = true
	}

	return ID, newValue
}

func (p *IDManager) GetID(metaInfo string) int32 {
	if id, ok := p.MetaInfo2ID[metaInfo]; ok {
		return id
	}

	return 0
}

//根据id获取api
func (p *IDManager) GetMetaInfo(id int32) string {
	if api, ok := p.ID2MetaInfo[id]; ok {
		return api
	}

	return ""
}