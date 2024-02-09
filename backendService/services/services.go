package services

type Service interface {
	construct() error
}

/*
*

	construct -  for construction of a  service
	@Param  s Service
*/
func construct(s Service) error {
	return s.construct()
}

/** Providers - will create  the servies  and if  the servoice  can  be created  fall the system
 */
func Providers() {

}
