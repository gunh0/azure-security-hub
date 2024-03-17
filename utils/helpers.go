package utils

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

func StringPtr(s string) *string {
	return &s
}

func AutoProvisionPtr(a armsecurity.AutoProvision) *armsecurity.AutoProvision {
	return &a
}

// MockPager creates a mock pager for testing
type MockPager struct {
	items []*armsecurity.AutoProvisioningSetting
	index int
}

func NewMockPager(items []armsecurity.AutoProvisioningSetting) *MockPager {
	ptrItems := make([]*armsecurity.AutoProvisioningSetting, len(items))
	for i := range items {
		ptrItems[i] = &items[i]
	}
	return &MockPager{
		items: ptrItems,
		index: 0,
	}
}

func (m *MockPager) More() bool {
	return m.index < len(m.items)
}

func (m *MockPager) NextPage(ctx context.Context) (armsecurity.AutoProvisioningSettingsClientListResponse, error) {
	if !m.More() {
		return armsecurity.AutoProvisioningSettingsClientListResponse{}, nil
	}
	start := m.index
	m.index = len(m.items)
	return armsecurity.AutoProvisioningSettingsClientListResponse{
		AutoProvisioningSettingList: armsecurity.AutoProvisioningSettingList{
			Value: m.items[start:],
		},
	}, nil
}
