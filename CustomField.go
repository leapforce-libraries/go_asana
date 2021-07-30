package asana

// CustomFieldProject stores CustomFieldProject from Asana
//
type CustomFieldProject struct {
	ID              string       `json:"gid"`
	ResourceType    string       `json:"resource_type"`
	DisplayValue    string       `json:"display_value"`
	Enabled         bool         `json:"enabled"`
	EnumOptions     []EnumOption `json:"enum_options"`
	Name            string       `json:"name"`
	NumberValue     float64      `json:"number_value"`
	ResourceSubtype string       `json:"resource_subtype"`
	TextValue       string       `json:"text_value"`
}

// CustomFieldTask stores CustomFieldTask from Asana
//
type CustomFieldTask struct {
	ID                      string        `json:"gid"`
	ResourceType            string        `json:"resource_type"`
	CreatedBy               Object        `json:"created_by"`
	CurrencyCode            *string       `json:"currency_code"`
	CustomLabel             *string       `json:"custom_label"`
	CustomLabelPosition     *string       `json:"custom_label_position"`
	Description             string        `json:"description"`
	DisplayValue            string        `json:"display_value"`
	Enabled                 bool          `json:"enabled"`
	EnumOptions             *[]EnumOption `json:"enum_options"`
	EnumValue               *EnumOption   `json:"enum_value"`
	Format                  string        `json:"custom"`
	HasNotificationsEnabled bool          `json:"has_notifications_enabled"`
	IsGlobalToWorkspace     bool          `json:"is_global_to_workspace"`
	Name                    string        `json:"name"`
	NumberValue             *float64      `json:"number_value"`
	Precision               *int64        `json:"precision"`
	ResourceSubtype         string        `json:"resource_subtype"`
	TextValue               *string       `json:"text_value"`
}
