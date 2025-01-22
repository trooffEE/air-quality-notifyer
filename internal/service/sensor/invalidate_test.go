package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetAllApiIds() ([]int64, error) {
	args := m.Called()
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockRepo) EvictSensor(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// unused methods? Don't know how to make it better to fit interface
func (m *MockRepo) GetSensorByApiId(id int64) (*models.AirqualitySensor, error) { return nil, nil }
func (m *MockRepo) SaveSensor(sensor models.AirqualitySensor) error             { return nil }
func (m *MockRepo) GetSensorsByDistrictId(id int64) ([]models.AirqualitySensor, error) {
	return nil, nil
}

func TestInvalidateSensors(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockRepo)

	currentlySavedSensors := []int64{1, 2, 3, 4, 5}
	mockRepo.On("GetAllApiIds").Return(currentlySavedSensors, nil)
	mockRepo.On("EvictSensor", int64(4)).Return(nil)
	mockRepo.On("EvictSensor", int64(5)).Return(nil)

	service := &Service{repo: mockRepo}

	aliveSensors := []AqiSensorScriptScrapped{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}
	service.invalidateSensors(aliveSensors)

	mockRepo.AssertCalled(t, "GetAllApiIds")
	mockRepo.AssertCalled(t, "EvictSensor", int64(4))
	mockRepo.AssertCalled(t, "EvictSensor", int64(5))
	mockRepo.AssertNotCalled(t, "EvictSensor", int64(1))
	mockRepo.AssertNotCalled(t, "EvictSensor", int64(2))
	mockRepo.AssertNotCalled(t, "EvictSensor", int64(3))
}

func TestInvalidateSensors_ErrorGetSensorByApiId(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockRepo)
	service := &Service{repo: mockRepo}

	mockRepo.On("GetAllApiIds").Return([]int64{}, errors.New("something went wrong"))

	service.invalidateSensors([]AqiSensorScriptScrapped{{Id: int64(1)}})

	mockRepo.AssertCalled(t, "GetAllApiIds")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "EvictSensor")
}

func TestInvalidateSensors_ErrorEvictSensor(t *testing.T) {
	t.Parallel()
	mockRepo := new(MockRepo)
	service := &Service{repo: mockRepo}

	currentlySavedSensors := []int64{1, 2, 3, 4, 5}
	mockRepo.On("GetAllApiIds").Return(currentlySavedSensors, nil)
	mockRepo.On("EvictSensor", int64(4)).Return(errors.New("something went wrong"))
	mockRepo.On("EvictSensor", int64(5)).Return(nil)

	aliveSensors := []AqiSensorScriptScrapped{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}
	service.invalidateSensors(aliveSensors)

	mockRepo.AssertCalled(t, "GetAllApiIds")
	mockRepo.AssertCalled(t, "EvictSensor", int64(4))
	mockRepo.AssertCalled(t, "EvictSensor", int64(5))

	mockRepo.AssertExpectations(t)
}
