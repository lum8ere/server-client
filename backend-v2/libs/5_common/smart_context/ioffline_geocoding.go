package smart_context

// IGeocoder определяет интерфейс для геокодирования (он совпадает с интерфейсом в инфраструктуре)
type IGeocoder interface {
	// LocalGeocode получает координаты (lat, lon) по IP.
	LocalGeocode(ip string) (float64, float64, error)
}
