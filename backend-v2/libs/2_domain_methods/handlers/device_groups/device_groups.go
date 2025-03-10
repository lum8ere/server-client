package device_groups

import (
	"backed-api-v2/libs/3_generated_models/model"
	"backed-api-v2/libs/5_common/smart_context"
	"backed-api-v2/libs/5_common/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// GetDeviceGroupsHandler returns all device groups.
func GetDeviceGroupsHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	var groups []model.DeviceGroup
	if err := sctx.GetDB().Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to get device groups: %w", err)
	}
	return groups, nil
}

// CreateDeviceGroupHandler creates a new device group.
func CreateDeviceGroupHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	name, ok := args.GetStringValue("name")
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required")
	}
	description, _ := args.GetStringValue("description")

	group := model.DeviceGroup{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := sctx.GetDB().Create(&group).Error; err != nil {
		return nil, fmt.Errorf("failed to create device group: %w", err)
	}
	return group, nil
}

// UpdateDeviceGroupHandler updates an existing device group.
func UpdateDeviceGroupHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	id, ok := args.GetStringValue("id")
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}

	var group model.DeviceGroup
	if err := sctx.GetDB().Where("id = ?", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device group not found")
		}
		return nil, fmt.Errorf("failed to find device group: %w", err)
	}

	if name, ok := args.GetStringValue("name"); ok && name != "" {
		group.Name = name
	}
	if description, ok := args.GetStringValue("description"); ok {
		group.Description = description
	}

	// При необходимости можно обновить и временную метку

	if err := sctx.GetDB().Save(&group).Error; err != nil {
		return nil, fmt.Errorf("failed to update device group: %w", err)
	}
	return group, nil
}

// DeleteDeviceGroupHandler deletes a device group by its ID.
func DeleteDeviceGroupHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	id, ok := args.GetStringValue("id")
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	if err := sctx.GetDB().Delete(&model.DeviceGroup{}, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to delete device group: %w", err)
	}
	return map[string]string{"status": "deleted"}, nil
}

// AssignDeviceToGroupHandler assigns a device to a group by updating the devices table.
func AssignDeviceToGroupHandler(sctx smart_context.ISmartContext, args types.ANY_DATA) (interface{}, error) {
	deviceId, ok := args.GetStringValue("device_id")
	if !ok || deviceId == "" {
		return nil, fmt.Errorf("device_id is required")
	}
	groupId, ok := args.GetStringValue("group_id")
	if !ok || groupId == "" {
		return nil, fmt.Errorf("group_id is required")
	}
	// Обновляем столбец group_id для указанного устройства
	if err := sctx.GetDB().Model(&model.Device{}).
		Where("id = ?", deviceId).
		Update("group_id", groupId).Error; err != nil {
		return nil, fmt.Errorf("failed to assign device to group: %w", err)
	}
	return map[string]string{"status": "assigned"}, nil
}
