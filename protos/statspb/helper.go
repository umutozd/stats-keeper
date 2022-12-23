package statspb

// GetComponentType is a helper function to get the ComponentType of this entity.
func (se *StatisticEntity) GetComponentType() ComponentType {
	switch se.Component.(type) {
	case *StatisticEntity_Counter:
		return ComponentType_COUNTER
	case *StatisticEntity_Date:
		return ComponentType_DATE
	default:
		return ComponentType_NONE
	}
}
