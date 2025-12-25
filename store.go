package main

import ("sync")

// Store será la estructura principal de los datos
type Store struct{ 
	mu sync.RWMutex // Mutex es para gestionar accesos concurrentes como lectura y escritura
	data map[string]string  //Mapa donde vamos a guardar los datos
}

//Función que servirá para inicializar la memoria
func NewStore() *Store{
	return &Store{
		data: make(map[string]string),
	}
}
//Set guardará un valor de forma segura
func (s *Store) Set(key, value string){
	s.mu.Lock()				// Bloqueamos para escritura
	defer s.mu.Unlock()	   // Desbloqueamos para salir de la función
	s.data[key] = value
}

//Get recuperará un valor de forma segura
func (s *Store) Get(key string) (string, bool){
	s.mu.RLock() 		   //Bloqueo de lectura en el que permitiremos múltiples lectores a la vez
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Delete borra una clave
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}