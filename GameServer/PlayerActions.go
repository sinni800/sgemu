package GameServer

import (
	. "github.com/sinni800/sgemu/Data"
)

func (c *GClient) RemoveUnit(unitID uint32) bool {
	unit, exists := c.Units[unitID]
	if exists {
		delete(c.Units, unitID)
		_, exists = c.Player.UnitsData[unit.DBID]
		if exists {
			delete(c.Player.UnitsData, unit.DBID)
		}

		SendRemoveUnit(c, unitID)
		return true
	}
	return false
}

func (c *GClient) AddUnit(unitName string, customName string) *Unit {
	unitdb := c.Player.AddUnit(unitName)
	unitdb.CustomName = customName
	id, r := c.Server.IDG.Next()
	if !r {
		c.Log().Println_Warning("No more ids left - server is full!")
		return nil
	}
	name, e := Units[unitdb.Name]
	if !e {
		c.Log().Println_Warning("Unit name does not exists")
		return nil
	}
	unit := &Unit{unitdb, id, c.Player, name, 0, 0}
	c.Units[id] = unit

	SendNewUnit(c, unit)
	SendUnitInventory(c, unit)

	return unit
}
